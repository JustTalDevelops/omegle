package requests

import "net/url"

// Request is an Omegle request from the slave to Omegle.
// It is used for things such as disconnections, messages, and typing.
type Request interface {
	// URL is the address that the request should be sent to.
	URL() string
	// FormData is the form data that the request should be sent to.
	FormData() *url.Values
}