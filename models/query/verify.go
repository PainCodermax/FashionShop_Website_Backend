package query

type VerifyRequest struct {
	User_Code  string `json:"userCode,omitempty"`
	VerifyCode string `json:"verifyCode,omitempty"`
}

type UpdatePasswordRequest struct {
	Email      string `json:"email,omitempty"`
	VerifyCode string `json:"verifyCode,omitempty"`
	Password   string `json:"password,omitempty"`
}
