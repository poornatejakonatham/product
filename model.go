package main

import (
	"database/sql"
	"time"
)

type product struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Price     float64   `json:"price"`
	CreatedOn time.Time `json:"created_on"`
}

func getProducts(db *sql.DB, start, count int) ([]product, error) {
	rows, err := db.Query(
		"SELECT id, name, price, createdOn FROM product LIMIT $1 OFFSET $2",
		count,
		start,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	products := []product{}

	for rows.Next() {
		var p product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.CreatedOn); err != nil {
			return nil, err
		}

		products = append(products, p)
	}

	return products, nil
}

func (p *product) getProduct(db *sql.DB) error {
	return db.QueryRow(
		"SELECT name, price, createdOn FROM product WHERE id = $1",
		p.ID,
	).Scan(&p.Name, &p.Price, &p.CreatedOn)
}

func (p *product) createProduct(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO product(name, price, createdOn) VALUES ($1, $2, $3) RETURNING id",
		p.Name,
		p.Price,
		p.CreatedOn,
	).Scan(&p.ID)
	if err != nil {
		return err
	}

	return nil
}

func (p *product) updateProduct(db *sql.DB) error {
	_, err := db.Exec(
		"UPDATE product SET name = $1, price = $2 WHERE id = $3",
		p.Name,
		p.Price,
		p.ID,
	)

	return err
}

func (p *product) deleteProduct(db *sql.DB) error {
	_, err := db.Exec(
		"DELETE FROM product WHERE id = $1",
		p.ID,
	)

	return err
}
