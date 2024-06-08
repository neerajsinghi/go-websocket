package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{conns: make(map[*websocket.Conn]bool)}
}
func (s *Server) handleWSOrderBook(ws *websocket.Conn) {
	fmt.Println("new connection from client", ws.RemoteAddr(), " To order book feed")
	for {
		paylload := fmt.Sprintf("orderbook data -> %d\n", time.Now().UnixNano())
		ws.Write([]byte(paylload))
		time.Sleep(2 * time.Second)
	}

}
func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("new connection from client", ws.RemoteAddr())
	s.conns[ws] = true
	s.readLoop(ws)
}
func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			fmt.Println("receive error:", err)
			if err == io.EOF {
				delete(s.conns, ws)
				ws.Close()
				return
			}
			continue
		}
		msg := buf[:n]
		s.broadcast(msg)

	}
}
func (s *Server) broadcast(b []byte) {
	for ws := range s.conns {
		go func(ws *websocket.Conn) {

			if _, err := ws.Write(b); err != nil {
				fmt.Println("send error:", err)
			}
		}(ws)
	}
}
func main() {
	server := NewServer()
	http.Handle("/ws", websocket.Handler(server.handleWS))
	http.Handle("/orderbookfeed", websocket.Handler(server.handleWSOrderBook))
	http.ListenAndServe(":8080", nil)
}
