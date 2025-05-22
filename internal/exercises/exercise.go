package exercises

import (
	"bufio"
	"duocli/internal/database"
	"duocli/internal/models"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

type ExerciseSession struct {
	UserID   uint
	LessonID uint
	Score    int
	Total    int
	XPEarned int
}

func StartLesson(userID, lessonID uint) error {
	// Get lesson details
	var lesson models.Lesson
	if err := database.DB.First(&lesson, lessonID).Error; err != nil {
		return fmt.Errorf("lesson not found: %w", err)
	}

	// Get exercises for this lesson
	var exercises []models.Exercise
	if err := database.DB.Where("lesson_id = ?", lessonID).Order("\"order\"").Find(&exercises).Error; err != nil {
		return fmt.Errorf("failed to load exercises: %w", err)
	}

	if len(exercises) == 0 {
		return fmt.Errorf("no exercises found for this lesson")
	}

	session := &ExerciseSession{
		UserID:   userID,
		LessonID: lessonID,
		Score:    0,
		Total:    len(exercises),
		XPEarned: 0,
	}

	color.Cyan("\n🎓 Starting Lesson: %s", lesson.Title)
	color.White("📝 %s", lesson.Description)
	color.Yellow("💪 %d exercises to complete\n", len(exercises))

	// Run through exercises
	for i, exercise := range exercises {
		color.Blue("\n📚 Exercise %d/%d", i+1, len(exercises))
		
		correct := runExercise(exercise)
		if correct {
			session.Score++
			session.XPEarned += 5
			color.Green("✅ Correct! (+5 XP)")
		} else {
			color.Red("❌ Incorrect. The correct answer was: %s", exercise.Answer)
			if exercise.Explanation != "" {
				color.Yellow("💡 %s", exercise.Explanation)
			}
		}

		// Record progress
		progress := models.Progress{
			UserID:      userID,
			LessonID:    lessonID,
			ExerciseID:  exercise.ID,
			IsCorrect:   correct,
			Attempts:    1,
			CompletedAt: time.Now(),
		}
		database.DB.Create(&progress)

		// Small delay for better UX
		time.Sleep(1 * time.Second)
	}

	// Calculate lesson completion
	completionPercentage := float64(session.Score) / float64(session.Total) * 100
	
	// Bonus XP for high performance
	if completionPercentage >= 80 {
		bonusXP := lesson.XPReward
		session.XPEarned += bonusXP
		color.Green("🎉 Great job! Bonus XP: +%d", bonusXP)
	}

	// Update user stats
	var user models.User
	database.DB.First(&user, userID)
	user.XP += session.XPEarned
	user.LastSeen = time.Now()
	
	// Check for level up
	newLevel := calculateLevel(user.XP)
	if newLevel > user.Level {
		user.Level = newLevel
		color.Magenta("🚀 LEVEL UP! You are now level %d!", newLevel)
	}
	
	database.DB.Save(&user)

	// Mark lesson as completed if score is good enough
	if completionPercentage >= 70 {
		lesson.IsCompleted = true
		database.DB.Save(&lesson)
	}

	// Show results
	showResults(session, completionPercentage)
	
	return nil
}

func runExercise(exercise models.Exercise) bool {
	fmt.Printf("\n%s\n", exercise.Question)
	
	if exercise.Hint != "" {
		color.Yellow("💡 Hint: %s", exercise.Hint)
	}

	switch exercise.Type {
	case "translation":
		return handleTranslation(exercise)
	case "multiple_choice":
		return handleMultipleChoice(exercise)
	case "fill_blank":
		return handleFillBlank(exercise)
	default:
		return handleTranslation(exercise)
	}
}

func handleTranslation(exercise models.Exercise) bool {
	fmt.Print("Your answer: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	answer := strings.TrimSpace(scanner.Text())
	
	return strings.EqualFold(answer, exercise.Answer)
}

func handleMultipleChoice(exercise models.Exercise) bool {
	var options []string
	json.Unmarshal([]byte(exercise.Options), &options)
	
	// Shuffle options
	rand.Shuffle(len(options), func(i, j int) {
		options[i], options[j] = options[j], options[i]
	})
	
	fmt.Println("\nChoose the correct answer:")
	for i, option := range options {
		fmt.Printf("%d. %s\n", i+1, option)
	}
	
	fmt.Print("Your choice (1-4): ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	choice := strings.TrimSpace(scanner.Text())
	
	choiceNum, err := strconv.Atoi(choice)
	if err != nil || choiceNum < 1 || choiceNum > len(options) {
		return false
	}
	
	selectedAnswer := options[choiceNum-1]
	return strings.EqualFold(selectedAnswer, exercise.Answer)
}

func handleFillBlank(exercise models.Exercise) bool {
	fmt.Print("Fill in the blank: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	answer := strings.TrimSpace(scanner.Text())
	
	return strings.EqualFold(answer, exercise.Answer)
}

func showResults(session *ExerciseSession, percentage float64) {
	fmt.Println("\n" + strings.Repeat("=", 50))
	color.Cyan("📊 LESSON COMPLETE!")
	fmt.Println(strings.Repeat("=", 50))
	
	color.White("Score: %d/%d (%.1f%%)", session.Score, session.Total, percentage)
	color.Green("XP Earned: +%d", session.XPEarned)
	
	if percentage >= 90 {
		color.Magenta("🏆 PERFECT! Outstanding work!")
	} else if percentage >= 80 {
		color.Green("🌟 EXCELLENT! Great job!")
	} else if percentage >= 70 {
		color.Yellow("👍 GOOD! Lesson completed!")
	} else {
		color.Red("📚 Keep practicing! You can retake this lesson.")
	}
	
	fmt.Println(strings.Repeat("=", 50))
}

func calculateLevel(xp int) int {
	// Simple level calculation: every 100 XP = 1 level
	return (xp / 100) + 1
}