package consumer

import (
	"database/sql"
	//"fmt"
	//"log"
	_ "mysql-master"
	s "strings"
)


type MassageOrder struct {
	Order struct {
		CustomerName           string  `json:"customer_name"`
		CustomerPhone          string  `json:"customer_phone"`
		CustId                 int     `json:"cust_id"`
		ItemID                 int     `json:"item_id"`
		ItemName               string  `json:"item_name"`
		ItemPrice              int     `json:"item_price"`
		OrderDetail            string  `json:"order_detail"`
		OrderFinishTime        string  `json:"order_finish_time"`
		OrderGenderPreferences int     `json:"order_gender_preferences"`
		OrderPaymentMethod     int  `json:"order_payment_method"`
		OrderPlace             string  `json:"order_place"`
		OrderPlaceMAP          string  `json:"order_place_MAP"`
		OrderPlaceLat          float64 `json:"order_place_lat"`
		OrderPlaceLng          float64 `json:"order_place_lng"`
		OrderPlaceNote         string  `json:"order_place_note"`
		OrderPrice             int     `json:"order_price"`
		OrderSpecialNote       string  `json:"order_special_note"`
		OrderStartTime         string  `json:"order_start_time"`
		OrderTime              int     `json:"order_time"`
		OrderType              int     `json:"order_type"`
		ServiceID              int     `json:"service_id"`
	} `json:"order"`
}

type ExtendMassageOrder struct {
	Order struct {
		ItemID                 int    `json:"item_id"`
		ItemName               string `json:"item_name"`
		ItemPrice              int    `json:"item_price"`
		OrderDetail            string `json:"order_detail"`
		OrderFinishTime        string `json:"order_finish_time"`
		OrderGenderPreferences int    `json:"order_gender_preferences"`
		OrderID                int    `json:"order_id"`
		OrderPaymentMethod     int    `json:"order_payment_method"`
		OrderPrice             int    `json:"order_price"`
		OrderSpecialNote       string `json:"order_special_note"`
		OrderStartTime         string `json:"order_start_time"`
		OrderTime              int    `json:"order_time"`
		OrderType              int    `json:"order_type"`
		ServiceID              int    `json:"service_id"`
	} `json:"order"`
}

type MassageOrderRespon struct {
	Error 		int		`json:"error"`
	Status      string	`json:"status"`
	App_version string	`json:"app_version"`
	Order_ID 	int64	`json:"order_id"`
		
	//Booking_id  int64
	//Masseur_id  int
}

type BidBlastOrderMassage struct{
	BlastType int   `json:"blast_type"`
	SpID      []int `json:"sp_id"`
	Message MessageBidBlastMassage
}

type MessageBidBlastMassage struct{
	Address        string `json:"address"`
	CustName       string `json:"cust_name"`
	CustPhone      string `json:"cust_phone"`
	EstimatedPrice int    `json:"estimated_price"`
	FinishTime     string `json:"finish_time"`
	OrderDetail    string `json:"order_detail"`
	OrderID        int    `json:"order_id"`
	PaymentMethod  int    `json:"payment_method"`
	SpecialNote    string `json:"special_note"`
	StartTime      string `json:"start_time"`
}

