package model

// User ...
type User struct {
	ID        int    `json:"id"`
	Email     string `json:"email,omitempty" gorm:"unique"`
	Password  string `json:"password,omitempty"`
	Token     string `json:"token"`
	Limit     int    `json:"limit"`
	Date      string `json:"time"`
	BotsCount int    `json:"bots"`
	Role      string `json:"role"`
}
