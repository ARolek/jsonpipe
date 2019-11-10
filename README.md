# JSONpipe

Golang TCP/IP socket server for handling JSON requests with ids. Each JSON message will need to have a reqId property which will be included in the response message.

### Features:

	- Each connection has it's own go routine
	- Each request on a connection has it's own go routine
	- If an attempted socket flood is happening, server will automatically close the connection.

### Install the package:

```
go get https://github.com/ARolek/jsonpipe
```

### Use the package:

```
package main

import (
	"github.com/ARolek/jsonpipe"
)

func main (){
	jsonpipe.Handle("MyAction", myHandler)
	jsonpipe.ListenAndServe(":8080")
}

//	data var does not include the reqId or the action properties. only the request data property.
func myHandler(data *json.RawMessage) (map[string]interface{}, error){
	//	typicaly I unmarshal the JSON data into a struct here
}
```
### Request requirements:

- reqId (string): will be sent back to the client with the response
- action (string): used to match up a request with a handler
- data (map[string]interface{}): will be sent to the registered handler as *json.RawMessage. Perfect for unmarshalling into a struct 

**IMPORTANT:** requests must be on a single line. JSONpipe scans for new lines (\n) to decide when a message ends. Do not send multiline JSON data.


### Example 

Request (indented for readability. must be on a single line)

```
{
	"reqId": "0",
	"action": "AnalyzeFile", 
	"data": {
		"src": "/path/to/some/file.png"
	}
}
```

Response (success):

```
{
	"reqId": "0",
	"success": true,
	"data": {
		//response data here
	}
}
```

Response(fail):

```
{
	"reqId": "0",
	"success": false,
	"error": "Invalid JSON format"
}
```