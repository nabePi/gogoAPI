package serviprov

import (
	"database/sql"
	//"fmt"
	_ "mysql-master"
	s "strings"
)


type BidStatusRequest struct {
	Sp_Id			int `json:"sp_id"`
	Token			string `json:"token"`
	Status_bid		int		`json:"status_bid"`                         
}

type BidStatusRespon struct {
	Error		int	`json:"error"`
	Status		string `json:"status"`
	Status_bid	int `json:"status_bid"`                       
}

type UpdateLocationRequest struct {
	Sp_Id		int		`json:"sp_id"`
	Token		string  `json:"token"`
	Latitude	float32 `json:"latitude"`
	Longitude	float32 `json:"longitude"`
}

type UpdateLocationRespon struct {
	Error		int		`json:"error"`
	Status		string	`json:"status"`                     
}

type SPPulsaRequest struct {
	Token	string	`json:"token"`
}

type SPPulsaRespon struct {
	Error                int    `json:"error"`
	SpIncomeToday        string `json:"sp_income_today"`
	SpOutstandingBalance string `json:"sp_outstanding_balance"`
	SpPulsa              string `json:"sp_pulsa"`
	SpTotalOrderLeft     string `json:"sp_total_order_left"`
	Status               string `json:"status"`	
}

func SetBidStatus(db *sql.DB, sp *BidStatusRequest, mor *BidStatusRespon) {
	row, err := db.Query(s.Join([]string{
		"SELECT sp_id FROM service_provider WHERE sp_token = ?",
	}, " "), sp.Token)
	
	if err != nil {
		//log.Fatal(err)
		s := err.Error()
		mor.Status = s
		mor.Error = 1
		return
	}
	
	mor.Status_bid = sp.Status_bid
	
	if (!row.Next()) {
		mor.Status = "Invalid Token"
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
		"UPDATE service_provider SET sp_status = ? WHERE sp_id = ? AND sp_token = ?",
	}, " "))
	
	if err != nil {
		//log.Fatal(err)
		s := err.Error()
		mor.Status = s
		mor.Error = 1
		return
	}
	
	_, err = stmt.Exec(
		sp.Status_bid,
		sp.Sp_Id,
		sp.Token,
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

func UpdateLocation(db *sql.DB, sp *UpdateLocationRequest, mor *UpdateLocationRespon) {
	row, err := db.Query(s.Join([]string{
		"SELECT sp_id FROM service_provider WHERE  sp_token = ?",
	}, " "), sp.Token)
	
	if err != nil {
		//log.Fatal(err)
		s := err.Error()
		mor.Status = s
		mor.Error = 1
		return
	}
	
	if (!row.Next()) {
		mor.Status = "Invalid Token"
		mor.Error = 1
		return	
	}

	err = row.Scan(&sp.Sp_Id)
	if (err != nil) {
		mor.Status = err.Error()
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
	
	stmt, err := db.Prepare(s.Join([]string{
		"UPDATE sp_location SET latitude = ?, longitude = ? WHERE sp_id = ?",
	}, " "))
	
	if err != nil {
		//log.Fatal(err)
		s := err.Error()
		mor.Status = s
		mor.Error = 1
		return
	}
	
	_, err = stmt.Exec(                              
		sp.Latitude,    
		sp.Longitude,
		sp.Sp_Id,
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

func SPCheckPulsa(db *sql.DB, sp *SPPulsaRequest, mor *SPPulsaRespon) {
	row, err := db.Query(s.Join([]string{
		"SELECT sp_id, sp_vertical_id, sp_pulsa, sp_outstanding_balance FROM service_provider WHERE  sp_token = ?",
	}, " "), sp.Token)
	
	if err != nil {
		//log.Fatal(err)
		s := err.Error()
		mor.Status = s
		mor.Error = 1
		return
	}
	
	if (!row.Next()) {
		mor.Status = "Invalid Token"
		mor.Error = 1
		return	
	}
	
	var sp_id int
	var vertical_id int
	
	err = row.Scan(&sp_id, &vertical_id, &mor.SpPulsa, &mor.SpOutstandingBalance)

}