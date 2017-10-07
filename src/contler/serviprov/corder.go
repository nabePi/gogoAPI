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

func BidOrderHandler(w http.ResponseWriter, r *http.Request, ctx *co.AppContext) {

	var (
		//sp common.SrvProv
		bor ds.BidOrderResponse
	)
	bor.Status = ctx.Status

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

			res := ds.BidOrderParam{}

			//if err := json.Unmarshal(buf.Bytes(), &dat); err != nil {
			if err := json.Unmarshal(buf.Bytes(), &res); err != nil {
				//			panic(err)
				s := err.Error()
				bor.Status = s

			} else {
				switch res.Vertical_id {
				case 1:
					ds.BidOrderGlam(ctx.Db, &res, &bor)
					break
					
				case 4:
					ds.BidOrderAuto(ctx.Db, &res, &bor)
				}
				

				//fmt.Printf("phonex = %v", res)

				/*lr.Service_provider.Phone = res.Sp_phone
				//fmt.Printf("phonep = %s", res.Sp_phone)

				ds.LoadSP(ctx.Db, &lr)
				if lr.Status == "success" {
					ds.LoadSPServices(ctx.Db, &lr)
					if lr.Status == "success" {
						ds.LoadSPEquip(ctx.Db, &lr)
					}
				}*/
			}
			//		fmt.Println(res)
			//		fmt.Println(res.Page)
		} else {

			bor.Status = "Error parameter in body required"
		}
	}

	//lr.Service_provider = sp

	//fmt.Printf("s= %v\n", lr)
	resp, err := json.Marshal(bor)

	if err == nil {
		fmt.Fprintf(w, bytes.NewBuffer(resp).String())

	} else {

		fmt.Fprintf(w, "{ \"Status\" : \""+err.Error()+"\"}")
	}
}

func GetOrderList(w http.ResponseWriter, r *http.Request, ctx *co.AppContext) {

	var olr ds.OrderListResponse

	olr.Status = ctx.Status

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

			res := ds.OrderParam{}

			//if err := json.Unmarshal(buf.Bytes(), &dat); err != nil {
			if err := json.Unmarshal(buf.Bytes(), &res); err != nil {
				//			panic(err)
				s := err.Error()
				olr.Status = s

			} else {
				ds.GetOrderList(ctx.Db, &res, &olr)

				//fmt.Printf("phonex = %v", res)

				/*lr.Service_provider.Phone = res.Sp_phone
				//fmt.Printf("phonep = %s", res.Sp_phone)

				ds.LoadSP(ctx.Db, &lr)
				if lr.Status == "success" {
					ds.LoadSPServices(ctx.Db, &lr)
					if lr.Status == "success" {
						ds.LoadSPEquip(ctx.Db, &lr)
					}
				}*/
			}
			//		fmt.Println(res)
			//		fmt.Println(res.Page)
		} else {

			olr.Status = "Error parameter in body required"
		}
	}

	//lr.Service_provider = sp

	//fmt.Printf("s= %v\n", lr)
	resp, err := json.Marshal(olr)

	if err == nil {
		fmt.Fprintf(w, bytes.NewBuffer(resp).String())

	} else {

		fmt.Fprintf(w, "{ \"Status\" : \""+err.Error()+"\"}")
	}
}

func GetOrderHistory(w http.ResponseWriter, r *http.Request, ctx *co.AppContext) {

	var ohr ds.OrderHistoryResponse

	ohr.Status = ctx.Status

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

			res := ds.OrderHistoryParam{}

			//if err := json.Unmarshal(buf.Bytes(), &dat); err != nil {
			if err := json.Unmarshal(buf.Bytes(), &res); err != nil {
				//			panic(err)
				s := err.Error()
				ohr.Status = s

			} else {
				ds.GetOrderHistory(ctx.Db, &res, &ohr)

				//fmt.Printf("phonex = %v", res)

				/*lr.Service_provider.Phone = res.Sp_phone
				//fmt.Printf("phonep = %s", res.Sp_phone)

				ds.LoadSP(ctx.Db, &lr)
				if lr.Status == "success" {
					ds.LoadSPServices(ctx.Db, &lr)
					if lr.Status == "success" {
						ds.LoadSPEquip(ctx.Db, &lr)
					}
				}*/
			}
			//		fmt.Println(res)
			//		fmt.Println(res.Page)
		} else {

			ohr.Status = "Error parameter in body required"
		}
	}

	//lr.Service_provider = sp

	//fmt.Printf("s= %v\n", lr)
	resp, err := json.Marshal(ohr)

	if err == nil {
		fmt.Fprintf(w, bytes.NewBuffer(resp).String())

	} else {

		fmt.Fprintf(w, "{ \"Status\" : \""+err.Error()+"\"}")
	}
}