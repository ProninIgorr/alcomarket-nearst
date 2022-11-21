package migrations

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/ProninIgorr/alcomarket-nearst/pkg/helpers/httphelper"
	"github.com/pressly/goose"
	"net/http"
	"time"
)

func init() {
	goose.AddMigration(upAddStores, downAddStores)
}

type Response struct {
	Type       string `json:"type"`
	Properties struct {
		ResponseMetaData struct {
			SearchResponse struct {
				Found     int         `json:"found"`
				Display   string      `json:"display"`
				BoundedBy [][]float64 `json:"boundedBy"`
			} `json:"SearchResponse"`
			SearchRequest struct {
				Request   string      `json:"request"`
				Skip      int         `json:"skip"`
				Results   int         `json:"results"`
				BoundedBy [][]float64 `json:"boundedBy"`
			} `json:"SearchRequest"`
		} `json:"ResponseMetaData"`
	} `json:"properties"`
	Features []Feature `json:"features"`
}

type Feature struct {
	Type     string `json:"type"`
	Geometry struct {
		Type        string    `json:"type"`
		Coordinates []float64 `json:"coordinates"`
	} `json:"geometry"`
	Properties struct {
		Name            string      `json:"name"`
		Description     string      `json:"description"`
		BoundedBy       [][]float64 `json:"boundedBy"`
		CompanyMetaData struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Address string `json:"address"`
			URL     string `json:"url"`
			Phones  []struct {
				Type      string `json:"type"`
				Formatted string `json:"formatted"`
			} `json:"Phones"`
			Categories []struct {
				Class string `json:"class"`
				Name  string `json:"name"`
			} `json:"Categories"`
			Hours struct {
				Text           string `json:"text"`
				Availabilities []struct {
					Intervals []struct {
						From string `json:"from"`
						To   string `json:"to"`
					} `json:"Intervals,omitempty"`
					Monday          bool `json:"Monday,omitempty"`
					TwentyFourHours bool `json:"TwentyFourHours,omitempty"`
					Tuesday         bool `json:"Tuesday,omitempty"`
					Wednesday       bool `json:"Wednesday,omitempty"`
					Thursday        bool `json:"Thursday,omitempty"`
					Friday          bool `json:"Friday,omitempty"`
					Saturday        bool `json:"Saturday,omitempty"`
					Sunday          bool `json:"Sunday,omitempty"`
				} `json:"Availabilities"`
			} `json:"Hours"`
		} `json:"CompanyMetaData"`
	} `json:"properties"`
}

func upAddStores(tx *sql.Tx) error {

	f, err := getStores()
	if err != nil {
		return err
	}

	for _, i := range f {
		_, err := tx.Exec(`INSERT INTO stores(name, addr, lon, lat) VALUES ($1, $2, $3, $4);`,
			i.Properties.Name,
			i.Properties.CompanyMetaData.Address,
			i.Geometry.Coordinates[1],
			i.Geometry.Coordinates[0],
		)
		if err != nil {
			return err
		}
	}

	products := [][]string{[]string{"Red Wine", "1700"}, []string{"White Wine", "1300"}, []string{"Beer", "400"}, []string{"Vodka", "700"}}
	for _, p := range products {
		_, err := tx.Exec(`INSERT INTO products(name, price) VALUES ($1, $2);`,
			p[0],
			p[1],
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func downAddStores(tx *sql.Tx) error {
	_, err := tx.Exec(`DELETE FROM stores;`)
	if err != nil {
		return err
	}
	return nil
}

func getStores() ([]Feature, error) {
	var result []Feature
	c := httphelper.Client{Timeout: time.Second * 5}
	urlTemplate := "https://search-maps.yandex.ru/v1/?apikey=<API_KEY>&text=%s&lang=ru_RU&results=500&skip=%d&ll=37.618920,55.756994&spn=0.552069,0.400552"

	for i := 0; i < 1500; i += 500 {
		url := fmt.Sprintf(urlTemplate, "Alcomarket", i)
		body, err := c.MakeRequest(context.Background(), http.MethodGet, url, []byte{}, map[string]string{})
		if err != nil {
			return nil, err
		}
		r := &Response{}
		err = json.Unmarshal(body, r)
		if err != nil {
			return nil, err
		}

		result = append(result, r.Features...)
	}

	return result, nil
}
