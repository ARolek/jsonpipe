package jsonpipe

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"runtime/debug"
)

const (
	MaxScanTokenSize = 64 * 1024
)

type Server struct {
	Conn     net.Conn
	Reader   *bufio.Reader
	Encoder  *json.Encoder
	LineData []byte
}

func (server *Server) Read() {
	for {
		var err error
		lineData, isPrefix, err := server.Reader.ReadLine()
		if err == io.EOF {
			log.Println("server disconnected: " + server.Conn.RemoteAddr().String())
			server.Conn.Close()
			break
		}

		if err != nil {
			log.Println("reader error: ", err)
			continue
		}

		if isPrefix {
			server.LineData = append(server.LineData, lineData...)

			//	check if the request is larger than our max allowed size.
			//	if so, we are probably being flooded, kill the connection
			if len(server.LineData) > MaxScanTokenSize {
				log.Println("connection flood detected. closing connection")
				server.Conn.Close()
				break
			}
			continue
		}

		server.LineData = append(server.LineData, lineData...)
		server.Decode()

		//	reset our line slice
		server.LineData = []byte{}
	}

	//	end go routine
	return
}

func (server *Server) Decode() {
	var err error
	var req Request
	err = json.Unmarshal(server.LineData, &req)
	if err != nil {
		log.Println(err)

		res := Response{
			ReqId:   "error",
			Success: false,
			Error:   "JSON decode fail: " + err.Error(),
		}

		err = server.Encoder.Encode(res)
		if err != nil {
			log.Println(err)
		}
		return
	}

	fmt.Printf("%+v \n", req)

	//	recover if our go routine crashes
	defer func() {
		if r := recover(); r != nil {
			log.Println(debug.Stack())
		}
	}()

	//	log.Println(runtime.NumGoroutine())

	//	make request concurrent
	go server.HandleRequest(&req)
}

func (server *Server) HandleRequest(req *Request) {
	var err error
	var res Response

	if req.Data == nil {
		res = Response{
			ReqId:   req.ReqId,
			Success: false,
			Error:   "no request data provided",
		}
	} else {
		//		data, err := req.HandleAction()
		data, err := Pipe.actions[req.Action].h(req.Data)
		if err != nil {
			res = Response{
				ReqId:   req.ReqId,
				Success: false,
				Error:   err.Error(),
			}
		} else {
			//	success!
			res = Response{
				ReqId:   req.ReqId,
				Success: true,
				Data:    data,
			}
		}
	}

	log.Printf("Res: %+v \n", res)

	err = server.Encoder.Encode(&res)
	if err != nil {
		log.Println("Encoder error: ", err)
	}

	//	TODO: remove request data

	//	end go routine
	return
}
