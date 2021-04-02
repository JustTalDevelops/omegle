package omegleapi

import (
	"errors"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Listener is an event listener for an Omegle slave.
type Listener struct {
	slave *Slave
}

// NewListener creates a new listener for a slave.
func NewListener(slave *Slave) *Listener {
	return &Listener{slave}
}

// Start starts the listener.
func (l *Listener) Start() error {
	for {
		if !l.slave.connected && !l.slave.waiting {
			return nil
		}
		v := &url.Values{}
		v.Set("id", l.slave.clientId)
		req, err := http.NewRequest("POST", l.slave.front+"/events", strings.NewReader(v.Encode()))
		if err != nil {
			return err
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Set("User-Agent", l.slave.userAgent)
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
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		for _, event := range gjson.ParseBytes(b).Array() {
			eventName := event.Array()[0].String()
			eventData := event.Array()[1:]
			switch eventName {
			case "connected":
				l.slave.log.Info("Connected to a user.")
				l.slave.waiting = false
				l.slave.connected = true
				l.slave.handler.OnConnected(l.slave)
			case "commonLikes":
				var commonLikes []string
				for _, like := range eventData[0].Array() {
					commonLikes = append(commonLikes, like.String())
				}
				l.slave.log.Debugf("Common likes: %v\n", commonLikes)
				l.slave.handler.OnCommonLikesReceived(l.slave, commonLikes)
			case "typing":
				l.slave.log.Debug("Stranger is typing...")
				l.slave.handler.OnTyping(l.slave)
			case "gotMessage":
				l.slave.log.Infof("Stranger: %v\n", eventData[0].String())
				l.slave.handler.OnMessage(l.slave, eventData[0].String())
			case "strangerDisconnected":
				l.slave.log.Info("Stranger disconnected.")
				l.slave.connected = false
				l.slave.handler.OnStrangerDisconnected(l.slave)
				return nil
			case "connectionDied":
				l.slave.log.Info("The connection died.")
				l.slave.connected = false
				l.slave.handler.OnServerDisconnected(l.slave)
				return errors.New(string(b))
			case "error":
				return errors.New(string(b))
			case "stoppedTyping":
				l.slave.log.Debug("Stranger stopped typing.")
				l.slave.handler.OnTypingStopped(l.slave)
			case "statusInfo":
				// Do nothing.
			case "identDigests":
				// Do nothing.
			default:
				l.slave.log.Infof("Unhandled event from Omegle: %v: %v\n", eventName, eventData)
			}
		}

		resp.Body.Close()

		time.Sleep(1)
	}
}
