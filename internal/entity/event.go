package entity

type Event struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	Active    bool   `json:"active"`
	Title     string `json:"name"`
	Daily     bool   `json:"daily"`
	IsSent    bool   `json:"is_sent"`
	Subject   string `json:"subject"`
	Author    string `json:"author"`
	Letter    Letter `json:"letter"`
	SendAt    int64  `json:"send_at"`
	SentAt    int64  `json:"sent_at"`
	CreatedAt int64  `json:"created_at"`
}
