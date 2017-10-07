// gservice
package common

import (
	"database/sql"
	//"fmt"
	_ "mysql-master"
	s "strings"
)

type Duration struct {
	Id       int
	Duration int
}

type Service struct {
	Id        int
	Name      string
	Price     float32
	Durations []Duration
}

type ServiceRespon struct {
	Status      string
	App_version string
	Services    []Service
}

func LoadService(db *sql.DB, idRole int, sp *ServiceRespon) {

	rows, err := db.Query(s.Join([]string{"SELECT s.id, s.name, s.price ",
		"FROM service AS s WHERE s.role = ? ORDER BY id"}, " "), idRole)
	//fmt.Printf("phone = %s", sp.Service_provider.Phone)

	if err != nil {
		s := err.Error()
		sp.Status = s

	} else {
		defer rows.Close()

		srv := make([]Service, 0, 20)
		ln := 1

		for rows.Next() {
			srv = srv[0:ln] // extend length by 1

			err := rows.Scan(&srv[ln-1].Id, &srv[ln-1].Name, &srv[ln-1].Price)

			if err != nil {
				s := err.Error()
				sp.Status = s
				break

			} else {

				LoadDuration(db, sp, &srv[ln-1])
				ln++
			}
		}

		err = rows.Err()
		if err != nil {
			s := err.Error()
			sp.Status = s
		} else {
			sp.Services = srv
			//for j := 0; j < len(sp.Service); j++ {
			//	fmt.Printf("v=%+v\n", sp.Service[j])
			//}
		}
	}
}

func LoadDuration(db *sql.DB, sr *ServiceRespon, sp *Service) {

	rows, err := db.Query(s.Join([]string{"SELECT s.id, s.duration ",
		"FROM service_duration AS s WHERE s.service_id = ? ORDER BY duration"}, " "), sp.Id)
	//fmt.Printf("phone = %s", sp.Service_provider.Phone)

	if err != nil {
		s := err.Error()
		sr.Status = s

	} else {
		defer rows.Close()

		srv := make([]Duration, 0, 20)
		ln := 1

		for rows.Next() {
			srv = srv[0:ln] // extend length by 1

			err := rows.Scan(&srv[ln-1].Id, &srv[ln-1].Duration)
			ln++

			if err != nil {
				s := err.Error()
				sr.Status = s
				break
			}
		}

		err = rows.Err()
		if err != nil {
			s := err.Error()
			sr.Status = s
		} else {
			sp.Durations = srv
			//for j := 0; j < len(sp.Service); j++ {
			//	fmt.Printf("v=%+v\n", sp.Service[j])
			//}
		}
	}
}
