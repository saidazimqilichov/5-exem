package models

type RegisterReq struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Age      int32  `json:"age"`
}

type RegisterResp struct {
	Message string `json:"message"`
}

type VerifyReq struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type VerifyResp struct {
	UserID  int32  `json:"user_id"`
	Token   string `json:"token"`
	Message string `json:"message"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResp struct {
	UserID  int32  `json:"user_id"`
	Token   string `json:"token"`
	Message string `json:"message"`
}

type GetUserReq struct {
	UserID int32 `json:"user_id"`
}

type GetUserResp struct {
	UserID   int32         `json:"user_id"`
	Response []RegisterReq `json:"response"`
}

type DeleteUserReq struct {
	UserID int32 `json:"user_id"`
}

type DeleteUserResp struct {
	Message string `json:"message"`
}
