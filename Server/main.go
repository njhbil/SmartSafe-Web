package main

import (
	api "SmartSafe/api"

	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/api/register", api.Register)
	http.HandleFunc("/api/accounts", api.Accounts)
	http.HandleFunc("/api/login", api.Login)
	http.HandleFunc("/api/forgetpassword", api.ForgetPassword)
	http.HandleFunc("/api/resetpassword", api.ResetPassword)
	log.Println("Server started at", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
