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

func OrderGlamHandler(w http.ResponseWriter, r *http.Request, ctx *co.AppContext) {

	var (
		lr ds.GlamOrderRespon
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

			res := ds.GlamOrder{}

			if err := json.Unmarshal(buf.Bytes(), &res); err != nil {
				//			panic(err)
				s := err.Error()
				lr.Status = s
				log.Fatal(err)
			} else {

				ds.AddGlamOrder(ctx.Db, &res, &lr)
				//fmt.Println(res)

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
		w.WriteHeader(201)
		fmt.Fprintf(w, bytes.NewBuffer(resp).String())

	} else {

		fmt.Fprintf(w, "{ \"Status\" : \""+err.Error()+"\"}")
	}
}

/*func GetOrderStatusById(w http.ResponseWriter, r *http.Request, ctx *co.AppContext) {

	var (
		glor ds.GlamOrderRespon
	)
	glor.Status = ctx.Status
	glor.App_version = "1"

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

			res := ds.GlamOrder{}

			if err := json.Unmarshal(buf.Bytes(), &res); err != nil {
				//			panic(err)
				s := err.Error()
				glor.Status = s

			} else {

				ds.GetStatusOrderById(ctx.Db, &res, &glor)
				//fmt.Println(res)

			}
			//		fmt.Println(res)
			//		fmt.Println(res.Page)
		} else {

			glor.Status = "Error parameter in body required"
		}
	}

	//lr.Service_provider = sp

	//fmt.Printf("s= %v\n", lr)
	resp, err := json.Marshal(glor)

	if err == nil && glor.Status == "success" {
		w.WriteHeader(201)
		fmt.Fprintf(w, bytes.NewBuffer(resp).String())

	} else {
		w.WriteHeader(400)
		if err != nil {
			fmt.Fprintf(w, "{ \"Status\" : \""+err.Error()+"\"}")
		} else {
			fmt.Fprintf(w, "{ \"Status\" : \""+glor.Status+"\"}")
		}
	}
}*/

/*
func GetOrderById(w http.ResponseWriter, r *http.Request, ctx *co.AppContext) {

	res := ds.GlamOrder{}
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
}*/
