# Natser
> Simple and handy module for making rpc calls for NATS

## Installation
```bash
go get -u github.com/m1ome/natser
```

## Usage
```golang
package main

import (
	"fmt"
	"log"

	"github.com/m1ome/natser"
)

type (
	Req struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	Res struct {
		Verified bool `json:"verified"`
	}
)

func main() {
	server, err := natser.New("0.0.0.0:4222")
	if err != nil {
		log.Fatalf("error connecting to nats: %v", err)
	}

	server.AddHandler("ping", func(r *natser.Request) error {
		var req Req
		if err := r.Unmarshal(&req); err != nil {
			return err
		}

		res := Res{Verified: req.Age >= 18}
		return r.SendResponse(res)
	})

	if err := server.Serve(); err != nil {
		log.Fatalf("error serving: %v", err)
	}

	req := Req{Name: "John Doe", Age: 23}
	var res Res
	if err := server.MakeRequest("ping", req, &res); err != nil {
		log.Fatalf("error making request: %v", err)
	}

	if !res.Verified {
		log.Fatalf("error on response, data missmatch")
	}

	if err := server.Stop(); err != nil {
		log.Fatalf("error on stopping server: %v", err)
	}

	fmt.Printf("Verified: %t", res.Verified)
}

```