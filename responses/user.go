package responses

// User response format for user
type User struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name,omitempty"`
	HighScore uint32 `json:"highscore,omitempty"`
}
