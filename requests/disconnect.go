package requests

import "net/url"

// Disconnect disconnects you from the stranger.
type Disconnect struct {}

func (d Disconnect) URL() string {
	return "/disconnect"
}

func (d Disconnect) FormData() *url.Values {
	return &url.Values{}
}
