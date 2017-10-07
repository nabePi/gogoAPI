// orderdetail
package serviprov

import (
	"database/sql"
	//"fmt"
	//"log"
	_ "mysql-master"
	s "strings"
	//"strconv"
)

type OrderDetailRequest struct{
	OrderID int    `json:"Order_id"`
	Token   string `json:"Token"`
}

type DetailResponse struct{
	CustomerName       string   `json:"Customer_name"`
	CustomerPhone      string   `json:"Customer_phone"`
	Ondemand           bool     `json:"Ondemand"`
	OrderDuration      int   `json:"Order_duration"`
	OrderFinishTime    string   `json:"Order_finish_time"`
	OrderID            int   `json:"Order_id"`
	OrderPaymentMethod int   `json:"Order_payment_method"`
	OrderPlaceInput    string   `json:"Order_place_input"`
	OrderPlaceLat      float64  `json:"Order_place_lat"`
	OrderPlaceLng      float64  `json:"Order_place_lng"`
	OrderPlaceMap      string   `json:"Order_place_map"`
	OrderPlaceNote     string   `json:"Order_place_note"`
	OrderPrice         float32      `json:"Order_price"`
	OrderSpecialNote   string   `json:"Order_special_note"`
	OrderStartTime     string   `json:"Order_start_time"`
	OrderStatus        int      `json:"Order_status"`
	ServiceDetail      string   `json:"Service_detail"`
	ServiceName        string   `json:"Service_name"`
	ServiceType        string   `json:"service_type"`
	Asignee            []string `json:"asignee"` //for go glam
	ItemDetail         int   `json:"item_detail"` //multi for auto
}

type OrderDetailResponse struct{
	Error  int    `json:"Error"`
	Order  DetailResponse `json:"Order"`
	Status string `json:"Status"`
}

