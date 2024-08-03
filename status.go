package utils

//状态码
const (
	StatusUnknown = iota + 10
	StatusNil
	StatusPing
	StatusPong
	StatusGoaway
	StatusRequest
	StatusResponse
	StatusNotify
	StatusCtxCancelFunc
	StatusSubscribe
	StatusUnsubscribe
	StatusStream
	StatusBroadcast
	StatusError
)

//状态码
const (
	StatusUnknown16 uint16 = iota + 10
	StatusNil16
	StatusPing16
	StatusPong16
	StatusGoaway16
	StatusRequest16
	StatusResponse16
	StatusNotify16
	StatusCtxCancelFunc16
	StatusSubscribe16
	StatusUnsubscribe16
	StatusStream16
	StatusBroadcast16
	StatusError16
)

// HTTP状态码，参见RFC 2616
const (
	StatusContinue           = 100 // RFC 7231, 6.2.1
	StatusSwitchingProtocols = 101 // RFC 7231, 6.2.2
	StatusProcessing         = 102 // RFC 2518, 10.1
	StatusEarlyHints         = 103 // RFC 8297

	StatusOK                   = 200 // RFC 7231, 6.3.1
	StatusCreated              = 201 // RFC 7231, 6.3.2
	StatusAccepted             = 202 // RFC 7231, 6.3.3
	StatusNonAuthoritativeInfo = 203 // RFC 7231, 6.3.4
	StatusNoContent            = 204 // RFC 7231, 6.3.5
	StatusResetContent         = 205 // RFC 7231, 6.3.6
	StatusMultiStatus          = 207 // RFC 4918, 11.1
	StatusAlreadyReported      = 208 // RFC 5842, 7.1
	StatusIMUsed               = 226 // RFC 3229, 10.4.1

	StatusMultipleChoices   = 300 // RFC 7231, 6.4.1
	StatusMovedPermanently  = 301 // RFC 7231, 6.4.2
	StatusFound             = 302 // RFC 7231, 6.4.3
	StatusSeeOther          = 303 // RFC 7231, 6.4.4
	StatusNotModified       = 304 // RFC 7232, 4.1
	StatusUseProxy          = 305 // RFC 7231, 6.4.5
	StatusTemporaryRedirect = 307 // RFC 7231, 6.4.7
	StatusPermanentRedirect = 308 // RFC 7538, 3

	StatusBadRequest                   = 400 // RFC 7231, 6.5.1
	StatusUnauthorized                 = 401 // RFC 7235, 3.1
	StatusPaymentRequired              = 402 // RFC 7231, 6.5.2
	StatusForbidden                    = 403 // RFC 7231, 6.5.3
	StatusNotFound                     = 404 // RFC 7231, 6.5.4
	StatusMethodNotAllowed             = 405 // RFC 7231, 6.5.5
	StatusNotAcceptable                = 406 // RFC 7231, 6.5.6
	StatusProxyAuthRequired            = 407 // RFC 7235, 3.2
	StatusRequestTimeout               = 408 // RFC 7231, 6.5.7
	StatusConflict                     = 409 // RFC 7231, 6.5.8
	StatusGone                         = 410 // RFC 7231, 6.5.9
	StatusLengthRequired               = 411 // RFC 7231, 6.5.10
	StatusPreconditionFailed           = 412 // RFC 7232, 4.2
	StatusRequestEntityTooLarge        = 413 // RFC 7231, 6.5.11
	StatusRequestURITooLong            = 414 // RFC 7231, 6.5.12
	StatusUnsupportedMediaType         = 415 // RFC 7231, 6.5.13
	StatusRequestedRangeNotSatisfiable = 416 // RFC 7233, 4.4
	StatusExpectationFailed            = 417 // RFC 7231, 6.5.14
	StatusTeapot                       = 418 // RFC 7168, 2.3.3
	StatusMisdirectedRequest           = 421 // RFC 7540, 9.1.2
	StatusUnprocessableEntity          = 422 // RFC 4918, 11.2
	StatusLocked                       = 423 // RFC 4918, 11.3
	StatusFailedDependency             = 424 // RFC 4918, 11.4
	StatusUpgradeRequired              = 426 // RFC 7231, 6.5.15
	StatusPreconditionRequired         = 428 // RFC 6585, 3
	StatusTooManyRequests              = 429 // RFC 6585, 4
	StatusRequestHeaderFieldsTooLarge  = 431 // RFC 6585, 5
	StatusUnavailableForLegalReasons   = 451 // RFC 7725, 3

	StatusInternalServerError           = 500 // RFC 7231, 6.6.1
	StatusNotImplemented                = 501 // RFC 7231, 6.6.2
	StatusBadGateway                    = 502 // RFC 7231, 6.6.3
	StatusServiceUnavailable            = 503 // RFC 7231, 6.6.4
	StatusGatewayTimeout                = 504 // RFC 7231, 6.6.5
	StatusHTTPVersionNotSupported       = 505 // RFC 7231, 6.6.6
	StatusVariantAlsoNegotiates         = 506 // RFC 2295, 8.1
	StatusInsufficientStorage           = 507 // RFC 4918, 11.5
	StatusLoopDetected                  = 508 // RFC 5842, 7.2
	StatusNotExtended                   = 510 // RFC 2774, 7
	StatusNetworkAuthenticationRequired = 511 // RFC 6585, 6
)
const (
	StatusContinue16           uint16 = 100
	StatusSwitchingProtocols16 uint16 = 101
	StatusProcessing16         uint16 = 102
	StatusEarlyHints16         uint16 = 103

	StatusOK16                   uint16 = 200
	StatusCreated16              uint16 = 201
	StatusAccepted16             uint16 = 202
	StatusNonAuthoritativeInfo16 uint16 = 203
	StatusNoContent16            uint16 = 204
	StatusResetContent16         uint16 = 205
	StatusPartialContent16       uint16 = 206
	StatusMultiStatus16          uint16 = 207
	StatusAlreadyReported16      uint16 = 208
	StatusIMUsed16               uint16 = 226

	StatusMultipleChoices16   uint16 = 300
	StatusMovedPermanently16  uint16 = 301
	StatusFound16             uint16 = 302
	StatusSeeOther16          uint16 = 303
	StatusNotModified16       uint16 = 304
	StatusUseProxy16          uint16 = 305
	StatusTemporaryRedirect16 uint16 = 307
	StatusPermanentRedirect16 uint16 = 308

	StatusBadRequest16                   uint16 = 400
	StatusUnauthorized16                 uint16 = 401
	StatusPaymentRequired16              uint16 = 402
	StatusForbidden16                    uint16 = 403
	StatusNotFound16                     uint16 = 404
	StatusMethodNotAllowed16             uint16 = 405
	StatusNotAcceptable16                uint16 = 406
	StatusProxyAuthRequired16            uint16 = 407
	StatusRequestTimeout16               uint16 = 408
	StatusConflict16                     uint16 = 409
	StatusGone16                         uint16 = 410
	StatusLengthRequired16               uint16 = 411
	StatusPreconditionFailed16           uint16 = 412
	StatusRequestEntityTooLarge16        uint16 = 413
	StatusRequestURITooLong16            uint16 = 414
	StatusUnsupportedMediaType16         uint16 = 415
	StatusRequestedRangeNotSatisfiable16 uint16 = 416
	StatusExpectationFailed16            uint16 = 417
	StatusTeapot16                       uint16 = 418
	StatusMisdirectedRequest16           uint16 = 421
	StatusUnprocessableEntity16          uint16 = 422
	StatusLocked16                       uint16 = 423
	StatusFailedDependency16             uint16 = 424
	StatusUpgradeRequired16              uint16 = 426
	StatusPreconditionRequired16         uint16 = 428
	StatusTooManyRequests16              uint16 = 429
	StatusRequestHeaderFieldsTooLarge16  uint16 = 431
	StatusUnavailableForLegalReasons16   uint16 = 451

	StatusInternalServerError16           uint16 = 500
	StatusNotImplemented16                uint16 = 501
	StatusBadGateway16                    uint16 = 502
	StatusServiceUnavailable16            uint16 = 503
	StatusGatewayTimeout16                uint16 = 504
	StatusHTTPVersionNotSupported16       uint16 = 505
	StatusVariantAlsoNegotiates16         uint16 = 506
	StatusInsufficientStorage16           uint16 = 507
	StatusLoopDetected16                  uint16 = 508
	StatusNotExtended16                   uint16 = 510
	StatusNetworkAuthenticationRequired16 uint16 = 511
)
