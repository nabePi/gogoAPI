package consumer

import (
	//"fmt"
	am "contler"
	//"fmt"
	"strconv"
	"database/sql"
	//"fmt"
	"log"
	_ "mysql-master"
	s "strings"
)

type AutoOrder struct {
	Order struct {
		CustomerName       string  `json:"customer_name"`
		CustomerPhone      string  `json:"customer_phone"`
		ItemIDs             string     `json:"item_ids"`
		ItemName           string  `json:"item_name"`
		ItemPrice          int     `json:"item_price"`
		OnDemand           bool    `json:"on_demand"`
		OrderDetail        string  `json:"order_detail"`
		OrderFinishTime    string  `json:"order_finish_time"`
		OrderPaymentMethod int     `json:"order_payment_method"`
		OrderPlace         string  `json:"order_place"`
		OrderPlaceMAP      string  `json:"order_place_MAP"`
		OrderPlaceLat      float64 `json:"order_place_lat"`
		OrderPlaceLng      float64 `json:"order_place_lng"`
		OrderPlaceNote     string  `json:"order_place_note"`
		OrderPrice         int     `json:"order_price"`
		OrderSpecialNote   string  `json:"order_special_note"`
		OrderStartTime     string  `json:"order_start_time"`
		OrderTime          int     `json:"order_time"`
		ServiceID          int     `json:"service_id"`
		ServiceIDs		   string  `json:"service_ids"`
	} `json:"order"`
}


type AutoOrderRespon struct {
	Error		int 	`json:"error"`
	Status      string	`json:"status"`
	App_version string	`json:"app_version"`
	Order_id 	int64	`json:"order_id"`
}

