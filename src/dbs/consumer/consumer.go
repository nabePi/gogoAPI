package consumer

import (
	"database/sql"
	//"fmt"
	//"log"
	_ "mysql-master"
	//"encoding/json"
	//"bytes"
	//"html"
	s "strings"
	"strconv"
	am "contler"
)

type CancelOrderRequest struct {
	Sp_Id		int `json:"sp_id"`	
	OrderID		int	`json:"order_id"`
	VerticalID	int	`json:"vertical_id"`
}

type CancelOrderRespon struct {
	Error		int 	`json:"error"`
	Status      string	`json:"status"`
}

type CancelOrderBlastMessage struct {
	Sp_id			[]int `json:"sp_id"'`
	Blast_type		int	`json:"blast_type"`
	Message_type 	string `json:"message_type"`
	Message blastMessage
}

type blastMessage struct {
	Order_id	int `json:"order_id"`
	Vertical_id	int	`json:"vertical_id"`	
}

func CancelOrder(db *sql.DB, sp *CancelOrderRequest, mor *CancelOrderRespon) {
	var orderTable [5]string
	orderTable[GO_GLAM_VERTICAL_ID] = "glam_order"
	orderTable[GO_CLEAN_VERTICAL_ID] = "clean_order"
	orderTable[GO_MASSAGE_VERTICAL_ID] = "massage_order"
	orderTable[GO_AUTO_VERTICAL_ID] = "auto_order"

	// check if order exist/non-exist
	row, err := db.Query(s.Join([]string{
		"SELECT sp_id FROM ",
		orderTable[sp.VerticalID],
		"WHERE ",
		s.Join([]string{orderTable[sp.VerticalID], "id"}, "_"),	" = ? ",
	}, " "), sp.OrderID)
	

	
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

	row.Scan(&sp.Sp_Id)

	tx, err := db.Begin()

	if err != nil {
		//log.Fatal(err)
		s := err.Error()
		mor.Status = s
		mor.Error = 1
		return
	}
	
	defer tx.Rollback()
	
	stmt, err := db.Prepare(s.Join([]string{
		"UPDATE ",
		orderTable[sp.VerticalID],
		"SET order_status = ", strconv.Itoa(ORDER_STATUS_CANCELED),
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
	
	cancelOrderBlastMessage := CancelOrderBlastMessage {
		[]int{sp.Sp_Id},
		BLAST_TYPE_CANCEL_ORDER,
		"cancelOrder",
		blastMessage {
			sp.OrderID,
			sp.VerticalID,
		},
	}	
	//jsonMsg, err := json.Marshal(cancelOrderBlastMessage)
	//am.PublishRoutine(am.PubnubChannel, html.UnescapeString(bytes.NewBuffer(jsonMsg).String()))

	am.PublishRoutine(am.PubnubChannel, cancelOrderBlastMessage)
	
	
	mor.Error = 0
	mor.Status = "success"
	
	stmt.Close()
}
