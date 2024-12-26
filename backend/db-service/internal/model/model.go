package model

type Card struct {
	Id          int      `json:"id"`
	Name        string   `json:"name"`
	Price       float64  `json:"price"`
	Photos      []string `json:"photos"`
	CreatorName string   `json:"creator_name"`
	DormNumber  string   `json:"dorm_number"`
	Description string   `json:"description"`
	Room        string   `json:"room"`
	Telegram    string   `json:"telegram"`
	Status      string   `json:"status"`
}

type UserData struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Room       string `json:"room"`
	DormNumber string `json:"dorm_number"`
}
