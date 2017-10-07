package serviprov

import (
	"bytes"
	co "contler"
	ds "dbs/serviprov"
	"encoding/json"
	"fmt"
	//"log"
	"net/http"
	//s "strings"
)

type LoginParam struct {
	Sp_phone     string
	Sp_imei      string
	Sp_pubnub_id string
	Sp_pin       string
}

func LoginHandler(w http.ResponseWriter, r *http.Request, ctx *co.AppContext) {

	var (
		//sp common.SrvProv
		lr ds.LoginResponse
	)
	lr.Status = ctx.Status
	lr.App_version = "1"

	if ctx.Status == "success" {

		if r.Body != nil {
			//var p []byte

			buf := new(bytes.Buffer)
			buf.ReadFrom(r.Body)
			//s := buf.String() // Does a complete copy of the bytes in the buffer.
			//fmt.Fprintf(w, "n = %s", s)

			//var byt []byte
			//r.Body.Read(byt)
			//		res := Response1{}
			//var dat map[string]interface{}

			res := LoginParam{}

			//if err := json.Unmarshal(buf.Bytes(), &dat); err != nil {
			if err := json.Unmarshal(buf.Bytes(), &res); err != nil {
				//			panic(err)
				s := err.Error()
				lr.Status = s

			} else {

				//fmt.Printf("phonex = %v", res)

				lr.Service_provider.Phone = res.Sp_phone
				//fmt.Printf("phonep = %s", res.Sp_phone)

				ds.LoadSP(ctx.Db, &lr)
				if lr.Status == "success" {
					ds.LoadSPServices(ctx.Db, &lr)
					if lr.Status == "success" {
						ds.LoadSPEquip(ctx.Db, &lr)
					}
				}
			}
			//		fmt.Println(res)
			//		fmt.Println(res.Page)
		} else {

			lr.Status = "Error parameter in body required"
		}
	}

	//lr.Service_provider = sp

	//fmt.Printf("s= %v\n", lr)
	resp, err := json.Marshal(lr)

	if err == nil {
		fmt.Fprintf(w, bytes.NewBuffer(resp).String())

	} else {

		fmt.Fprintf(w, "{ \"Status\" : \""+err.Error()+"\"}")
	}
}
