package xcode

var (
	OK                 = NewCode(0, "OK")
	Nologin            = NewCode(101, "NOT_LOGIN")
	RequestErr         = NewCode(400, "INVALID_ARGUMENT")
	Unauthorized       = NewCode(401, "UNAUTHENTICATED")
	AccessDenied       = NewCode(403, "PERMISSION_DENIED")
	NotFound           = NewCode(404, "NOT_FOUND")
	MethodNotAllowed   = NewCode(405, "METHOD_NOT_ALLOWED")
	Canceled           = NewCode(498, "CANCELED")
	ServerErr          = NewCode(500, "INTERNAL_ERROR")
	ServiceUnavailable = NewCode(503, "UNAVAILABLE")
	Deadline           = NewCode(504, "DEADLINE_EXCEEDED")
	LimitExceed        = NewCode(509, "RESOURCE_EXHAUSTED")
	FindByMobileErr    = NewCode(999, "db_err")
)
