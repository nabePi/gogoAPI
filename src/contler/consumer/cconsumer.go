package consumer

import (
	"bytes"
	co "contler"
	ds "dbs/consumer"
	"encoding/json"
	"fmt"
	"net/http"
)

func CancelOrderHandler(w http.ResponseWriter, r *http.Request, ctx *co.AppContext) {
	var (
		lr ds.CancelOrderRespon
	)
	lr.Status = ctx.Status
	//lr.App_version = "1"

	w.Header().Set("Content-Type", "application/json")

	if ctx.Status == "success" {
		if r.Body != nil {
			//var p []byte

			buf := new(bytes.Buffer)
			buf.ReadFrom(r.Body)

			res := ds.CancelOrderRequest{}

			if err := json.Unmarshal(buf.Bytes(), &res); err != nil {
				//			panic(err)
				s := err.Error()
				lr.Error = 1
				lr.Status = s

			} else {
				ds.CancelOrder(ctx.Db, &res, &lr)
			}
		} else {
			lr.Status = "Error parameter in body required"
		}
	}

	resp, err := json.Marshal(lr)

	if err == nil {
		fmt.Fprintf(w, bytes.NewBuffer(resp).String())
	} else {
		fmt.Fprintf(w, "{ \"Status\" : \""+err.Error()+"\"}")
	}
}