func AddAutoOrder(db *sql.DB, sp *AutoOrder, mor *AutoOrderRespon) {
	//serviceId := sp.Order.ServiceID
	serviceIds := s.Split(sp.Order.ServiceIDs, "#")
	itemId := s.Split(sp.Order.ItemIDs, "#")
	itemIds := s.Join(itemId, ",")

	
	servicesPrice := 0
	for _, v := range serviceIds {
		stmt, err := db.Prepare("SELECT service_price FROM service WHERE service_id=?")
	    if err != nil {
			log.Fatal(err)
			mor.Error = 1
			mor.Status = "Something not right"
			return
		}
	
		rows, err := stmt.Query(v)
	
		servicePrice := 0
	    if rows.Next() {
	        err = rows.Scan(&servicePrice)
		    if err != nil {
				log.Fatal(err)
				mor.Error = 1
				mor.Status = "An error occured"
				return
			} else {
				servicesPrice = servicesPrice + servicePrice	
			}
	    }		
	}
	
	
	rows, err := db.Query("SELECT item_price FROM item WHERE item_id IN ("+itemIds+")")
    if err != nil {
		log.Fatal(err)
		mor.Error = 1
		mor.Status = "Something not right"
		return
	}

	//rows, err := stmt.Query(itemIds)

    for rows.Next() {
		itemPrice := 0
		
        err = rows.Scan(&itemPrice)
	    if err != nil {
			log.Fatal(err)
			mor.Error = 1
			mor.Status = "An error occured"
			return
		} else {
			servicesPrice = servicesPrice + itemPrice
		}
    }
	
	
	tx, err := db.Begin()
	
	if err != nil {
		//log.Fatal(err)
		s := err.Error()
		mor.Status = s
		mor.Error = 1
	} else {
		mor.Error = 0
		defer tx.Rollback()

		stmt, err := tx.Prepare(s.Join([]string{"INSERT INTO golife_db.auto_order ",
			"(sp_id, ",
			"cust_id, ",
			"cust_name, ",
			"cust_phone, ",
			"service_id, ",
			"service_ids, ",
			"order_status, ",
			"order_detail, ",
			"order_item_ids, ",
			"order_item_code, ",
			"order_gender_preference, ",
			"order_lat, ",
			"order_lng, ",
			"order_address, ",
			"order_time, ",
			"order_start_time, ",
			"order_actual_start_time, ",
			"order_finish_time, ",
			"order_actual_finish_time, ",
			"order_real_total_price, ",			
			"order_total_price, ",
			"order_special_note, ",      
			"order_payment_method, ",
			"order_on_demand) ",
			"VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"}, " "))

		if err != nil {
			//log.Fatal(err)
			s := err.Error()
			mor.Status = s
			mor.Error = 1

		} else {

			//defer stmt.Close() // danger!
			/*
			if sp.Voucher_code == "null" {
				sp.vVoucher_code.Valid = false
			}
			*/

			res, err := stmt.Exec(
				nil,
				nil,
				sp.Order.CustomerName,
				sp.Order.CustomerPhone,
				sp.Order.ServiceID,
				sp.Order.ServiceIDs,
				ORDER_STATUS_RECEIVED,
				sp.Order.OrderDetail,
				sp.Order.ItemIDs,
				nil,
				nil,
				sp.Order.OrderPlaceLat,
				sp.Order.OrderPlaceLng,
				s.Join([]string{sp.Order.OrderPlaceMAP, sp.Order.OrderPlace}, " - "),
				sp.Order.OrderTime,
				sp.Order.OrderStartTime,
				nil,
				sp.Order.OrderFinishTime,
				nil,
				servicesPrice,				
				sp.Order.OrderPrice,
				sp.Order.OrderSpecialNote,      
				sp.Order.OrderPaymentMethod,
				sp.Order.OnDemand,
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

					orderId, _ := res.LastInsertId()
					mor.Order_id = orderId
					
					// query untuk mendapat service provider yang berjarak <= 25 Km
					mystr := s.Join([]string{`
						SELECT sp.sp_id,
							(
								(
									(
										acos(
											sin((go.order_lat * pi() / 180)) * sin((spl.latitude * pi() / 180)) + cos((go.order_lat * pi() / 180)) * cos((spl.latitude * pi() / 180)) * cos(
												(
													(go.order_lng - spl.longitude) * pi() / 180
												)
											)
										)
									) * 180 / pi()
								) * 60 * 1.1515
							) AS distance
						
						FROM
							sp_location spl 
							JOIN 
							service_provider sp ON sp.sp_id = spl.sp_id
							,auto_order go

						WHERE
							go.auto_order_id = `, strconv.FormatInt(orderId, 10),
							`AND sp.sp_status = 1`,
							`
						HAVING
							distance <= 25
					`}, " ")
					
					if sp.Order.OrderPlaceLat == 0 || sp.Order.OrderPlaceLng == 0 {
					mystr = s.Join([]string{`
						SELECT sp.sp_id,
							(
								(
									(
										acos(
											sin((go.order_lat * pi() / 180)) * sin((spl.latitude * pi() / 180)) + cos((go.order_lat * pi() / 180)) * cos((spl.latitude * pi() / 180)) * cos(
												(
													(go.order_lng - spl.longitude) * pi() / 180
												)
											)
										)
									) * 180 / pi()
								) * 60 * 1.1515
							) AS distance
						
						FROM
							sp_location spl 
							JOIN 
							service_provider sp ON sp.sp_id = spl.sp_id
							,auto_order go

						WHERE
							go.auto_order_id = `, strconv.FormatInt(orderId, 10),
							`AND sp.sp_status = 1`,
						}, " ")
												
					}
					
		            rows, err := db.Query(mystr)
					
					if err != nil {
						mor.Error = 1
						mor.Status = err.Error()
					}
					
					var nearestSP []int
					
					for rows.Next() {
						var spId int
						var distance float64
						rows.Scan(&spId, &distance)
						
						nearestSP = append(nearestSP, spId)	
					}
									
					type blastMessage struct {
						Order_id	int64 `json:"order_id"`
						Vertical_id	int	`json:"vertical_id"`
						Cust_name	string `json:"cust_name"`
						Start_time	string `json:"start_time"`
						Finish_time string `json:"finish_time"`
						Cust_phone	string `json:"cust_phone"`
						Address 	string `json:"address"`
						Order_detail string `json:"order_detail"`
						Special_note string `json:"special_note"`	
						Estimated_price int `json:"estimated_price"`
						Payment_method  int `json:"payment_method"`
						On_Demand		bool `json:"on_demand"`
					}	
				
					type AutoOrderBlastMessage struct {
						Sp_id			[]int `json:"sp_id"'`
						Blast_type		int	`json:"blast_type"`
						Message_type 	string `json:"message_type"`
						Message blastMessage
					}
					
					autoOrderBlastMessage := AutoOrderBlastMessage {
						nearestSP,
						BLAST_TYPE_ADD_ORDER,
						"addAutoOrder",
						blastMessage {
							orderId,
        					GO_AUTO_VERTICAL_ID,
							sp.Order.CustomerName,
							sp.Order.OrderStartTime,
							sp.Order.OrderFinishTime,
							sp.Order.CustomerPhone,
							s.Join([]string{sp.Order.OrderPlaceMAP, sp.Order.OrderPlace}, " - "),
							sp.Order.OrderDetail,
							sp.Order.OrderSpecialNote,
							servicesPrice,
							sp.Order.OrderPaymentMethod,
							sp.Order.OnDemand,
						},
					}	
				
					go am.PublishRoutine(am.PubnubChannel, autoOrderBlastMessage)	
					
					UpdateAutoOrderStatus(db, orderId, ORDER_STATUS_PUBLISHED)				
				}
			}
		}
		stmt.Close()
	}
}

func UpdateAutoOrderStatus(db *sql.DB, orderId int64, orderStatus int) {

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)

	} else {

		defer tx.Rollback()

		stmt, err := tx.Prepare(`UPDATE auto_order SET order_status = ?
				WHERE auto_order_id=?;`)

		if err != nil {
			log.Fatal(err)
		} else {

			//defer stmt.Close() // danger!
			/*if glo.Voucher_code == "null" {
				glo.vVoucher_code.Valid = false
			}*/

			_, err := stmt.Exec(orderStatus, orderId)

			if err != nil {
				log.Fatal(err)

			} else {

				err = tx.Commit()

				if err != nil {
					log.Fatal(err)

				}
			}
		}
		stmt.Close()
	}
}

