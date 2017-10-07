package consumer

import (
	"bytes"
	co "contler"
	ds "dbs/common"
	"encoding/json"
	"fmt"
	//"log"
	"net/http"
	//s "strings"
)

type ServiceParam struct {
	Type_id int
}

func ServiceOfferHandler(w http.ResponseWriter, r *http.Request, ctx *co.AppContext) {

	var (
		lr ds.ServiceRespon
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

			res := ServiceParam{}

			if err := json.Unmarshal(buf.Bytes(), &res); err != nil {
				//			panic(err)
				s := err.Error()
				lr.Status = s

			} else {

				ds.LoadService(ctx.Db, res.Type_id, &lr)

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
