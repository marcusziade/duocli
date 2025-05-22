package models

import (
	"time"
)

// User represents a learner
type User struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Name      string    `json:"name"`
	Level     int       `json:"level" gorm:"default:1"`
	XP        int       `json:"xp" gorm:"default:0"`
	Streak    int       `json:"streak" gorm:"default:0"`
	LastSeen  time.Time `json:"last_seen"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Lesson represents a language lesson
type Lesson struct {
	ID          uint   `gorm:"primarykey" json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Level       int    `json:"level"`
	Order       int    `json:"order"`
	XPReward    int    `json:"xp_reward" gorm:"default:10"`
	IsCompleted bool   `json:"is_completed" gorm:"default:false"`
	Language    string `json:"language" gorm:"default:german"`
}

// Exercise represents individual exercises within lessons
type Exercise struct {
	ID           uint   `gorm:"primarykey" json:"id"`
	LessonID     uint   `json:"lesson_id"`
	Type         string `json:"type"` // translation, multiple_choice, fill_blank, speaking
	Question     string `json:"question"`
	Answer       string `json:"answer"`
	Options      string `json:"options"` // JSON string for multiple choice
	Hint         string `json:"hint"`
	Explanation  string `json:"explanation"`
	Order        int    `json:"order"`
	Difficulty   int    `json:"difficulty" gorm:"default:1"`
	Lesson       Lesson `gorm:"foreignKey:LessonID"`
}

// Progress tracks user progress
type Progress struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	UserID      uint      `json:"user_id"`
	LessonID    uint      `json:"lesson_id"`
	ExerciseID  uint      `json:"exercise_id"`
	IsCorrect   bool      `json:"is_correct"`
	Attempts    int       `json:"attempts" gorm:"default:1"`
	CompletedAt time.Time `json:"completed_at"`
	User        User      `gorm:"foreignKey:UserID"`
	Lesson      Lesson    `gorm:"foreignKey:LessonID"`
	Exercise    Exercise  `gorm:"foreignKey:ExerciseID"`
}

// Vocabulary represents words to learn
type Vocabulary struct {
	ID          uint   `gorm:"primarykey" json:"id"`
	German      string `json:"german"`
	English     string `json:"english"`
	Category    string `json:"category"`
	Difficulty  int    `json:"difficulty" gorm:"default:1"`
	AudioURL    string `json:"audio_url"`
	Example     string `json:"example"`
	Translation string `json:"translation"`
}