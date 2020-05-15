package server

import (
	"fmt"
	"friend-api/routes"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	*http.Server
}

func NewServer() *Server {
	port := ":3000"
	if p := os.Getenv("PORT"); p != "" {
		port = fmt.Sprintf(":%s", p)
	}
	s := &http.Server{
		Addr:    port,
		Handler: routes.Handler(),
	}
	return &Server{s}
}

func (srv *Server) Start() {
	fmt.Println("Starting server at port", srv.Addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		panic(err)
	}
	fmt.Println("Listening on &s", srv.Addr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	fmt.Println("Exiting")
}
