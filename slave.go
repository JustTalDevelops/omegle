package omegleapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/JustTalDevelops/omegleapi/requests"
	"github.com/corpix/uarand"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net/http"
	"strings"
)

// randIdSelection is a string to create random IDs from.
const randIdSelection = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"

// Slave is a worker slave. It has it's own front manager, status, and current connection.
type Slave struct {
	frontManager *FrontManager
	randId       string
	topics       []string
	userAgent    string
	client       *http.Client
	waiting      bool
	connected    bool
	clientId     string
	handler      Handler
	front        string
	log          *logrus.Logger
}

// SlaveOptions contain options for the slave.
// If certain options are unfilled, they will be filled with defaults.
type SlaveOptions struct {
	// UserAgent is the user agent that will be used by the slave.
	UserAgent string
	// Client is the HTTP client that will be used by the slave.
	Client *http.Client
	// Logger is the Logrus logger that the slave should use.
	Logger *logrus.Logger
}

// LoginResponse is a response to the login request sent in the Start method of Client.
type LoginResponse struct {
	// Events are the new events added from Omegle.
	Events [][]string `json:"events"`
	// ClientID is the Omegle client ID.
	ClientID string `json:"clientID"`
}

// NewSlave creates a new slave.
func NewSlave(handler Handler, topics []string, opts ...SlaveOptions) *Slave {
	if len(opts) == 0 {
		opts = append(opts, SlaveOptions{})
	}

	b := make([]byte, 8)
	for i := range b {
		b[i] = randIdSelection[rand.Intn(len(randIdSelection))]
	}

	if opts[0].Client == nil {
		opts[0].Client = http.DefaultClient
	}

	if opts[0].UserAgent == "" {
		opts[0].UserAgent = uarand.GetRandom()
	}

	return &Slave{log: opts[0].Logger, frontManager: NewFrontManagerFromDefaultFronts(), randId: string(b), topics: topics, userAgent: opts[0].UserAgent, client: opts[0].Client, handler: handler, connected: false}
}

// SendRequest sends a new request to Omegle.
func (s *Slave) SendRequest(request requests.Request) error {
	if !s.connected {
		return errors.New("attempted to send a request when connection isn't alive")
	}
	data := request.FormData()
	data.Set("id", s.clientId)
	req, err := http.NewRequest("POST", s.front+"/"+request.URL(), strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Length", "0")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", s.userAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Sec-Gpc", "1")
	req.Header.Set("Origin", "https://www.omegle.com")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://www.omegle.com/")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if _, ok := request.(requests.Disconnect); ok {
		s.connected = false
		s.log.Info("You disconnected from the stranger.")
	}

	return resp.Body.Close()
}

// Alive returns true if the connection is alive.
func (s *Slave) Alive() bool {
	return s.connected || s.waiting
}

// Start starts the slave on the current thread.
func (s *Slave) Start() error {
	front, err := s.frontManager.FindFront()
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", front+"/start?caps=recaptcha2&firstevents=1&spid=&randid="+s.randId+"&topics="+serializeTopics(s.topics)+"&lang=en", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Length", "0")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", s.userAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("Sec-Gpc", "1")
	req.Header.Set("Origin", "https://www.omegle.com")
	req.Header.Set("Sec-Fetch-Site", "same-site")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Referer", "https://www.omegle.com/")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var loginResponse LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&loginResponse)
	if err != nil {
		return err
	}

	for _, e := range loginResponse.Events {
		for _, d := range e {
			if d == "waiting" {
				s.front = front
				s.waiting = true
				s.clientId = loginResponse.ClientID
				s.log.Info("Logged in, awaiting a user to speak to...")

				return NewListener(s).Start()
			}
		}
	}

	return fmt.Errorf("%v", loginResponse)
}
