// wserver
package main

import (
	am "contler"
	csm "contler/consumer"
	csp "contler/serviprov"
	//"github.com/pubnub/go/messaging"
	"database/sql"
	"fmt"
	"log"
	_ "mysql-master"
	"net/http"
	"os"
)

func main() {
	ctx := &am.AppContext{Status: "success"}

	//db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golive")
	db, err := sql.Open("mysql", "golive:ShdP6ScS9E4Dr2vS@tcp(202.138.229.148:3306)/golive_db")
	//db, err := sql.Open("mysql", "golive:ShdP6ScS9E4Dr2vS@tcp(localhost:3309)/golive_db")

	if err != nil {
		//log.Fatal(err)
		ctx.Status = err.Error()
		fmt.Println(ctx.Status)
		//http.Redirect(w, r, "/edit/"+title, http.StatusOK)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		//log.Fatal(err)
		ctx.Status = err.Error()
		fmt.Println(ctx.Status)

	} else {
		ctx.Db = db
		am.PubInit()
	}
	log.Println("Mulai2")

	http.HandleFunc("/ag/view", func(w http.ResponseWriter, r *http.Request) {
		loginHandler(w, r, ctx)
	})
	http.HandleFunc("/ag/spa/login", makeHandler(csp.LoginHandler, ctx))
	http.HandleFunc("/ag/ca/svoffer", makeHandler(csm.ServiceOfferHandler, ctx))
	http.HandleFunc("/ag/spa/update_order_status", makeHandler(csp.UpdateOdrStatusHandler, ctx))
	http.HandleFunc("/ag/spa/order_detail", makeHandler(csp.OrderDetailHandler, ctx))
	http.HandleFunc("/ag/ca/add_order_massage", makeHandler(csm.OrderMassageHandler, ctx))
	http.HandleFunc("/ag/ca/add_order_massage_extend", makeHandler(csm.ExtendMassageHandler, ctx))
	http.HandleFunc("/ag/ca/addOrderGlam", makeHandler(csm.OrderGlamHandler, ctx))
	http.HandleFunc("/ag/ca/getOrderStatusById", makeHandler(csm.GetOrderStatusById, ctx))
	http.HandleFunc("/ag/ca/getOrderById", makeHandler(csm.GetOrderById, ctx))
	http.HandleFunc("/ag/spa/bidOrder", makeHandler(csp.BidOrderHandler, ctx))
	http.HandleFunc("/ag/spa/orderList", makeHandler(csp.GetOrderList, ctx))
	http.HandleFunc("/ag/spa/orderHistory", makeHandler(csp.GetOrderHistory, ctx))
	http.HandleFunc("/ag/stop", handler)

	http.HandleFunc("/ca/addOrderAuto", makeHandler(csm.OrderAutoHandler, ctx))
	http.HandleFunc("/ca/cancelOrder", makeHandler(csm.CancelOrderHandler, ctx))
	http.HandleFunc("/ag/spa/bidStatus", makeHandler(csp.SetBidStatusHandler, ctx))
	http.HandleFunc("/ag/spa/updateLocation", makeHandler(csp.UpdateLocationHandler, ctx))
	http.HandleFunc("/ag/spa/cekItemGoAuto", makeHandler(csp.CekItemAutoHandler, ctx))
	http.HandleFunc("/ag/spa/pulsa", makeHandler(csp.SPPulsaHandler, ctx))

	http.ListenAndServe(":8082", nil)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, *am.AppContext), ctx *am.AppContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//m := validPath.FindStringSubmatch(r.URL.Path)
		//if m == nil {
		//	http.NotFound(w, r)
		//	return
		//}

		apik := r.Header.Get("apik")

		if apik == "1234abcd" {
			ctx.Status = "success"

		} else {

			ctx.Status = "Not valid API key"
		}
		fn(w, r, ctx)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	/*
		if rest.Status != "success" {
			resp, err := json.Marshal(rest)
			if err == nil {
				fmt.Fprintf(w, bytes.NewBuffer(resp).String())
			}
		} else {
			apik := r.Header.Get("apik")
			fmt.Printf("apik = %s\n", apik)

			if apik == "1234abcd" {
				switch {
				case r.URL.Path[1:] == "apilogin":
					//process(w, r, &db)
					//var status string
					fmt.Printf("db = %v\n", db)
				}
				//
			} else {
				rest.Status = "Not valid API key"
				resp, err := json.Marshal(rest)
				if err == nil {
					fmt.Fprintf(w, bytes.NewBuffer(resp).String())
				}

			}

		}
	*/
	/*
		type Response1 struct {
			Page   int
			Fruits []string
		}

		type SrvProv struct {
			id       int
			name     string
			gender   int
			role     int
			phone    string
			profile  string
			rate     float32
			foto_url string
			service  []struct {
				tipe      string
				price     int
				price_min int
			}
			equipment []struct {
				id   int
				nama string
			}
		}

		type LoginRsp struct {
			status      string
			app_version string
			sp          SrvProv
		}

		res1D := &Response1{
			Page:   1,
			Fruits: []string{"apple", "peach", "pear"}}
		res1B, _ := json.Marshal(res1D)
		fmt.Println(string(res1B))

		fmt.Fprintf(w, "Hi there, I love %s!\n", r.URL.Path[1:])

		if r.Body != nil {
			//var p []byte

			buf := new(bytes.Buffer)
			buf.ReadFrom(r.Body)
			//s := buf.String() // Does a complete copy of the bytes in the buffer.
			//fmt.Fprintf(w, "n = %s", s)

			//var byt []byte
			//r.Body.Read(byt)
			res := Response1{}
			//var dat map[string]interface{}

			//if err := json.Unmarshal(buf.Bytes(), &dat); err != nil {
			if err := json.Unmarshal(buf.Bytes(), &res); err != nil {
				panic(err)
			}
			fmt.Println(res)
			fmt.Println(res.Page)
		}
	*/

	//rest.Status = "success"

	//if r.URL.Path[1:] == "end" {
	//	//dbs.Main2()
	os.Exit(0)
	//}

}
