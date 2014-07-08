package jsonpipe

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
)

type JSONpipe struct {
	actions map[string]Action
	port    string
}

type Action struct {
	h       Handler
	pattern string
}

type Handler func(*json.RawMessage) (map[string]interface{}, error)

// NewServeMux allocates and returns a new ServeMux.
func NewJSONpipe() *JSONpipe {
	return &JSONpipe{
		actions: make(map[string]Action),
	}
}

// our pipe instance
var Pipe = NewJSONpipe()

func Handle(action string, handler Handler) {
	if action == "" {
		panic("jsonstream: action can't be an empty string")
	}
	if handler == nil {
		panic("jsonstream: nil handler")
	}

	//	TODO: check for already created entry
	Pipe.actions[action] = Action{pattern: action, h: handler}
}

//	pass in the port to bind to. ie. :8080
func ListenAndServe(port string) {
	ln, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	Pipe.port = port
	log.Println("jsonstream listening on localhost" + port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			// handle error
			continue
		}

		server := Server{
			Conn:    conn,
			Reader:  bufio.NewReader(conn),
			Encoder: json.NewEncoder(conn),
		}

		log.Println("new connection from " + server.Conn.RemoteAddr().String())

		//	create a new go routine for each connection
		go server.Read()
	}
}
