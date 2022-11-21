package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/ProninIgorr/alcomarket-nearst/internal/store"
	"github.com/ProninIgorr/alcomarket-nearst/pkg/helpers/response"
	"github.com/ProninIgorr/alcomarket-nearst/pkg/helpers/utils"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type MetaDistance struct {
	Exist    bool
	Distance float64
	Shop     store.Shop
}

type Response struct {
	Url  string `json:"url"`
	Name string `json:"name"`
	Addr string `json:"addr"`
}

type SortDistance []MetaDistance

func (a SortDistance) Len() int {
	return len(a)
}
func (a SortDistance) Less(i, j int) bool {
	return a[i].Distance < a[j].Distance
}
func (a SortDistance) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (h *Handler) Stores(w http.ResponseWriter, r *http.Request) {
	coord := r.URL.Query().Get("coord")
	productId := r.URL.Query().Get("product")
	if coord == "" || productId == "" {
		response.JSONError(w, http.StatusBadRequest, "Invalid params")
		return
	}

	coordinates := strings.Split(coord, ",")
	if len(coordinates) != 2 || coordinates[0] == "" || coordinates[1] == "" {
		response.JSONError(w, http.StatusBadRequest, "Wrong coordinates")
		return
	}
	lon := coordinates[0]
	lat := coordinates[1]

	latitude, _ := strconv.ParseFloat(lat, 64)
	longtitude, _ := strconv.ParseFloat(lon, 64)

	app := make(chan MetaDistance, 1000)
	for _, i := range h.Store.CachedShops {
		go func(app chan<- MetaDistance, shop store.Shop, productId string) {
			q, _ := h.Store.GetQuantity(shop.Id, productId)
			var dis float64
			isExist := q > 0
			if isExist {
				dis = utils.Distance(latitude, longtitude, shop.Lat, shop.Lon)
			}
			app <- MetaDistance{Exist: isExist, Distance: dis, Shop: shop}
		}(app, i, productId)
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	var points SortDistance
	go func(app <-chan MetaDistance) {
		var counter int = 1
		for v := range app {
			if v.Exist {
				points = append(points, v)
			}
			counter++
			if counter == len(h.Store.CachedShops) {
				break
			}
		}
		wg.Done()
	}(app)

	wg.Wait()

	sort.Sort(points)
	resultShop := points[0]

	urlTemplate := "https://static-maps.yandex.ru/1.x/?l=map&pt=%v,%v,pm2rdm~%v,%v,pm2ywl"
	url := fmt.Sprintf(urlTemplate, lat, lon, resultShop.Shop.Lat, resultShop.Shop.Lon)

	resp := &Response{}
	resp.Url = url
	resp.Name = resultShop.Shop.Name
	resp.Addr = resultShop.Shop.Addr
	responseBody, _ := json.Marshal(resp)
	response.WriteBody(w, http.StatusOK, responseBody)
}
