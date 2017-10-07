package consumer

import (
	am "contler"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	_ "mysql-master"
	s "strings"
)

type GlamOrder struct {
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

type GlamOrderRespon struct {
	Status       string
	App_version  string
	Order_id     int64
	Order_status int
}

type Service struct {
	Service_id          int
	Service_level       int8
	Service_parent      int
	Service_name        string
	Service_description string
	Service_price       float32
	Service_comodity    int8
	Vertical_id         int
}

type ServiceProvider struct {
	Id        int
	Name      string
	Gender    int
	Phone     string
	Pubnub_id string
	Rate      float64
	Foto_Url  string
	Distance  float64
}

type BidBlast struct {
	Sp_id           []int
	Order_id        int64
	Cust_name       string
	Start_time      string
	Finish_time     string
	Cust_phone      string
	Address         string
	Order_detail    string
	Special_note    string
	Estimated_price float32
	Payment_method  int
}

/*func GetOrderById(db *sql.DB, glo *GlamOrder) {
	fmt.Println("GET ORDER")
	if glo.Vertical_id == 1 {
		query := `SELECT glam_order_id, cust_id, cust_name, cust_phone, cust_pubnub_id, service_id,
		order_status, order_detail,order_item,order_item_color,
		order_gender_preference, order_lat, order_lng,order_address, order_place_note,
		order_city, order_time,
		order_start_time, order_finish_time, order_estimated_price, order_special_note,order_payment_method
		FROM glam_order WHERE glam_order_id = ?`

		if err := db.QueryRow(query, glo.Order_Id).Scan(&glo.Order_Id, &glo.Cust_Id, &glo.Cust_name, &glo.Cust_phone,
			&glo.Cust_pubnub_id, &glo.Service_id, &glo.Order_status, &glo.Order_detail,
			&glo.Item, &glo.Item_color, &glo.Gender_pref, &glo.Latitude, &glo.Longitude,
			&glo.Address, &glo.Place_note, &glo.City, &glo.Order_time, &glo.Start_time, &glo.Finish_time, &glo.Estimated_price,
			&glo.Special_note, &glo.Payment_method); err == nil {
			glo.Status = "success"
		} else if err == sql.ErrNoRows {
			glo.Status = "No Record Found"
		} else {
			log.Fatal(err)
		}

	} else {
		//glor.Status = "Wrong vertical id"
	}

}*/

/*
func GetStatusOrderById(db *sql.DB, glo *GlamOrder, glor *GlamOrderRespon) {
	fmt.Println("GET STATUS")
	if glo.Vertical_id == 1 {
		query := "SELECT order_status FROM glam_order WHERE glam_order_id = ?"

		orderStatus := 0

		if err := db.QueryRow(query, glo.Order_Id).Scan(&orderStatus); err == nil {
			glor.Order_status = orderStatus
			glor.Order_id = glo.Order_Id
		} else if err == sql.ErrNoRows {
			glor.Status = "No Record Found"
		} else {
			log.Fatal(err)
		}

	} else {
		glor.Status = "Wrong vertical id"
	}

}*/

func AddGlamOrder(db *sql.DB, glo *GlamOrder, glor *GlamOrderRespon) {

	tx, err := db.Begin()
	if err != nil {
		//log.Fatal(err)
		s := err.Error()
		glor.Status = s

	} else {

		defer tx.Rollback()

		stmt, err := tx.Prepare(`INSERT INTO golive_db.glam_order ( 
		cust_id, cust_name, cust_phone, cust_pubnub_id, service_id,
		order_status, order_detail,order_item_id, order_item,order_item_color, 
		order_gender_preference, order_lat, order_lng,order_address, order_place_note,
		order_city, order_time, 
		order_start_time, order_finish_time, additional_service, order_estimated_price, order_special_note,order_payment_method) 
		VALUES (?, ?, ?, ?, ?,
		 ?, ?, ?, ?, ?,
		 ?, ?, ?, ?, ?,
		 ?, ?, ?, ?, ?, ?, ?, ?)`)

		if err != nil {
			log.Fatal(err)
			s := err.Error()
			glor.Status = s
		} else {

			//defer stmt.Close() // danger!
			/*if glo.Voucher_code == "null" {
				glo.vVoucher_code.Valid = false
			}*/

			res, err := stmt.Exec(glo.Cust_Id, glo.Cust_name, glo.Cust_phone,
				glo.Cust_pubnub_id, glo.Service_id, 1, glo.Order_detail, 0,
				glo.Item, glo.Item_color, glo.Gender_pref, glo.Latitude, glo.Longitude,
				glo.Address, glo.Place_note, glo.City, glo.Order_time, glo.Start_time, glo.Finish_time, glo.Additional_service, glo.Estimated_price,
				glo.Special_note, glo.Payment_method)

			if err != nil {
				log.Fatal(err)
				s := err.Error()
				glor.Status = s

			} else {

				err = tx.Commit()

				if err != nil {
					//log.Fatal(err)
					s := err.Error()
					glor.Status = s

				} else {

					bookId, _ := res.LastInsertId()
					glo.Order_Id = bookId
					glor.Order_id = bookId
					getServiceProviders(db, glo, glor)
					calculateService(db, glo, glor)

				}
			}
		}
		stmt.Close()
	}
}

func getServiceProviders(db *sql.DB, glo *GlamOrder, glor *GlamOrderRespon) {

	query := fmt.Sprintf(`SELECT sp.sp_id, sp.sp_name, sp.sp_gender, sp.sp_phone, sp.sp_pubnub_id, sp.sp_rate,
sp.sp_foto_url, (((acos(sin((` + fmt.Sprintf("%.6f", glo.Latitude) + `*pi()/180)) * 
            sin((sl.latitude*pi()/180))+cos((` + fmt.Sprintf("%.6f", glo.Latitude) + `*pi()/180)) * 
            cos((sl.latitude*pi()/180)) * cos(((` + fmt.Sprintf("%.6f", glo.Longitude) + `- sl.longitude)* 
            pi()/180))))*180/pi())*60*1.1515
        ) as distance 
FROM service_provider sp
INNER JOIN sp_service ss
	ON sp.sp_id = ss.sp_id JOIN sp_location sl
	ON sl.sp_id = sp.sp_id 
WHERE ss.service_id = ? AND sp.sp_status != 0
  HAVING distance <= 25`)

	stmt, _ := db.Prepare(query)
	rows, err := stmt.Query(glo.Service_id)

	var avai_sp []ServiceProvider
	if err != nil {
		log.Fatal(err)
	} else {
		defer rows.Close()

		sp := make([]ServiceProvider, 0, 20)
		ln := 1

		for rows.Next() {
			sp = sp[0:ln] // extend length by 1

			err := rows.Scan(&sp[ln-1].Id, &sp[ln-1].Name, &sp[ln-1].Gender, &sp[ln-1].Phone,
				&sp[ln-1].Pubnub_id, &sp[ln-1].Rate, &sp[ln-1].Foto_Url, &sp[ln-1].Distance)

			if err != nil {

				log.Fatal(err)
				break

			} else {

				ln++
			}
		}

		err = rows.Err()
		if err != nil {
			//			s := err.Error()
			log.Fatal(err)
		} else {

			query = fmt.Sprintf(`SELECT go.sp_id
					FROM glam_order go
					WHERE go.sp_id IS NOT NULL AND go.service_id =? AND order_finish_time BETWEEN  ? AND ?`)

			stmt, _ = db.Prepare(query)
			rows, err = stmt.Query(glo.Service_id, glo.Start_time, glo.Finish_time)
			if err != nil {
				log.Fatal(err)
			} else {
				defer rows.Close()
				order_sp := make([]ServiceProvider, 0, 20)
				ln = 1

				for rows.Next() {
					order_sp = order_sp[0:ln] // extend length by 1

					err := rows.Scan(&order_sp[ln-1].Id)

					if err != nil {
						//				s := err.Error()
						log.Fatal(err)
						break

					} else {

						ln++
					}
				}

				err = rows.Err()
				if err != nil {
					//			s := err.Error()
					log.Fatal(err)
				} else {

					for _, s := range sp {
						found := false
						for _, os := range order_sp {
							if s.Id == os.Id {
								found = true
								break
							}
						}
						if !found {

							avai_sp = append(avai_sp, s)
						}
					}

				}
			}

			var sp_ids []int
			// send sp list ke pubnub
			for _, element := range avai_sp {
				sp_ids = append(sp_ids, element.Id)
			}

			blast := &BidBlast{
				Sp_id:           sp_ids,
				Order_id:        glo.Order_Id,
				Cust_name:       glo.Cust_name,
				Start_time:      glo.Start_time,
				Finish_time:     glo.Finish_time,
				Cust_phone:      glo.Cust_phone,
				Address:         glo.Address,
				Order_detail:    glo.Order_detail,
				Special_note:    glo.Special_note,
				Estimated_price: glo.Estimated_price,
				Payment_method:  glo.Payment_method}

			blastJsn, _ := json.Marshal(blast)

			go am.PublishRoutine("golife.wgs.channel", string(blastJsn))

			UpdateOrderStatus(db, glo, 2)
			//for j := 0; j < len(sp.Service); j++ {
			//	fmt.Printf("v=%+v\n", sp.Service[j])
			//}
		}
	}
}

func calculateService(db *sql.DB, glo *GlamOrder, glor *GlamOrderRespon) {

	addService := s.Split(glo.Additional_service, "-")
	args := make([]interface{}, len(addService))
	for index, element := range addService {
		args[index] = element
	}
	query := fmt.Sprintf("SELECT *	FROM service AS s WHERE s.service_id IN (?" + s.Repeat(",?", len(addService)-1) + ")")

	stmt, _ := db.Prepare(query)
	rows, err := stmt.Query(args...)
	//rows, err := db.Query("SELECT *	FROM service AS s WHERE s.service_id IN (?" + s.Repeat(",?", len(addService)-1) + ")")

	fmt.Println("SELECT *	FROM service AS s WHERE s.service_id IN (?" + s.Repeat(",?", len(addService)-1) + ")")
	if err != nil {
		log.Fatal(err)
	} else {
		defer rows.Close()

		srv := make([]Service, 0, 20)
		ln := 1

		for rows.Next() {
			srv = srv[0:ln] // extend length by 1

			err := rows.Scan(&srv[ln-1].Service_id, &srv[ln-1].Service_level, &srv[ln-1].Service_parent, &srv[ln-1].Service_name,
				&srv[ln-1].Service_description, &srv[ln-1].Service_price, &srv[ln-1].Service_comodity, &srv[ln-1].Vertical_id)

			if err != nil {
				//				s := err.Error()
				log.Fatal(err)
				break

			} else {

				//LoadDuration(db, sp, &srv[ln-1])
				ln++
			}
		}

		err = rows.Err()
		if err != nil {
			//			s := err.Error()
			log.Fatal(err)
		} else {
			//TODO itungan price non komodity
			totalPrice := float32(0.0)
			for _, element := range srv {
				totalPrice += element.Service_price
				fmt.Println("ELEMENT ", element)
			}
			UpdateTotalPrice(db, glo, totalPrice)
			fmt.Println("TOTAL PRICE ", totalPrice)
			//for j := 0; j < len(sp.Service); j++ {
			//	fmt.Printf("v=%+v\n", sp.Service[j])
			//}
		}
	}
}

func UpdateTotalPrice(db *sql.DB, glo *GlamOrder, totalPrice float32) {

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)

	} else {

		defer tx.Rollback()

		stmt, err := tx.Prepare(`UPDATE glam_order SET order_total_price =?
				WHERE glam_order_id=?;`)

		if err != nil {
			log.Fatal(err)
		} else {

			//defer stmt.Close() // danger!
			/*if glo.Voucher_code == "null" {
				glo.vVoucher_code.Valid = false
			}*/

			_, err := stmt.Exec(totalPrice, glo.Order_Id)

			if err != nil {
				log.Fatal(err)

			} else {

				err = tx.Commit()

				if err != nil {
					log.Fatal(err)

				} else {

					glo.Total_price = totalPrice

				}
			}
		}
		stmt.Close()
	}
}

func UpdateOrderStatus(db *sql.DB, glo *GlamOrder, orderStatus int) {

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)

	} else {

		defer tx.Rollback()

		stmt, err := tx.Prepare(`UPDATE glam_order SET order_status =?
				WHERE glam_order_id=?;`)

		if err != nil {
			log.Fatal(err)
		} else {

			//defer stmt.Close() // danger!
			/*if glo.Voucher_code == "null" {
				glo.vVoucher_code.Valid = false
			}*/

			_, err := stmt.Exec(orderStatus, glo.Order_Id)

			if err != nil {
				log.Fatal(err)

			} else {

				err = tx.Commit()

				if err != nil {
					log.Fatal(err)

				} else {

					glo.Order_status = orderStatus

				}
			}
		}
		stmt.Close()
	}
}
