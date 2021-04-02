package requests

import "net/url"

// SendTyping displays the typing status to the stranger.
type SendTyping struct {}

func (s SendTyping) URL() string {
	return "/typing"
}

func (s SendTyping) FormData() *url.Values {
	return &url.Values{}
}