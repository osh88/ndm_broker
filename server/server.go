package server

import (
	"fmt"
	"log"
	"ndm_broker/broker"
	"net/http"
	"strconv"
	"time"
)

func New(broker *broker.Broker) *Server {
	o := &Server{
		srv:    http.Server{},
		mux:    http.ServeMux{},
		broker: broker,
	}

	o.mux.HandleFunc("PUT /{queue}", o.put)
	o.mux.HandleFunc("GET /{queue}", o.get)
	o.srv.Handler = &o.mux

	return o
}

type Server struct {
	srv    http.Server
	mux    http.ServeMux
	broker *broker.Broker
}

func (o *Server) Run(port int) error {
	o.srv.Addr = ":" + strconv.Itoa(port)
	log.Printf("Listening %s", o.srv.Addr)
	return o.srv.ListenAndServe()
}

func (o *Server) put(w http.ResponseWriter, r *http.Request) {
	queueName := r.PathValue("queue")

	msg := r.URL.Query().Get("v")
	if msg == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	if err := o.broker.Put(queueName, msg); err != nil {
		log.Printf("Server.put(): %v", err)
	}
}

func (o *Server) get(w http.ResponseWriter, r *http.Request) {
	queueName := r.PathValue("queue")
	timeout, _ := strconv.Atoi(r.URL.Query().Get("timeout"))

	if timeout < 0 {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	msgs, err := o.broker.Subscribe(queueName)
	if err != nil {
		log.Printf("Server.get(): %v", err)
		return
	}

	// Ждем сообщение до конца
	if timeout == 0 {
		fmt.Fprint(w, <-msgs)
		return
	}

	// Ждем сообщение некоторое время
	tmr := time.NewTimer(time.Duration(timeout) * time.Second)
	defer tmr.Stop()

	select {
	case <-tmr.C:
		http.Error(w, "", http.StatusNotFound)
		return

	case msg := <-msgs:
		fmt.Fprint(w, msg)
	}
}
