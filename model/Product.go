package model

import (
	"github.com/jinzhu/gorm"
)

//Product ...
// This struct is used to model the Product Table in the PostGRESQL database
type Product struct {
	gorm.Model
	Name  string  `json:"name"`
	Price float64 `json:"price,string"`
}

// GetProduct ...
// Parameters: db *sql.DB
// Description: This function is passed the database which it then
// 	uses to get the rest of the Product Data from the database.
func (p *Product) GetProduct() error {
	err := GetDB().Where("id = ?", p.ID).First(p).Error
	return err
}

// UpdateProduct ...
// Parameters: db *sql.DB
// Description: This function is passed the database which it then
//	uses to update the product based upon the current values of
//	Name and Price.
func (p *Product) UpdateProduct() error {
	err := GetDB().Model(p).Update(map[string]interface{}{"name": p.Name, "price": p.Price}).Error
	return err
}

// DeleteProduct ...
// Parameters: db. *sql.DB
// Description: This function is passed the database which it then
//	uses to delete the product by the product ID/
func (p *Product) DeleteProduct() error {
	err := GetDB().Table("products").Delete(p).Error
	return err
}

func (p *Product) HardDeleteProduct() error {
	err := GetDB().Table("products").Unscoped().Delete(p).Error
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
func (p *Product) CreateProduct() error {
	err := GetDB().Create(p).Error
	if err != nil {
		return err
	}
	return nil
}

//GetDeletedProducts ... This function retrieves soft deleted items.
func GetDeletedProducts(start, count int, p *[]Product) error {
	err := GetDB().Unscoped().Offset(start).Limit(count).Find(p).Error

	if err != nil {
		return err
	}

	return nil
}

// GetProducts ...
// Parameters: db *sql.DB
// Description: This function is passed the Database, a start value, and a count value.
//	This
func GetProducts(start, count int, p *[]Product) error {

	err := GetDB().Offset(start).Limit(count).Find(p).Error

	if err != nil {
		return err
	}

	return nil
}
