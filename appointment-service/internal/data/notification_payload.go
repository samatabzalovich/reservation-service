package data
type RequestBody struct {
	Token        string   `json:"token"`
	Notification *Message `json:"notification"`
	Data         *Message `json:"data"`
}

type Message struct {
	Title    string `json:"title"`
	Body     string `json:"body"`
	ImageUrl string `json:"imageUrl"`
}