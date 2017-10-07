package serviprov

import (
	"database/sql"
	"fmt"
	_ "mysql-master"
	//	s "strings"
	"log"
)

type BidOrderParam struct {
	Sp_id       int
	Service_id  int
	Order_id    int
	Vertical_id int
}

type BidOrderResponse struct {
	Status       string
	Order_status int
}

type OrderParam struct {
	Token string
}

type OrderList struct {
	Order_id    int
	Cust_name   string
	Start_time  string
	Vertical_id int
}

type OrderListResponse struct {
	Status string
	Orders []OrderList
}

type OrderHistoryParam struct {
	Token  string
	Limit  int
	Offset int
}
type OrderHistory struct {
	Order_id       int
	Cust_name      string
	Service_name   string
	Service_detail string
	Gojek_credit   int
	Total_payment  float32
	Start_time     string
	Vertical_id    int
}

type OrderHistoryResponse struct {
	Status string
	Orders []OrderHistory
}

func BidOrderGlam(db *sql.DB, bo *BidOrderParam, bor *BidOrderResponse) {

	maxSp := 1
	serviceComodity := getComodity(db, bo.Service_id)
	if serviceComodity == 0 {
		maxSp = 5
	}

	if countBidOrderGlam(db, bo.Order_id) < maxSp {
		tx, err := db.Begin()
		if err != nil {
			//log.Fatal(err)
			s := err.Error()
			bor.Status = s

		} else {

			defer tx.Rollback()

			stmt, err := tx.Prepare(`INSERT INTO glam_order_asignee ( glam_order_id, sp_id) 
			VALUES (?, ?)`)

			if err != nil {
				log.Fatal(err)
				s := err.Error()
				bor.Status = s
			} else {
				_, err := stmt.Exec(bo.Order_id, bo.Sp_id)
				if err != nil {
					log.Fatal(err)
					s := err.Error()
					bor.Status = s

				} else {

					err = tx.Commit()

					if err != nil {
						//log.Fatal(err)
						s := err.Error()
						bor.Status = s

					}
				}
			}
			stmt.Close()
			if countBidOrderGlam(db, bo.Order_id) == maxSp {
				if serviceComodity == 0 {
					//TODO cust id blast
					fmt.Println("TODO cust id blast")
				} else {
					fmt.Println("UPDATE SP")
					UpdateOrderSP(db, bo)
				}

			}

		}
	} else {
		bor.Status = "Not available anymore"
	}
	bor.Order_status = getOrderStatus(db, bo.Order_id)

}

func countBidOrderGlam(db *sql.DB, orderId int) int {
	stmt, _ := db.Prepare("SELECT COUNT(*) FROM glam_order_asignee WHERE glam_order_id = ?")
	rows, err := stmt.Query(orderId)
	count := 0

	if err != nil {
		log.Fatal(err)
	} else {
		defer rows.Close()

		for rows.Next() {

			err := rows.Scan(&count)
			if err != nil {
				log.Fatal(err)
				break
			}
		}

		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
	}
	return count

}

func getComodity(db *sql.DB, serviceId int) int {
	stmt, _ := db.Prepare("SELECT service_comodity FROM service WHERE service_id = ?")
	rows, err := stmt.Query(serviceId)
	comodity := 0

	if err != nil {
		log.Fatal(err)
	} else {
		defer rows.Close()

		for rows.Next() {

			err := rows.Scan(&comodity)
			if err != nil {
				log.Fatal(err)
				break
			}
		}

		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
	}
	return comodity

}

func getOrderStatus(db *sql.DB, orderId int) int {
	stmt, _ := db.Prepare("SELECT order_status FROM glam_order WHERE glam_order_id = ?")
	rows, err := stmt.Query(orderId)
	var glamOrderStatus int

	if err != nil {
		log.Fatal(err)
	} else {
		defer rows.Close()

		for rows.Next() {

			err := rows.Scan(&glamOrderStatus)
			if err != nil {
				log.Fatal(err)
				break
			}
		}

		err = rows.Err()
		if err != nil {
			log.Fatal(err)
		}
	}
	return glamOrderStatus

}

func UpdateOrderSP(db *sql.DB, bo *BidOrderParam) {

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)

	} else {
		defer tx.Rollback()
		stmt, err := tx.Prepare(`UPDATE glam_order SET sp_id =?
				WHERE glam_order_id=?;`)
		if err != nil {
			log.Fatal(err)
		} else {
			_, err := stmt.Exec(bo.Sp_id, bo.Order_id)
			if err != nil {
				log.Fatal(err)
			} else {
				err = tx.Commit()
				if err != nil {
					log.Fatal(err)
				} else {
					//Asign sp in order success

				}
			}
		}
		stmt.Close()
	}
}

