package jsonpipe

import ()

type Response struct {
	ReqId   string                 `json:"reqId"`
	Success bool                   `json:"success"`
	Error   string                 `json:"error,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
}
