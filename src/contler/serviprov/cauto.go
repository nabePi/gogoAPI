
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

func CekItemAutoHandler(w http.ResponseWriter, r *http.Request, ctx *co.AppContext) {
	var (
		lr ds.CekItemAutoRespon
	)
	lr.Status = ctx.Status
	//lr.App_version = "1"

	//fmt.Fprintf(w, "%v", ctx.Status)
	w.Header().Set("Content-Type", "application/json")

	if ctx.Status == "success" {
		if r.Body != nil {
			//var p []byte

			buf := new(bytes.Buffer)
			buf.ReadFrom(r.Body)

			res := ds.CekItemAutoRequest{}

			if err := json.Unmarshal(buf.Bytes(), &res); err != nil {
				//			panic(err)
				s := err.Error()
				lr.Error = 1
				lr.Status = s

			} else {
				ds.CekItemAuto(ctx.Db, &res, &lr)

			}

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
