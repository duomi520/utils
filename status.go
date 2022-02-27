package utils

//状态码
const (
	StatusUnknown = iota + 1000
	StatusNil
	StatusPing
	StatusPong
	StatusGoaway
	StatusRequest
	StatusResponse
	StatusSubscribe
	StatusUnsubscribe
	StatusStream
	StatusBroadcast
	StatusError
	StatusEvent
	StateBreakerClosed
	StateBreakerOpen
	StateBreakerHalfOpen
)

//状态码
const (
	StatusUnknown16 uint16 = iota + 1000
	StatusNil16
	StatusPing16
	StatusPong16
	StatusGoaway16
	StatusRequest16
	StatusResponse16
	StatusSubscribe16
	StatusUnsubscribe16
	StatusStream16
	StatusBroadcast16
	StatusError16
	StatusEvent16
	StateBreakerClosed16
	StateBreakerOpen16
	StateBreakerHalfOpen16
)

// HTTP状态码，参见RFC 2616
const (
	StatusContinue                     = 100
	StatusSwitchingProtocols           = 101
	StatusOK                           = 200
	StatusCreated                      = 201
	StatusAccepted                     = 202
	StatusNonAuthoritativeInfo         = 203
	StatusNoContent                    = 204
	StatusResetContent                 = 205
	StatusPartialContent               = 206
	StatusMultipleChoices              = 300
	StatusMovedPermanently             = 301
	StatusFound                        = 302
	StatusSeeOther                     = 303
	StatusNotModified                  = 304
	StatusUseProxy                     = 305
	StatusTemporaryRedirect            = 307
	StatusBadRequest                   = 400
	StatusUnauthorized                 = 401
	StatusPaymentRequired              = 402
	StatusForbidden                    = 403
	StatusNotFound                     = 404
	StatusMethodNotAllowed             = 405
	StatusNotAcceptable                = 406
	StatusProxyAuthRequired            = 407
	StatusRequestTimeout               = 408
	StatusConflict                     = 409
	StatusGone                         = 410
	StatusLengthRequired               = 411
	StatusPreconditionFailed           = 412
	StatusRequestEntityTooLarge        = 413
	StatusRequestURITooLong            = 414
	StatusUnsupportedMediaType         = 415
	StatusRequestedRangeNotSatisfiable = 416
	StatusExpectationFailed            = 417
	StatusTeapot                       = 418
	StatusInternalServerError          = 500
	StatusNotImplemented               = 501
	StatusBadGateway                   = 502
	StatusServiceUnavailable           = 503
	StatusGatewayTimeout               = 504
	StatusHTTPVersionNotSupported      = 505
)
const (
	StatusContinue16                     = 100
	StatusSwitchingProtocols16           = 101
	StatusOK16                           = 200
	StatusCreated16                      = 201
	StatusAccepted16                     = 202
	StatusNonAuthoritativeInfo16         = 203
	StatusNoContent16                    = 204
	StatusResetContent16                 = 205
	StatusPartialContent16               = 206
	StatusMultipleChoices16              = 300
	StatusMovedPermanently16             = 301
	StatusFound16                        = 302
	StatusSeeOther16                     = 303
	StatusNotModified16                  = 304
	StatusUseProxy16                     = 305
	StatusTemporaryRedirect16            = 307
	StatusBadRequest16                   = 400
	StatusUnauthorized16                 = 401
	StatusPaymentRequired16              = 402
	StatusForbidden16                    = 403
	StatusNotFound16                     = 404
	StatusMethodNotAllowed16             = 405
	StatusNotAcceptable16                = 406
	StatusProxyAuthRequired16            = 407
	StatusRequestTimeout16               = 408
	StatusConflict16                     = 409
	StatusGone16                         = 410
	StatusLengthRequired16               = 411
	StatusPreconditionFailed16           = 412
	StatusRequestEntityTooLarge16        = 413
	StatusRequestURITooLong16            = 414
	StatusUnsupportedMediaType16         = 415
	StatusRequestedRangeNotSatisfiable16 = 416
	StatusExpectationFailed16            = 417
	StatusTeapot16                       = 418
	StatusInternalServerError16          = 500
	StatusNotImplemented16               = 501
	StatusBadGateway16                   = 502
	StatusServiceUnavailable16           = 503
	StatusGatewayTimeout16               = 504
	StatusHTTPVersionNotSupported16      = 505
)