func GetOrderDetail(db *sql.DB, sp *OrderDetailRequest, mor *OrderDetailResponse) {
	// check if token exist/non-exist
	sp_id := 0
	vertical_id := 0
	rowssp, err:= db.Query(s.Join([]string{"SELECT sp_id, sp_vertical_id ",
	"FROM golife_spapps.service_provider ",
	"WHERE sp_token = ? "}, " "), sp.Token)
	if err != nil{
		s := err.Error()
		mor.Status = s
		mor.Error = 1
	}else{
		defer rowssp.Close()
		
		if rowssp.Next(){
			err := rowssp.Scan(&sp_id, &vertical_id)
			
			if err != nil{
				s := err.Error()
				mor.Status = s
				mor.Error = 1
			}
			
		}else{
			mor.Status = "Service Provider Not Found"
			mor.Error = 1
			return
		}
	}
	
	var orderTable [5]string
	orderTable[1] = "glam_order"
	orderTable[2] = "massage_order"
	orderTable[3] = "clean_order"
	orderTable[4] = "auto_order"
	
	//detail GO GLAM
	if vertical_id == 1{
		
		rowdetail, err := db.Query(s.Join([]string{
			"SELECT A.cust_name, A.cust_phone, A.order_item_id, A.order_time, A.order_finish_time,",
			"A.glam_order_id, A.order_payment_method, A.order_address, A.order_lat, A.order_lng, ",
			"A.order_address, A.order_address, A.order_total_price, A.order_special_note, A.order_start_time,",
			"A.order_status, A.order_detail,B.service_name, D.vertical_name ",
			"FROM ",
			orderTable[vertical_id] + " A ",
			"left join service B on A.service_id = B.service_id",
			"left join service_provider C on A.sp_id = C.sp_id",
			"left join vertical D on C.sp_vertical_id = D.vertical_id",
			"WHERE ",
			s.Join([]string{orderTable[vertical_id], "id"}, "_"),	" = ? ",
		}, " "), sp.OrderID)
	
		/*fmt.Println(s.Join([]string{
			"SELECT A.cust_name, A.cust_phone, A.order_item_id, A.order_time, A.order_finish_time,",
			"A.massage_order_id, A.order_payment_method, A.order_address, A.order_lat, A.order_lng, ",
			"A.order_address, A.order_address, A.order_total_price, A.order_special_note, A.order_start_time,",
			"A.order_status, A.order_detail,B.service_name, D.vertical_name ",
			"FROM ",
			orderTable[vertical_id] + " A ",
			"left join service B on A.service_id = B.service_id",
			"left join service_provider C on A.sp_id = C.sp_id",
			"left join vertical D on C.sp_vertical_id = D.vertical_id",
			"WHERE ",
			s.Join([]string{orderTable[vertical_id], "id"}, "_"),	" = ? ",
		}, " "))*/
	
		if err != nil{
			s := err.Error()
			mor.Status = s
			mor.Error = 1
		}else{
			defer rowdetail.Close()
			if !rowdetail.Next(){
				mor.Status = "Order Not Found"
				mor.Error = 1
			}else{
				var svname sql.NullString
				var svtype sql.NullString
				
				var detorder DetailResponse
				
				err := rowdetail.Scan(&detorder.CustomerName,
				&detorder.CustomerPhone, 
				&detorder.ItemDetail,
				&detorder.OrderDuration,
				&detorder.OrderFinishTime,
				&detorder.OrderID,
				&detorder.OrderPaymentMethod,
				&detorder.OrderPlaceInput,
				&detorder.OrderPlaceLat,
				&detorder.OrderPlaceLng,
				&detorder.OrderPlaceMap,
				&detorder.OrderPlaceNote,
				&detorder.OrderPrice,
				&detorder.OrderSpecialNote, 
				&detorder.OrderStartTime,
				&detorder.OrderStatus, 
				&detorder.ServiceDetail,
				&svname,
				&svtype)
				
				detorder.ServiceName = svname.String
				detorder.ServiceType = svtype.String
			
				if err != nil{
					s := err.Error()
					mor.Status = s
					mor.Error = 1
				}else{
					mor.Error = 0
					mor.Status = "success"
					mor.Order = detorder
				}
			}
		}
		
	}
	
	//detail GO MASSAGE
	if vertical_id == 2{
		
		rowdetail, err := db.Query(s.Join([]string{
			"SELECT A.cust_name, A.cust_phone, A.order_item_id, A.order_time, A.order_finish_time,",
			"A.massage_order_id, A.order_payment_method, A.order_address, A.order_lat, A.order_lng, ",
			"A.order_address, A.order_address, A.order_total_price, A.order_special_note, A.order_start_time,",
			"A.order_status, A.order_detail,B.service_name, D.vertical_name ",
			"FROM ",
			orderTable[vertical_id] + " A ",
			"left join service B on A.service_id = B.service_id",
			"left join service_provider C on A.sp_id = C.sp_id",
			"left join vertical D on C.sp_vertical_id = D.vertical_id",
			"WHERE ",
			s.Join([]string{orderTable[vertical_id], "id"}, "_"),	" = ? ",
		}, " "), sp.OrderID)
	
		/*fmt.Println(s.Join([]string{
			"SELECT A.cust_name, A.cust_phone, A.order_item_id, A.order_time, A.order_finish_time,",
			"A.massage_order_id, A.order_payment_method, A.order_address, A.order_lat, A.order_lng, ",
			"A.order_address, A.order_address, A.order_total_price, A.order_special_note, A.order_start_time,",
			"A.order_status, A.order_detail,B.service_name, D.vertical_name ",
			"FROM ",
			orderTable[vertical_id] + " A ",
			"left join service B on A.service_id = B.service_id",
			"left join service_provider C on A.sp_id = C.sp_id",
			"left join vertical D on C.sp_vertical_id = D.vertical_id",
			"WHERE ",
			s.Join([]string{orderTable[vertical_id], "id"}, "_"),	" = ? ",
		}, " "))*/
	
		if err != nil{
			s := err.Error()
			mor.Status = s
			mor.Error = 1
		}else{
			defer rowdetail.Close()
			if !rowdetail.Next(){
				mor.Status = "Order Not Found"
				mor.Error = 1
			}else{
				var svname sql.NullString
				var svtype sql.NullString
				
				var detorder DetailResponse
				
				err := rowdetail.Scan(&detorder.CustomerName,
				&detorder.CustomerPhone, 
				&detorder.ItemDetail,
				&detorder.OrderDuration,
				&detorder.OrderFinishTime,
				&detorder.OrderID,
				&detorder.OrderPaymentMethod,
				&detorder.OrderPlaceInput,
				&detorder.OrderPlaceLat,
				&detorder.OrderPlaceLng,
				&detorder.OrderPlaceMap,
				&detorder.OrderPlaceNote,
				&detorder.OrderPrice,
				&detorder.OrderSpecialNote, 
				&detorder.OrderStartTime,
				&detorder.OrderStatus, 
				&detorder.ServiceDetail,
				&svname,
				&svtype)
				
				detorder.ServiceName = svname.String
				detorder.ServiceType = svtype.String
			
				if err != nil{
					s := err.Error()
					mor.Status = s
					mor.Error = 1
				}else{
					mor.Error = 0
					mor.Status = "success"
					mor.Order = detorder
				}
			}
		}
		
	}
	
	//detail GO CLEAN
	if vertical_id == 3{
		
	}
	
	//detail GO AUTO
	if vertical_id == 4{
		
		rowdetail, err := db.Query(s.Join([]string{
			"SELECT A.cust_name, A.cust_phone, A.order_item_id, A.order_time, A.order_finish_time,",
			"A.massage_order_id, A.order_payment_method, A.order_address, A.order_lat, A.order_lng, ",
			"A.order_address, A.order_address, A.order_total_price, A.order_special_note, A.order_start_time,",
			"A.order_status, A.order_detail,B.service_name, D.vertical_name ",
			"FROM ",
			orderTable[vertical_id] + " A ",
			"left join service B on A.service_id = B.service_id",
			"left join service_provider C on A.sp_id = C.sp_id",
			"left join vertical D on C.sp_vertical_id = D.vertical_id",
			"WHERE ",
			s.Join([]string{orderTable[vertical_id], "id"}, "_"),	" = ? ",
		}, " "), sp.OrderID)
	
		/*fmt.Println(s.Join([]string{
			"SELECT A.cust_name, A.cust_phone, A.order_item_id, A.order_time, A.order_finish_time,",
			"A.massage_order_id, A.order_payment_method, A.order_address, A.order_lat, A.order_lng, ",
			"A.order_address, A.order_address, A.order_total_price, A.order_special_note, A.order_start_time,",
			"A.order_status, A.order_detail,B.service_name, D.vertical_name ",
			"FROM ",
			orderTable[vertical_id] + " A ",
			"left join service B on A.service_id = B.service_id",
			"left join service_provider C on A.sp_id = C.sp_id",
			"left join vertical D on C.sp_vertical_id = D.vertical_id",
			"WHERE ",
			s.Join([]string{orderTable[vertical_id], "id"}, "_"),	" = ? ",
		}, " "))*/
	
		if err != nil{
			s := err.Error()
			mor.Status = s
			mor.Error = 1
		}else{
			defer rowdetail.Close()
			if !rowdetail.Next(){
				mor.Status = "Order Not Found"
				mor.Error = 1
			}else{
				var svname sql.NullString
				var svtype sql.NullString
				
				var detorder DetailResponse
				
				err := rowdetail.Scan(&detorder.CustomerName,
				&detorder.CustomerPhone, 
				&detorder.ItemDetail,
				&detorder.OrderDuration,
				&detorder.OrderFinishTime,
				&detorder.OrderID,
				&detorder.OrderPaymentMethod,
				&detorder.OrderPlaceInput,
				&detorder.OrderPlaceLat,
				&detorder.OrderPlaceLng,
				&detorder.OrderPlaceMap,
				&detorder.OrderPlaceNote,
				&detorder.OrderPrice,
				&detorder.OrderSpecialNote, 
				&detorder.OrderStartTime,
				&detorder.OrderStatus, 
				&detorder.ServiceDetail,
				&svname,
				&svtype)
				
				detorder.ServiceName = svname.String
				detorder.ServiceType = svtype.String
			
				if err != nil{
					s := err.Error()
					mor.Status = s
					mor.Error = 1
				}else{
					mor.Error = 0
					mor.Status = "success"
					mor.Order = detorder
				}
			}
		}
		
	}
			
}