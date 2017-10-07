package consumer

import (
	"database/sql"
	//	"encoding/json"
	"fmt"
	"log"
	_ "mysql-master"
	//s "strings"
)

type Order struct {
	Status             string
	Order_Id           int64
	Sp_Id              uint
	Cust_Id            uint
	Cust_name          string
	Cust_phone         string
	Cust_pubnub_id     string
	Service_id         uint
	Order_status       int
	Order_detail       string
	Item_id            int
	Item_code          string //go-auto
	Item               string
	Item_color         string
	Gender_pref        int
	Latitude           float32
	Longitude          float32
	Address            string
	Place_note         string
	City               string
	Order_time         int
	Start_time         string
	Act_start_time     string
	Finish_time        string
	Act_finish_time    string
	Additional_service string
	Total_price        float32
	Estimated_price    float32
	Special_note       string
	Payment_method     int
	Vertical_id        int
}

type OrderStatusResponse struct {
	Status       string
	App_version  string
	Order_id     int64
	Order_status int
}

func GetOrderById(db *sql.DB, o *Order) {
	fmt.Println("GET ORDER")
	var query string
	switch o.Vertical_id {
	case 1:
		query = `SELECT glam_order_id, cust_id, cust_name, cust_phone, cust_pubnub_id, service_id,
		order_status, order_detail,order_item,order_item_color,
		order_gender_preference, order_lat, order_lng,order_address, order_place_note,
		order_city, order_time,
		order_start_time, order_finish_time, additional_service, order_estimated_price, order_special_note,order_payment_method
		FROM glam_order WHERE glam_order_id = ?`
		if err := db.QueryRow(query, o.Order_Id).Scan(&o.Order_Id, &o.Cust_Id, &o.Cust_name, &o.Cust_phone,
			&o.Cust_pubnub_id, &o.Service_id, &o.Order_status, &o.Order_detail,
			&o.Item, &o.Item_color, &o.Gender_pref, &o.Latitude, &o.Longitude,
			&o.Address, &o.Place_note, &o.City, &o.Order_time, &o.Start_time, &o.Finish_time, &o.Additional_service, &o.Estimated_price,
			&o.Special_note, &o.Payment_method); err == nil {
			o.Status = "success"
		} else if err == sql.ErrNoRows {
			o.Status = "No Record Found"
		} else {
			log.Fatal(err)
		}
	case 2:
		query = `SELECT massage_order_id, cust_id, cust_name, cust_phone, cust_pubnub_id, service_id,
		order_status, order_detail, order_item_id,
		order_gender_preference, order_lat, order_lng,order_address, order_time,
		order_start_time, order_finish_time, order_total_price, order_special_note,order_payment_method
		FROM massage_order WHERE massage_order_id = ?`
		if err := db.QueryRow(query, o.Order_Id).Scan(&o.Order_Id, &o.Cust_Id, &o.Cust_name, &o.Cust_phone,
			&o.Cust_pubnub_id, &o.Service_id, &o.Order_status, &o.Order_detail, &o.Item_id,
			&o.Gender_pref, &o.Latitude, &o.Longitude,
			&o.Address, &o.Order_time, &o.Start_time, &o.Finish_time, &o.Total_price,
			&o.Special_note, &o.Payment_method); err == nil {
			o.Status = "success"
		} else if err == sql.ErrNoRows {
			o.Status = "No Record Found"
		} else {
			log.Fatal(err)
		}
	case 3:
		query = `SELECT clean_order_id, cust_id, cust_name, cust_phone, cust_pubnub_id, service_id,
		order_status, order_detail, order_item_id,
		order_gender_preference, order_lat, order_lng,order_address, order_time,
		order_start_time, order_finish_time, order_total_price, order_special_note,order_payment_method
		FROM clean_order WHERE clean_order_id = ?`
		if err := db.QueryRow(query, o.Order_Id).Scan(&o.Order_Id, &o.Cust_Id, &o.Cust_name, &o.Cust_phone,
			&o.Cust_pubnub_id, &o.Service_id, &o.Order_status, &o.Order_detail, &o.Item_id,
			&o.Gender_pref, &o.Latitude, &o.Longitude,
			&o.Address, &o.Order_time, &o.Start_time, &o.Finish_time, &o.Total_price,
			&o.Special_note, &o.Payment_method); err == nil {
			o.Status = "success"
		} else if err == sql.ErrNoRows {
			o.Status = "No Record Found"
		} else {
			log.Fatal(err)
		}
	case 4:
		query = `SELECT auto_order_id, cust_id, cust_name, cust_phone, cust_pubnub_id, service_id,
		order_status, order_detail,order_item_id,order_item_code,
		order_gender_preference, order_lat, order_lng,order_address, order_time,
		order_start_time, order_finish_time, order_total_price, order_special_note,order_payment_method
		FROM auto_order WHERE auto_order_id = ?`
		if err := db.QueryRow(query, o.Order_Id).Scan(&o.Order_Id, &o.Cust_Id, &o.Cust_name, &o.Cust_phone,
			&o.Cust_pubnub_id, &o.Service_id, &o.Order_status, &o.Order_detail, &o.Item_id, &o.Item_code,
			&o.Gender_pref, &o.Latitude, &o.Longitude,
			&o.Address, &o.Order_time, &o.Start_time, &o.Finish_time, &o.Total_price,
			&o.Special_note, &o.Payment_method); err == nil {
			o.Status = "success"
		} else if err == sql.ErrNoRows {
			o.Status = "No Record Found"
		} else {
			log.Fatal(err)
		}
	}

}

func GetStatusOrderById(db *sql.DB, o *Order, osr *OrderStatusResponse) {
	fmt.Println("GET STATUS")
	var query string
	switch o.Vertical_id {
	case 1:
		query = "SELECT order_status FROM glam_order WHERE glam_order_id = ?"
	case 2:
		query = "SELECT order_status FROM massage_order WHERE massage_order_id = ?"
	case 3:
		query = "SELECT order_status FROM clean_order WHERE clean_order_id = ?"
	case 4:
		query = "SELECT order_status FROM auto_order WHERE auto_order_id = ?"
	}

	orderStatus := 0

	if err := db.QueryRow(query, o.Order_Id).Scan(&orderStatus); err == nil {
		osr.Order_status = orderStatus
		osr.Order_id = o.Order_Id
	} else if err == sql.ErrNoRows {
		osr.Status = "No Record Found"
	} else {
		log.Fatal(err)
	}

}
