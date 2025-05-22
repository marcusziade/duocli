package cmd

import (
	"duocli/internal/database"
	"duocli/internal/exercises"
	"duocli/internal/models"
	"duocli/internal/ui"
	"fmt"
	"strconv"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start [lesson_id]",
	Short: "Start a specific lesson",
	Long:  `Start a specific lesson by providing its ID`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ensureUser()
		
		if len(args) == 0 {
			ui.ShowLessons()
			return
		}
		
		lessonID, err := strconv.ParseUint(args[0], 10, 32)
		if err != nil {
			color.Red("❌ Invalid lesson ID!")
			return
		}
		
		// Check if lesson exists and is unlocked
		var lesson models.Lesson
		if err := database.DB.First(&lesson, uint(lessonID)).Error; err != nil {
			color.Red("❌ Lesson not found!")
			return
		}
		
		if !lesson.IsCompleted && !isLessonUnlocked(lesson) {
			color.Red("🔒 This lesson is locked! Complete previous lessons first.")
			return
		}
		
		err = exercises.StartLesson(currentUser.ID, uint(lessonID))
		if err != nil {
			color.Red("❌ Error starting lesson: %v", err)
		}
	},
}

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Show user profile and progress",
	Long:  `Display detailed information about your learning progress`,
	Run: func(cmd *cobra.Command, args []string) {
		ensureUser()
		ui.ShowUserProfile(currentUser.ID)
	},
}

var lessonsCmd = &cobra.Command{
	Use:   "lessons",
	Short: "List all available lessons",
	Long:  `Show all lessons with their completion status and requirements`,
	Run: func(cmd *cobra.Command, args []string) {
		ui.ShowLessons()
	},
}

var vocabCmd = &cobra.Command{
	Use:   "vocab [category]",
	Short: "Show vocabulary words",
	Long:  `Display vocabulary words, optionally filtered by category`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		category := ""
		if len(args) > 0 {
			category = args[0]
		}
		ui.ShowVocabulary(category)
	},
}

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show learning statistics",
	Long:  `Display detailed statistics about your learning progress`,
	Run: func(cmd *cobra.Command, args []string) {
		ensureUser()
		ui.ShowStats(currentUser.ID)
	},
}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset all progress (dangerous!)",
	Long:  `Reset all user progress and start fresh. This cannot be undone!`,
	Run: func(cmd *cobra.Command, args []string) {
		color.Red("⚠️  WARNING: This will delete ALL your progress!")
		color.Yellow("This action cannot be undone. Are you sure? (type 'yes' to confirm)")
		
		var response string
		_, err := fmt.Scan(&response)
		if err != nil || response != "yes" {
			color.Green("✅ Reset cancelled.")
			return
		}
		
		// Reset database
		database.DB.Exec("DELETE FROM progresses")
		database.DB.Exec("DELETE FROM users")
		database.DB.Exec("UPDATE lessons SET is_completed = false")
		
		color.Green("✅ All progress has been reset!")
	},
}