package model

type Product struct {
	Id     int     `json:"id" gorm:"primaryKey"`
	UserId int     `json:"userId"`
	Name   string  `json:"name"`
	Price  float64 `json:"price"`
}
