package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/ProninIgorr/alcomarket-nearst/pkg/helpers/response"
	"github.com/ProninIgorr/alcomarket-nearst/pkg/helpers/utils"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

func (h *Handler) Slow(w http.ResponseWriter, r *http.Request) {

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
	var points SortDistance
	for _, shop := range h.Store.CachedShops {
		q, _ := h.Store.GetQuantity(shop.Id, productId)
		var dis float64
		isExist := q > 0
		if isExist {
			dis = utils.Distance(latitude, longtitude, shop.Lat, shop.Lon)
		}
		points = append(points, MetaDistance{Exist: isExist, Distance: dis, Shop: shop})
	}
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
