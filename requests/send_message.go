package requests

import "net/url"

// SendMessage is a request used when attempting to send a message to the stranger.
type SendMessage struct {
	// Content is the content of the message being sent.
	Content string
}

func (s SendMessage) URL() string {
	return "/send"
}

func (s SendMessage) FormData() (v *url.Values) {
	v = &url.Values{}
	v.Set("msg", s.Content)

	return
}