func AddOrder(db *sql.DB, sp *MassageOrder, mor *MassageOrderRespon) {
	
	rowsitem, err := db.Query(s.Join([]string{"SELECT item_price FROM golife_spapps.item ",
		"WHERE item_id = ? "}, " "), sp.Order.ItemID)
	//fmt.Printf("phone = %s", sp.Service_provider.Phone)
	if err != nil{
		s := err.Error()
		mor.Status = s
		mor.Error = 1
	}else{
		defer rowsitem.Close()
		
		item_price := 0
		
		if rowsitem.Next(){
			err := rowsitem.Scan(&item_price)
			
			if err != nil{
				s := err.Error()
				mor.Status = s
				mor.Error = 1
			}else{
				sp.Order.ItemPrice = item_price
			}
		}
	}
	
	rows, err := db.Query(s.Join([]string{"SELECT service_avg_price FROM golife_spapps.service ",
		"WHERE service_id = ? "}, " "), sp.Order.ServiceID)
	//fmt.Printf("phone = %s", sp.Service_provider.Phone)
	if err != nil{
		s := err.Error()
		mor.Status = s
		mor.Error = 1
	}else{
		defer rows.Close()
		
		service_price := 0
		
		if rows.Next(){
			err := rows.Scan(&service_price)
			
			if err != nil{
				s := err.Error()
				mor.Status = s
				mor.Error = 1
			}else{
				sp.Order.OrderPrice = service_price
			}
		}
	}
	
	tx, err := db.Begin()
	if err != nil {
		//log.Fatal(err)
		s := err.Error()
		mor.Status = s
		mor.Error = 1

	} else {

		defer tx.Rollback()

		stmt, err := tx.Prepare(s.Join([]string{"INSERT INTO golife_spapps.massage_order ",
			"(sp_id, ",
			"cust_id, ",
			"cust_name, ",
			"cust_phone, ",
			"service_id, ",
			"order_status,",
			"order_type, ",
			"order_detail, ",
			"order_item_id, ",
			"order_gender_preference, ",
			"order_lat, ",
			"order_lng, ",
			"order_address, ",
			"order_time, ",
			"order_start_time, ",
			"order_actual_start_time, ",
			"order_finish_time, ",
			"order_actual_finish_time, ",
			"work_otw_time, ",
			"order_total_price, ",
			"order_real_total_price, ",
			"order_special_note, ",
			"order_payment_method",
			") VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"}, " "))

		if err != nil {
			//log.Fatal(err)
			s := err.Error()
			mor.Status = s

		} else {

			//defer stmt.Close() // danger!
			/*
			   if sp.Voucher_code == "null" {
			       sp.vVoucher_code.Valid = false
			   }
			*/

			res, err := stmt.Exec(
				nil,//sp id masih null untuk pertama
				sp.Order.CustId,//cust id ?
				sp.Order.CustomerName,
				sp.Order.CustomerPhone,
				sp.Order.ServiceID,
				1,
				sp.Order.OrderType,//order type ?
				sp.Order.OrderDetail,
				sp.Order.ItemID,
				sp.Order.OrderGenderPreferences,
				sp.Order.OrderPlaceLat,
				sp.Order.OrderPlaceLng,
				s.Join([]string{sp.Order.OrderPlaceMAP, sp.Order.OrderPlace, sp.Order.OrderPlaceNote},","),
				sp.Order.OrderTime,
				sp.Order.OrderStartTime,
				nil,//order actul start time?
				sp.Order.OrderFinishTime,
				nil,//order actual finish time
				nil,//work otw time
				sp.Order.OrderPrice,//order total price dari front end 
				sp.Order.OrderPrice * (sp.Order.OrderTime/60) + sp.Order.ItemPrice,
				sp.Order.OrderSpecialNote,
				sp.Order.OrderPaymentMethod,
			)

			if err != nil {
				//log.Fatal(err)
				s := err.Error()
				mor.Status = s

			} else {

				err = tx.Commit()

				if err != nil {
					//log.Fatal(err)
					s := err.Error()
					mor.Status = s

				} else {

					bookId, _ := res.LastInsertId()
					//mor.Booking_id = bookId
					mor.Order_ID = bookId
					//mor.Masseur_id = 53
				}
			}
		}
		stmt.Close()
	}
}

