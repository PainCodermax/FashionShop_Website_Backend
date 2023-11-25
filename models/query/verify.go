package query

type VerifyRequest struct {
	User_Code  string `json:"userCode,omitempty"`
	VerifyCode string `json:"verifyCode,omitempty"`
}
