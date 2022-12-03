package users

import "gorm.io/gorm"

// User model
type User struct {
	ID        string `json:"id" gorm:"primaryKey"`
	Name      string `json:"name"`
	HighScore uint32 `json:"highscore"`
	gorm.Model
}
