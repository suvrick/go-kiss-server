package until

import (
	"encoding/json"
	"log"
	"net/http"
)

// Response ...
type Response struct {
	Code  int         `json:code`
	Error string      `json:error`
	Data  interface{} `json:data`
}

// WriteResponse ...
func WriteResponse(w http.ResponseWriter, r *http.Request, code int, data interface{}, err error) {
	var msgErr string
	if err != nil {
		msgErr = err.Error()
	}

	res := &Response{
		Code:  code,
		Data:  data,
		Error: msgErr,
	}

	w.WriteHeader(res.Code)

	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Println(err.Error())
	}
}

// JSONBind ...
func JSONBind(r *http.Request, strct interface{}) error {
	return json.NewDecoder(r.Body).Decode(&strct)
}
