// cmassage
package consumer

import (
	"bytes"
	co "contler"
	ds "dbs/consumer"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func GetOrderStatusById(w http.ResponseWriter, r *http.Request, ctx *co.AppContext) {

	var (
		osr ds.OrderStatusResponse
	)
	osr.Status = ctx.Status
	osr.App_version = "1"

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

			res := ds.Order{}

			if err := json.Unmarshal(buf.Bytes(), &res); err != nil {
				//			panic(err)
				s := err.Error()
				osr.Status = s
				log.Fatal(err)

			} else {

				ds.GetStatusOrderById(ctx.Db, &res, &osr)
				//fmt.Println(res)

			}
			//		fmt.Println(res)
			//		fmt.Println(res.Page)
		} else {

			osr.Status = "Error parameter in body required"
		}
	}

	//lr.Service_provider = sp

	//fmt.Printf("s= %v\n", lr)
	resp, err := json.Marshal(osr)

	if err == nil && osr.Status == "success" {
		w.WriteHeader(201)
		fmt.Fprintf(w, bytes.NewBuffer(resp).String())

	} else {
		w.WriteHeader(400)
		if err != nil {
			fmt.Fprintf(w, "{ \"Status\" : \""+err.Error()+"\"}")
		} else {
			fmt.Fprintf(w, "{ \"Status\" : \""+osr.Status+"\"}")
		}
	}
}

func GetOrderById(w http.ResponseWriter, r *http.Request, ctx *co.AppContext) {

	res := ds.Order{}
	if ctx.Status == "success" {

		if r.Body != nil {

			buf := new(bytes.Buffer)
			buf.ReadFrom(r.Body)

			if err := json.Unmarshal(buf.Bytes(), &res); err != nil {
				//			panic(err)
				//				s := err.Error()
				log.Fatal(err)
			} else {

				ds.GetOrderById(ctx.Db, &res)

			}

		} else {

			//glor.Status = "Error parameter in body required"
		}
	}

	resp, err := json.Marshal(&res)

	if err == nil && res.Status == "success" {
		w.WriteHeader(201)
		fmt.Fprintf(w, bytes.NewBuffer(resp).String())

	} else {
		if err != nil {
			fmt.Fprintf(w, "{ \"Status\" : \""+err.Error()+"\"}")
		} else {
			fmt.Fprintf(w, "{ \"Status\" : \""+res.Status+"\"}")
		}
	}
}
