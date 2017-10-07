// updateorderstatus
package serviprov

import (
	"database/sql"
	//"fmt"
	//"log"
	_ "mysql-master"
	s "strings"
	"strconv"
	"time"
)

type UpdateOrderStatusRequest struct {
	OrderID     int    `json:"Order_id"`
	OrderStatus int    `json:"Order_status"`
	Token       string `json:"Token"`
	VerticalID  int    `json:"Vertical_id"`
}

type UpdateOrderStatusRespon struct {
	Error   int    `json:"Error"`
	OrderID int    `json:"Order_id"`
	Status  string `json:"Status"`
}

func UpdateOrderStatus(db *sql.DB, sp *UpdateOrderStatusRequest, mor *UpdateOrderStatusRespon) {
	var orderTable [5]string
	orderTable[1] = "glam_order"
	orderTable[2] = "massage_order"
	orderTable[3] = "clean_order"
	orderTable[4] = "auto_order"
	
	
	// check if order exist/non-exist
	row, err := db.Query(s.Join([]string{
		"SELECT * FROM ",
		orderTable[sp.VerticalID],
		"WHERE ",
		s.Join([]string{orderTable[sp.VerticalID], "id"}, "_"),	" = ? ",
	}, " "), sp.OrderID)
	
	/*fmt.Println(s.Join([]string{
		"SELECT * FROM ",
		orderTable[sp.VerticalID],
		"WHERE ",
		s.Join([]string{orderTable[sp.VerticalID], "id"}, "_"),	" = ? ",
	}, " "))*/
	
	if err != nil {
		//log.Fatal(err)
		s := err.Error()
		mor.Status = s
		mor.Error = 1
		return
	}
	
	if (!row.Next()) {
		mor.Status = "Order not found"
		mor.Error = 1
		return	
	}
	
	tx, err := db.Begin()

	if err != nil {
		//log.Fatal(err)
		s := err.Error()
		mor.Status = s
		mor.Error = 1
		return
	}
	
	defer tx.Rollback()
	
	current_time := time.Now()
	//current_time.String()
	current_time.Format("2006-01-02 15:04:05")
	//fmt.Println(t.Format("2006-01-02 15:04:05"))
	
	if sp.OrderStatus == 5{
		stmt, err := db.Prepare(s.Join([]string{
			"UPDATE ",
			orderTable[sp.VerticalID],
			"SET order_status = ", strconv.Itoa(5),
			",work_otw_time = ", "'"+current_time.Format("2006-01-02 15:04:05")+"'", 
			"WHERE ",
			s.Join([]string{orderTable[sp.VerticalID], "id"}, "_"),	" = ? ",
		}, " "))
	
	/*fmt.Println(s.Join([]string{
			"UPDATE ",
			orderTable[sp.VerticalID],
			"SET order_status = ", strconv.Itoa(5),
			",work_otw_time = ", current_time.String(), 
			"WHERE ",
			s.Join([]string{orderTable[sp.VerticalID], "id"}, "_"),	" = ? ",
		}, " "))*/
	
		if err != nil {
			//log.Fatal(err)
			s := err.Error()
			mor.Status = s
			mor.Error = 1
			return
		}
		
		_, err = stmt.Exec(
			sp.OrderID,
		)
		
		if err != nil {
			s := err.Error()
			mor.Status = s
			mor.Error = 1
			return
		}
		
		err = tx.Commit()
		
		if err != nil {
			s := err.Error()
			mor.Status = s
			mor.Error = 1
			return
		}
		
		mor.Error = 0
		mor.Status = "success"
		
		stmt.Close()
	}
	
	if sp.OrderStatus == 6{
		stmt, err := db.Prepare(s.Join([]string{
			"UPDATE ",
			orderTable[sp.VerticalID],
			"SET order_status = ", strconv.Itoa(6),
			",order_actual_start_time = ", "'"+current_time.Format("2006-01-02 15:04:05")+"'", 
			"WHERE ",
			s.Join([]string{orderTable[sp.VerticalID], "id"}, "_"),	" = ? ",
		}, " "))
	
		if err != nil {
			//log.Fatal(err)
			s := err.Error()
			mor.Status = s
			mor.Error = 1
			return
		}
		
		_, err = stmt.Exec(
			sp.OrderID,
		)
		
		if err != nil {
			s := err.Error()
			mor.Status = s
			mor.Error = 1
			return
		}
		
		err = tx.Commit()
		
		if err != nil {
			s := err.Error()
			mor.Status = s
			mor.Error = 1
			return
		}
		
		mor.Error = 0
		mor.Status = "success"
		
		stmt.Close()
	}
	
	if sp.OrderStatus == 7{
		stmt, err := db.Prepare(s.Join([]string{
			"UPDATE ",
			orderTable[sp.VerticalID],
			"SET order_status = ", strconv.Itoa(7),
			",order_actual_finish_time = ", "'"+current_time.Format("2006-01-02 15:04:05")+"'", 
			"WHERE ",
			s.Join([]string{orderTable[sp.VerticalID], "id"}, "_"),	" = ? ",
		}, " "))
	
		if err != nil {
			//log.Fatal(err)
			s := err.Error()
			mor.Status = s
			mor.Error = 1
			return
		}
		
		_, err = stmt.Exec(
			sp.OrderID,
		)
		
		if err != nil {
			s := err.Error()
			mor.Status = s
			mor.Error = 1
			return
		}
		
		err = tx.Commit()
		
		if err != nil {
			s := err.Error()
			mor.Status = s
			mor.Error = 1
			return
		}
		
		mor.Error = 0
		mor.Status = "success"
		
		stmt.Close()
	}
	
	
}
