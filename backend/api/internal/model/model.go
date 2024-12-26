package model

type Tokens struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

type UserData struct {
	Id         int     `json:"id"`
	Name       string  `json:"name"`
	Room       string  `json:"room"`
	DormNumber string  `json:"dorm_number"`
	Role       string  `json:"role"`
	NewTokens  *Tokens `json:"new_tokens"`
}

type Card struct {
	Id          int      `json:"id"`
	Name        string   `json:"name"`
	Price       float64  `json:"price"`
	Photos      []string `json:"photos"`
	CreatorName string   `json:"creator_name"`
	DormNumber  string   `json:"dorm_number"`
	Room        string   `json:"room"`
	Telegram    string   `json:"telegram"`
	Description string   `json:"description"`
	Status      string   `json:"status"`
}
