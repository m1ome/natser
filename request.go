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

func (r *Request) Parse(v interface{}) error {
	return json.Unmarshal(r.body, &v)
}

func (r *Request) Json(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}

	r.data = data
	return nil
}
