package store

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"strconv"
)

type Store struct {
	Db          *pgxpool.Pool
	CachedShops []Shop
}

type Shop struct {
	Id   int
	Name string
	Addr string
	Lon  float64
	Lat  float64
}

func New(db *pgxpool.Pool) *Store {
	s := &Store{Db: db}
	return s
}

func (s *Store) CacheShops() error {
	q := "SELECT id, lat, lon, name, addr FROM stores;"
	rows, err := s.Db.Query(context.Background(), q)
	if err != nil {
		return err
	}
	defer rows.Close()
	var shops []Shop
	for rows.Next() {
		shop := Shop{}
		var lat, lon string
		err := rows.Scan(&shop.Id, &lat, &lon, &shop.Name, &shop.Addr)
		if err != nil {
			panic(err)
		}
		latitude, _ := strconv.ParseFloat(lat, 64)
		longtitude, _ := strconv.ParseFloat(lon, 64)
		shop.Lat = latitude
		shop.Lon = longtitude

		shops = append(shops, shop)
	}
	s.CachedShops = shops

	return nil
}

func (s *Store) GetQuantity(storeId int, productId string) (int, error) {
	var quantity int
	q := "SELECT quantity FROM quantity where store_id = $1 AND product_id = $2;"
	err := s.Db.QueryRow(context.Background(), q, storeId, productId).Scan(&quantity)
	if err != nil {
		return 0, err
	}

	return quantity, nil
}
