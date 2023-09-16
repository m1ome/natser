package natser

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

const (
	errHeader = "natser-err"
)

type (
	Server struct {
		nc            *nats.Conn
		handlers      map[string]Handler
		subscriptions []*nats.Subscription
		hmu           sync.RWMutex
	}

	Handler func(req *Request) error
)

// New creates a new Server instance, with a connection to nats
// you should provide url in format `localhost:4222`.
func New(url string) (*Server, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("error connecting to nats: %w", err)
	}

	s := &Server{
		nc:            nc,
		handlers:      make(map[string]Handler),
		subscriptions: make([]*nats.Subscription, 0),
	}

	return s, nil
}

// MakeRequest submits a request and waits for response
func (s *Server) MakeRequest(method string, body interface{}, v interface{}) error {
	reqData, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error encoding request data: %w", err)
	}

	req, err := s.nc.Request(method, reqData, time.Second)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}

	if err := req.Header.Get(errHeader); err != "" {
		return errors.New(err)
	}

	return json.Unmarshal(req.Data, v)
}

// AddHandler adds new handler to server with
func (s *Server) AddHandler(method string, fn Handler) {
	s.hmu.Lock()
	s.handlers[method] = fn
	s.hmu.Unlock()
}

// Serve is starting subscriptions for all registered handlers
func (s *Server) Serve() error {
	for method, handler := range s.handlers {
		sub, err := s.nc.Subscribe(method, func(msg *nats.Msg) {
			req := &Request{
				method: method,
				body:   msg.Data,
			}

			res := nats.NewMsg(msg.Reply)
			if err := handler(req); err != nil {
				res.Header.Add(errHeader, err.Error())
			} else {
				res.Data = req.data
			}

			s.nc.PublishMsg(res)
		})

		if err != nil {
			return fmt.Errorf("error creating nats subscription: %w", err)
		}

		s.subscriptions = append(s.subscriptions, sub)
	}

	return nil
}

// Stop unsubscribe from all subscriptions and drain all events
func (s *Server) Stop() error {
	for _, sub := range s.subscriptions {
		if err := sub.Unsubscribe(); err != nil {
			return fmt.Errorf("error unsubscribing: %w", err)
		}
	}

	return s.nc.Drain()
}
