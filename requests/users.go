package requests

// User request format for user
type User struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty" form:"highscore" binding:"required"`
	HighScore uint32 `json:"highscore,omitempty" form:"highscore" binding:"required"`
}
