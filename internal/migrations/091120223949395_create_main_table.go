package migrations

import (
	"database/sql"

	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up, Down)
}

func Up(tx *sql.Tx) error {
	_, err := tx.Exec(`
		CREATE TABLE IF NOT EXISTS stores (
			id          serial PRIMARY KEY,
			name        varchar(200) NOT NULL,
			addr        text NOT NULL,
			lon         varchar(100),
			lat         varchar(100)
		);`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS products (
			id          serial PRIMARY KEY,
			name        varchar(200) NOT NULL,
			price       integer
		);`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		CREATE TABLE IF NOT EXISTS quantity (
			store_id          integer,
			product_id        integer,
			quantity		  integer,
			FOREIGN KEY (store_id) REFERENCES stores (id),  
			FOREIGN KEY (product_id) REFERENCES products (id)
		);`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
			create index if not exists idx_store_id_product_id on quantity (store_id, product_id);
		`)
	if err != nil {
		return err
	}

	return nil
}

func Down(tx *sql.Tx) error {
	_, err := tx.Exec("drop table stores;")
	if err != nil {
		return err
	}
	_, err = tx.Exec("drop table products;")
	if err != nil {
		return err
	}
	_, err = tx.Exec("drop table quantity;")
	if err != nil {
		return err
	}
	return nil
}
