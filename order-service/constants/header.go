package constants

import "net/textproto"

var (
	XserviceName  = textproto.CanonicalMIMEHeaderKey("x-service-name")
	XApiKey       = textproto.CanonicalMIMEHeaderKey("x-api-key")
	XrequestAt    = textproto.CanonicalMIMEHeaderKey("x-request-at")
	Authorization = textproto.CanonicalMIMEHeaderKey("authorization")
)
