package dao

import (
	"AliPayService/db"
	"AliPayService/entity"
	"github.com/jinzhu/gorm"
)

type Order entity.Order

func (o *Order) Create(dbLink *gorm.DB) error {
	create := dbLink.Create(o)
	return create.Error
}

func (o *Order) GetOneOrder() error {
	first := db.DbLink.Where(o).First(o)
	return first.Error
}

func (o *Order) GetOrderList() ([]Order, error) {
	var orderList []Order
	find := db.DbLink.Where(o).Find(&orderList)
	return orderList, find.Error
}

func (o Order) UpdateOrderStatus(dbLink *gorm.DB, orderId string) error {
	updates := dbLink.Model(&Order{}).Updates(o).Where("order_id = ?", orderId)
	return updates.Error
}
