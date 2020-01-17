package errno

//	[1, 10000)
var (
	Success = &Errno{Code: New(1), Msg: "ok"}

	ErrBind   = &Errno{Code: New(10), Msg: "bind json error"}
	ErrParams = &Errno{Code: New(11), Msg: "params empty error"}
)
