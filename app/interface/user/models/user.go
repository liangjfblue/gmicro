package models

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Age      int32  `json:"age"`
	Addr     string `json:"addr"`
}

type RegisterRespond struct {
	Code int32  `json:"code"`
	Uid  string `json:"uid"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRespond struct {
	Code  int32  `json:"code"`
	Token string `json:"token"`
}

type InfoRequest struct {
	Uid string `json:"uid"`
}

type InfoRespond struct {
	Username string `json:"username"`
	Age      int32  `json:"age"`
	Addr     string `json:"addr"`
}

type AuthRequest struct {
	Token string `json:"token"`
}

type AuthResponse struct {
	UID string `json:"uid"`
}
