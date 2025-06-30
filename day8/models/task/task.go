package taskmodel

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	UserID      int    `json:"user_id"`
	IsCompleted bool   `json:"is_completed"`
}
