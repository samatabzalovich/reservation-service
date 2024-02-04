package main

type RequestPayload struct {
	Action      string       `json:"action"`
	Auth        AuthPayload  `json:"auth,omitempty"`
	Reg         RegPayload   `json:"reg,omitempty"`
	Token       TokenPayload `json:"token,omitempty"`
	Sms         SmsPayload   `json:"sms,omitempty"`
	Institution InstPayload  `json:"institution,omitempty"`
}

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

type InstPayload struct {
	Id           int64       `json:"id,omitempty"`
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	Website      string      `json:"website"`
	OwnerId      int64       `json:"owner_id"`
	Latitude     string      `json:"latitude"`
	Longitude    string      `json:"longitude"`
	Address      string      `json:"address"`
	Phone        string      `json:"phone"`
	Country      string      `json:"country"`
	City         int32       `json:"city"`
	Category     []int64     `json:"categories"`
	WorkingHours []workHours `json:"workingHours"`
}

type FilterPayload struct {
	PageSize   int
	Page       int
	SearchText string
	Sort       string
	Categories []int64
}

type workHours struct {
	Day   int    `json:"day"`
	Open  string `json:"open"`
	Close string `json:"close"`
}
