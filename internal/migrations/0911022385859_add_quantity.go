package migrations

import (
	"crypto/rand"
	"database/sql"
	"math/big"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(UpAddQuantity, DownAddQuantity)
}

func UpAddQuantity(tx *sql.Tx) error {

	var productsId []int
	pRows, err := tx.Query("SELECT id FROM products;")
	if err != nil {
		return err
	}
	defer pRows.Close()
	for pRows.Next() {
		var id int
		pRows.Scan(&id)
		productsId = append(productsId, id)
	}

	var storesId []int
	rows, err := tx.Query("SELECT id FROM stores;")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		rows.Scan(&id)
		storesId = append(storesId, id)
	}

	for _, i := range storesId {
		for _, k := range productsId {
			max := 3
			r, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
			_, err := tx.Exec("INSERT INTO quantity(store_id, product_id, quantity) VALUES ($1, $2, $3);", i, k, r.Int64())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func DownAddQuantity(tx *sql.Tx) error {
	_, err := tx.Exec("DELETE FROM quantity;")
	if err != nil {
		return err
	}
	return nil
}