func AddExtendOrder(db *sql.DB, sp *ExtendMassageOrder, mor *MassageOrderRespon) {
	
	rowsitem, err := db.Query(s.Join([]string{"SELECT item_price FROM golife_spapps.item ",
		"WHERE item_id = ? "}, " "), sp.Order.ItemID)
	//fmt.Printf("phone = %s", sp.Service_provider.Phone)
	if err != nil{
		s := err.Error()
		mor.Status = s
		mor.Error = 1
	}else{
		defer rowsitem.Close()
		
		item_price := 0
		
		if rowsitem.Next(){
			err := rowsitem.Scan(&item_price)
			
			if err != nil{
				s := err.Error()
				mor.Status = s
				mor.Error = 1
			}else{
				sp.Order.ItemPrice = item_price
			}
		}
	}
	
	rowprice, err := db.Query(s.Join([]string{"SELECT service_avg_price FROM golife_spapps.service ",
		"WHERE service_id = ? "}, " "), sp.Order.ServiceID)
	//fmt.Printf("phone = %s", sp.Service_provider.Phone)
	if err != nil{
		s := err.Error()
		mor.Status = s
		mor.Error = 1
	}else{
		defer rowprice.Close()
		
		service_price := 0
		
		if rowprice.Next(){
			err := rowprice.Scan(&service_price)
			
			if err != nil{
				s := err.Error()
				mor.Status = s
				mor.Error = 1
			}else{
				sp.Order.OrderPrice = service_price
			}
		}
	}
	
	sp_id := 0
	cust_id := 0
	cust_name := ""
	cust_phone := ""
	order_place_lat := 0.0000
	order_place_long := 0.0000
	order_address := ""
	order_act_start_time := "0000-00-00 00:00:00"
	order_act_finish_time := "0000-00-00 00:00:00"
	work_otw_time := "0000-00-00 00:00:00"
	
	rows, err:= db.Query(s.Join([]string{"SELECT sp_id, cust_id, cust_name,cust_phone, ",
	"order_lat, order_lng, order_address, order_start_time, ",
	"order_finish_time,work_otw_time FROM golife_spapps.massage_order ",
		"WHERE massage_order_id = ? "}, " "), sp.Order.OrderID)
	if err != nil{
		s := err.Error()
		mor.Status = s
		mor.Error = 1
	}else{
		defer rows.Close()
		
		if rows.Next(){
			err := rows.Scan(&sp_id, &cust_id, &cust_name, &cust_phone, 
			&order_place_lat, &order_place_long, &order_address, &order_act_start_time, 
			&order_act_finish_time, &work_otw_time)
			
			if err != nil{
				s := err.Error()
				mor.Status = s
				mor.Error = 1
			}
		}
	}
	
	tx, err := db.Begin()
	if err != nil {
		//log.Fatal(err)
		s := err.Error()
		mor.Status = s
		mor.Error = 1

	} else {

		defer tx.Rollback()

		stmt, err := tx.Prepare(s.Join([]string{"INSERT INTO golife_spapps.massage_order ",
			"(sp_id, ",
			"cust_id, ",
			"cust_name, ",
			"cust_phone, ",
			"service_id, ",
			"order_status,",
			"order_type, ",
			"order_detail, ",
			"order_item_id, ",
			"order_gender_preference, ",
			"order_lat, ",
			"order_lng, ",
			"order_address, ",
			"order_time, ",
			"order_start_time, ",
			"order_actual_start_time, ",
			"order_finish_time, ",
			"order_actual_finish_time, ",
			"work_otw_time, ",
			"order_total_price, ",
			"order_real_total_price, ",
			"order_special_note, ",
			"order_payment_method",
			") VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"}, " "))

		if err != nil {
			//log.Fatal(err)
			s := err.Error()
			mor.Status = s

		} else {

			//defer stmt.Close() // danger!
			/*
			   if sp.Voucher_code == "null" {
			       sp.vVoucher_code.Valid = false
			   }
			*/

			res, err := stmt.Exec(
				sp_id,//+current spid
				cust_id,//+current custid
				cust_name,//+ current customername
				cust_phone,// + current customerphone
				sp.Order.ServiceID,
				5,//+ current order status
				sp.Order.OrderType,//+ current order type
				sp.Order.OrderDetail,
				sp.Order.ItemID,
				sp.Order.OrderGenderPreferences,
				order_place_lat,//+ current lat
				order_place_long,//+curren long	
				order_address,//+ current orderaddress
				sp.Order.OrderTime,
				sp.Order.OrderStartTime,//+ request
				order_act_start_time,//+ current start time
				sp.Order.OrderFinishTime,//+ request
				order_act_finish_time,//+ current finish time
				work_otw_time,//+ current work otw time
				sp.Order.OrderPrice,//+ request ordertotalprice
				sp.Order.OrderPrice * (sp.Order.OrderTime / 60) + sp.Order.ItemPrice,
				sp.Order.OrderSpecialNote,
				sp.Order.OrderPaymentMethod,
			)

			if err != nil {
				//log.Fatal(err)
				s := err.Error()
				mor.Status = s

			} else {

				err = tx.Commit()

				if err != nil {
					//log.Fatal(err)
					s := err.Error()
					mor.Status = s

				} else {

					bookId, _ := res.LastInsertId()
					//mor.Booking_id = bookId
					mor.Order_ID = bookId
					//mor.Masseur_id = 53
				}
			}
		}
		stmt.Close()
	}
}
