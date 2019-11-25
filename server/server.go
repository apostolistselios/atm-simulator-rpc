package main

// var database = api.GetDatabase("../database.db")

import (
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/apostolistselios/atm-simulator-rpc/api"
)

func main() {
	atm := api.ATM{DB: api.GetDatabase("../database.db")}
	err := rpc.Register(&atm)
	if err != nil {
		log.Fatal(err)
	}

	rpc.HandleHTTP()
	listener, e := net.Listen("tcp", ":8080")
	if e != nil {
		log.Fatal("listen error", e)
	}

	log.Printf("server listening on :%d", 8080)
	http.Serve(listener, nil)
}
