package requests

import "net/url"

// StopTyping stops displays the typing status to the stranger.
type StopTyping struct {}

func (s StopTyping) URL() string {
	return "/stoppedtyping"
}

func (s StopTyping) FormData() *url.Values {
	return &url.Values{}
}