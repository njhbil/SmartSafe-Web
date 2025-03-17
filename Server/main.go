package main

import (
	api "SmartSafe/api"

	"flag"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")

	reader(ws)
}

func setupRoutes() {
	http.HandleFunc("/ws", wsEndpoint)
	http.HandleFunc("/api/register", api.Register)
	http.HandleFunc("/api/accounts", api.Accounts)
	http.HandleFunc("/api/login", api.Login)
	http.HandleFunc("/api/forgetpassword", api.ForgetPassword)
	http.HandleFunc("/api/resetpassword", api.ResetPassword)
	http.HandleFunc("/api/refreshtoken", api.RefreshToken)
	http.HandleFunc("/api/verifytoken", api.VerifyToken)
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	setupRoutes()
	log.Println("Server started at", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
