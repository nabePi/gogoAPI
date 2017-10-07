package serviprov

import (
	"database/sql"
	//"fmt"
	_ "mysql-master"
	s "strings"
)

type LoginResponse struct {
	Status           string
	App_version      string
	Service_provider SrvProv
}

type Service struct {
	Tipe      string
	Price     int
	Price_max int
}

type Equipment struct {
	Name string
}

type SrvProv struct {
	Id         int
	Name       string
	Gender     int
	Role       int
	Phone      string
	Profile    string
	Rate       float32
	Foto_url   string
	Services   []Service
	Equipments []Equipment
}

func LoadSP(db *sql.DB, sp *LoginResponse) {

	rows, err := db.Query(s.Join([]string{"SELECT sp.id, sp.name, sp.gender, sp.role,",
		"sp.phone, sp.profile, sp.rate, sp.foto_url FROM service_provider AS sp",
		"WHERE sp.phone = ?"}, " "), sp.Service_provider.Phone)
	//fmt.Printf("phone = %s", sp.Service_provider.Phone)

	if err != nil {
		s := err.Error()
		sp.Status = s

	} else {
		defer rows.Close()

		if rows.Next() {
			var spr SrvProv

			err := rows.Scan(&spr.Id, &spr.Name, &spr.Gender, &spr.Role, &spr.Phone,
				&spr.Profile, &spr.Rate, &spr.Foto_url)

			if err != nil {
				s := err.Error()
				sp.Status = s

			} else {
				sp.Service_provider = spr
			}

		} else {

			sp.Status = "Not found"

		}

		err = rows.Err()
		if err != nil {
			s := err.Error()
			sp.Status = s
		}
	}
}

func LoadSPServices(db *sql.DB, sp *LoginResponse) {

	rows, err := db.Query(s.Join([]string{"SELECT s.name, ss.price, ss.price_max",
		" FROM service AS s INNER JOIN sp_service AS ss ON ss.service_id = s.id",
		"WHERE ss.sp_id = ?"}, " "), sp.Service_provider.Id)

	if err != nil {
		s := err.Error()
		sp.Status = s

	} else {

		defer rows.Close()

		srv := make([]Service, 0, 20)
		ln := 1

		for rows.Next() {
			srv = srv[0:ln] // extend length by 1

			err := rows.Scan(&srv[ln-1].Tipe, &srv[ln-1].Price, &srv[ln-1].Price_max)
			ln++

			if err != nil {
				s := err.Error()
				sp.Status = s
				break
			}
		}

		err = rows.Err()
		if err != nil {
			s := err.Error()
			sp.Status = s
		} else {
			sp.Service_provider.Services = srv
			//for j := 0; j < len(sp.Service); j++ {
			//	fmt.Printf("v=%+v\n", sp.Service[j])
			//}
		}
	}
}

func LoadSPEquip(db *sql.DB, sp *LoginResponse) {

	rows, err := db.Query(s.Join([]string{"SELECT eqp.name FROM equipment AS eqp ",
		"INNER JOIN sp_equipment AS spe ON spe.equip_id = eqp.id ",
		"INNER JOIN service_provider AS sp ON spe.sp_id = sp.id WHERE sp.id = ?"}, " "),
		sp.Service_provider.Id)

	if err != nil {
		s := err.Error()
		sp.Status = s

	} else {

		defer rows.Close()

		srv := make([]Equipment, 0, 20)
		ln := 1

		for rows.Next() {
			srv = srv[0:ln] // extend length by 1

			err := rows.Scan(&srv[ln-1].Name)
			ln++

			if err != nil {
				s := err.Error()
				sp.Status = s
				break
			}
		}

		err = rows.Err()
		if err != nil {
			s := err.Error()
			sp.Status = s

		} else {
			sp.Service_provider.Equipments = srv
			//for j := 0; j < len(sp.Service); j++ {
			//	fmt.Printf("v=%+v\n", sp.Service[j])
			//}
		}
	}
}
