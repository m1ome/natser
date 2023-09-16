package natser

import (
	"errors"
	"testing"
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

func TestNatser(t *testing.T) {
	t.Run("correct workflow", func(t *testing.T) {
		server, err := New("0.0.0.0:4222")
		if err != nil {
			t.Fatalf("error connecting to nats: %v", err)
		}

		server.AddHandler("ping", func(r *Request) error {
			var req Req
			if err := r.Unmarshal(&req); err != nil {
				return err
			}

			res := Res{Verified: req.Age >= 18}
			return r.SendResponse(res)
		})

		if err := server.Serve(); err != nil {
			t.Fatalf("error serving: %v", err)
		}

		req := Req{Name: "John Doe", Age: 23}
		var res Res
		if err := server.MakeRequest("ping", req, &res); err != nil {
			t.Fatalf("error making request: %v", err)
		}

		if !res.Verified {
			t.Fatalf("error on response, data missmatch")
		}

		if err := server.Stop(); err != nil {
			t.Fatalf("error on stopping server: %v", err)
		}
	})

	t.Run("bad workflow", func(t *testing.T) {
		server, err := New("0.0.0.0:4222")
		if err != nil {
			t.Fatalf("error connecting to nats: %v", err)
		}

		server.AddHandler("ping", func(r *Request) error {
			return errors.New("i am an error")
		})

		if err := server.Serve(); err != nil {
			t.Fatalf("error serving: %v", err)
		}

		req := Req{Name: "John Doe", Age: 23}
		var res Res
		err = server.MakeRequest("ping", req, &res)
		if err == nil || err.Error() != "i am an error" {
			t.Fatalf("wanna error on response, got wrong one: %v", err)
		}

		if err := server.Stop(); err != nil {
			t.Fatalf("error on stopping server: %v", err)
		}
	})
}
