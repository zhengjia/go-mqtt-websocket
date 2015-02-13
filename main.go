package main

import (
	"log"
	"net/http"

	"code.google.com/p/go.net/websocket"
	"github.com/zhengjia/go-mqtt-websocket/mqtt"
)

const listenAddr = "localhost:9292"

func main() {
	http.Handle("/connect", websocket.Handler(connectHandler))
	log.Println("server started")
	http.ListenAndServe(listenAddr, nil)
}

func connectHandler(conn *websocket.Conn) {
	log.Println("Connection started")
	// TODO Reject if authentication fails
	c, err := mqtt.GetClient()
	if err != nil {
		conn.Close()
		return
	}
	// TODO Reject if GetClient() returns error
	p := &mqtt.Proxy{Conn: conn, Client: c, Done: make(chan bool)}
	log.Println("Connection accepted")
	go p.Start()
	<-p.Done
	log.Println("Connection closed")
}
