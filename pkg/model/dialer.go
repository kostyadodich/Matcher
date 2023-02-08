package model

type Dialer struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	PhoneNumber int    `json:"phone_number"`
	Items       []Item `json:"items"`
}