func GetOrderList(db *sql.DB, o *OrderParam, olr *OrderListResponse) {

	query := `SELECT glo.glam_order_id, glo.cust_name, glo.order_start_time, 1 as vertical_id
						FROM glam_order glo, service_provider sp
							WHERE glo.sp_id = sp.sp_id AND sp.sp_token =? AND glo.order_status = 3
						UNION ALL
			SELECT mo.massage_order_id, mo.cust_name, mo.order_start_time, 2 as vertical_id
						FROM massage_order mo, service_provider sp
						WHERE mo.sp_id = sp.sp_id AND sp.sp_token =? AND mo.order_status = 3
						UNION ALL
			SELECT co.clean_order_id, co.cust_name, co.order_start_time, 3 as vertical_id
						FROM clean_order co, service_provider sp
						WHERE co.sp_id = sp.sp_id AND sp.sp_token =? AND co.order_status = 3
						UNION ALL
			SELECT ao.auto_order_id, ao.cust_name, ao.order_start_time, 4 as vertical_id
						FROM auto_order ao, service_provider sp
						WHERE ao.sp_id = sp.sp_id AND sp.sp_token =? AND ao.order_status = 3
						
			`

	stmt, _ := db.Prepare(query)
	rows, err := stmt.Query(o.Token, o.Token, o.Token, o.Token)

	//var avai_sp []ServiceProvider
	if err != nil {
		log.Fatal(err)
	} else {
		defer rows.Close()

		var ols []OrderList
		ln := 1

		for rows.Next() {

			ol := OrderList{}
			err := rows.Scan(&ol.Order_id, &ol.Cust_name, &ol.Start_time, &ol.Vertical_id)

			if err != nil {

				log.Fatal(err)
				break

			} else {
				ols = append(ols, ol)
				ln++
			}
		}

		err = rows.Err()
		if err != nil {
			//			s := err.Error()
			log.Fatal(err)
		} else {
			olr.Orders = ols
		}
	}
}

func GetOrderHistory(db *sql.DB, o *OrderHistoryParam, ohr *OrderHistoryResponse) {

	query := `SELECT glo.glam_order_id, glo.cust_name, s.service_name, glo.order_detail, sp.sp_pulsa, glo.order_total_price,glo.order_start_time as start_time, 1 as vertical_id
				FROM glam_order glo, service s, service_provider sp
				WHERE glo.sp_id =sp.sp_id AND s.service_id=glo.service_id AND sp.sp_token =? AND glo.order_status = 7
			UNION ALL
			SELECT mo.massage_order_id, mo.cust_name, s.service_name, mo.order_detail, sp.sp_pulsa, mo.order_total_price,mo.order_start_time as start_time, 2 as vertical_id
				FROM massage_order mo, service s, service_provider sp
				WHERE mo.sp_id =sp.sp_id AND s.service_id=mo.service_id AND sp.sp_token =? AND mo.order_status = 7
			UNION ALL
			SELECT co.clean_order_id, co.cust_name, s.service_name, co.order_detail, sp.sp_pulsa, co.order_total_price,co.order_start_time as start_time, 3 as vertical_id
				FROM clean_order co, service s, service_provider sp
				WHERE co.sp_id =sp.sp_id AND s.service_id=co.service_id AND sp.sp_token =? AND co.order_status = 7
			UNION ALL
			SELECT ao.auto_order_id, ao.cust_name, s.service_name, ao.order_detail, sp.sp_pulsa, ao.order_total_price,ao.order_start_time as start_time, 4 as vertical_id
				FROM auto_order ao, service s, service_provider sp
				WHERE ao.sp_id =sp.sp_id AND s.service_id=ao.service_id AND sp.sp_token =? AND ao.order_status = 7
			ORDER BY start_time ASC LIMIT ?, ?`

	stmt, _ := db.Prepare(query)
	rows, err := stmt.Query(o.Token, o.Token, o.Token, o.Token, o.Offset, o.Limit)

	//var avai_sp []ServiceProvider
	if err != nil {
		log.Fatal(err)
	} else {
		defer rows.Close()

		var ohs []OrderHistory
		ln := 1

		for rows.Next() {

			oh := OrderHistory{}
			err := rows.Scan(&oh.Order_id, &oh.Cust_name, &oh.Service_name, &oh.Service_detail, &oh.Gojek_credit, &oh.Total_payment, &oh.Start_time, &oh.Vertical_id)

			if err != nil {

				log.Fatal(err)
				break

			} else {
				ohs = append(ohs, oh)
				ln++
			}
		}

		err = rows.Err()
		if err != nil {
			//			s := err.Error()
			log.Fatal(err)
		} else {
			ohr.Orders = ohs
		}
	}
}
