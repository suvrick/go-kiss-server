package model

type KissResponse struct {
	Status int           `json:"status"`
	Code   int           `json:"code"`
	Data   []interface{} `json:"data"`
	Delay  int           `json:"delay"`
}

// KissUser ...
type KissUser struct {
	UserID    int    `json:"user_id" gorm:"unique ,primaryKey"`
	IsTrial   bool   `json:"is_trial"`
	KissCount int    `json:"kiss_count"`
	DateUse   string `json:"date"`
}

// KissState ...
type KissState struct {
	IP    string `json:"ip"`
	Count int    `json:"count"`
	Date  string `json:"date"`
}
