package main

type TokenPayload struct {
	Bearer string `json:"bearer"`
}

// AuthPayload is the embedded type (in RequestPayload) that describes an authentication request
type AuthPayload struct {
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
}
type SmsPayload struct {
	PhoneNumber string `json:"phoneNumber"`
	Code        string `json:"code"`
}
type RegPayload struct {
	UserName    string `json:"userName"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
	Type        string `json:"type"`
}
