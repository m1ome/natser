package natser

import (
	"encoding/json"
)

type (
	Request struct {
		method string
		body   []byte
		data   []byte
	}
)

func (r *Request) Method() string {
	return r.method
}

func (r *Request) Unmarshal(v interface{}) error {
	return json.Unmarshal(r.body, v)
}

func (r *Request) SendResponse(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	r.data = data
	return nil
}
