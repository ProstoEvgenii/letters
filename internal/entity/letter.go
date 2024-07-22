package entity

type Letter struct {
	ID      int64  `json:"id"`
	UserID  int64  `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
