package model

import (
	"database/sql"

	"github.com/jinzhu/gorm"
)

//Product ...
// This struct is used to model the Product Table in the PostGRESQL database
type Product struct {
	gorm.Model
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// GetProduct ...
// Parameters: db *sql.DB
// Description: This function is passed the database which it then
// 	uses to get the rest of the Product Data from the database.
func (p *Product) GetProduct(db *sql.DB) error {
	return db.QueryRow("SELECT name, price FROM products WHERE id=$1",
		p.ID).Scan(&p.Name, &p.Price)
}

// UpdateProduct ...
// Parameters: db *sql.DB
// Description: This function is passed the database which it then
//	uses to update the product based upon the current values of
//	Name and Price.
func (p *Product) UpdateProduct(db *sql.DB) error {
	_, err :=
		db.Exec("UPDATE products SET name=$1, price=$2 WHERE id=$3",
			p.Name, p.Price, p.ID)

	return err
}

// DeleteProduct ...
// Parameters: db. *sql.DB
// Description: This function is passed the database which it then
//	uses to delete the product by the product ID/
func (p *Product) DeleteProduct(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM products WHERE id=$1", p.ID)

	return err
}

// CreateProduct ...
// Parameters: db *sql.DB
// Description: This function is passed the database which it then
//	uses to create the product based upon the values of name and price.
//	The database returns the newly created id which the Database struct
//	function .Scan can use be passed the pointer to the Product ID
//	value and then update that value to the newly created Product ID from
//	the database.
func (p *Product) CreateProduct(db *sql.DB) error {
	err := db.QueryRow(
		"INSERT INTO products(name, price) VALUES($1, $2) RETURNING id",
		p.Name, p.Price).Scan(&p.ID)

	if err != nil {
		return err
	}

	return nil
}

// GetProducts ...
// Parameters: db *sql.DB
// Description: This function is passed the Database, a start value, and a count value.
//	This
func GetProducts(db *sql.DB, start, count int) ([]Product, error) {
	rows, err := db.Query(
		"SELECT id, name,  price FROM products LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	products := []Product{}

	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}
