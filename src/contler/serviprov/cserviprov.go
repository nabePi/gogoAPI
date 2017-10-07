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

func SetBidStatusHandler(w http.ResponseWriter, r *http.Request, ctx *co.AppContext) {

	var (
		//sp common.SrvProv
		lr ds.BidStatusRespon
	)
	lr.Status = ctx.Status

	w.Header().Set("Content-Type", "application/json")


	if ctx.Status == "success" {

		if r.Body != nil {
			buf := new(bytes.Buffer)
			buf.ReadFrom(r.Body)

			res := ds.BidStatusRequest{}

			if err := json.Unmarshal(buf.Bytes(), &res); err != nil {
				//			panic(err)
				s := err.Error()
				lr.Error = 1
				lr.Status = s

			} else {
				ds.SetBidStatus(ctx.Db, &res, &lr)
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

func UpdateLocationHandler(w http.ResponseWriter, r *http.Request, ctx *co.AppContext) {

	var (
		//sp common.SrvProv
		lr ds.UpdateLocationRespon
	)
	lr.Status = ctx.Status

	w.Header().Set("Content-Type", "application/json")


	if ctx.Status == "success" {

		if r.Body != nil {
			buf := new(bytes.Buffer)
			buf.ReadFrom(r.Body)

			res := ds.UpdateLocationRequest{}

			if err := json.Unmarshal(buf.Bytes(), &res); err != nil {
				//			panic(err)
				s := err.Error()
				lr.Error = 1
				lr.Status = s

			} else {
				ds.UpdateLocation(ctx.Db, &res, &lr)
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

func SPPulsaHandler(w http.ResponseWriter, r *http.Request, ctx *co.AppContext) {

	var (
		//sp common.SrvProv
		lr ds.SPPulsaRespon
	)
	lr.Status = ctx.Status

	w.Header().Set("Content-Type", "application/json")


	if ctx.Status == "success" {

		if r.Body != nil {
			buf := new(bytes.Buffer)
			buf.ReadFrom(r.Body)

			res := ds.SPPulsaRequest{}

			if err := json.Unmarshal(buf.Bytes(), &res); err != nil {
				//			panic(err)
				s := err.Error()
				lr.Error = 1
				lr.Status = s

			} else {
				ds.SPCheckPulsa(ctx.Db, &res, &lr)
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