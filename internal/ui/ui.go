package ui

import (
	"duocli/internal/database"
	"duocli/internal/models"
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
)

func ShowWelcome() {
	color.Cyan(`
╔══════════════════════════════════════╗
║            DUOCLI - GERMAN           ║
║       Learn German in your CLI      ║
╚══════════════════════════════════════╝
`)
}

func ShowUserProfile(userID uint) {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		color.Red("User not found!")
		return
	}

	// Get completed lessons count
	var completedLessons int64
	database.DB.Model(&models.Lesson{}).Where("is_completed = ?", true).Count(&completedLessons)

	// Get total lessons count
	var totalLessons int64
	database.DB.Model(&models.Lesson{}).Count(&totalLessons)

	// Calculate streak
	streak := calculateStreak(user)

	fmt.Println("\n" + strings.Repeat("=", 50))
	color.Cyan("👤 USER PROFILE")
	fmt.Println(strings.Repeat("=", 50))
	
	color.White("Name: %s", user.Name)
	color.Green("Level: %d", user.Level)
	color.Yellow("XP: %d", user.XP)
	color.Magenta("Streak: %d days 🔥", streak)
	color.Blue("Lessons Completed: %d/%d", completedLessons, totalLessons)
	
	// Progress bar
	progress := float64(completedLessons) / float64(totalLessons) * 100
	progressBar := createProgressBar(int(progress), 30)
	color.White("Progress: %s %.1f%%", progressBar, progress)
	
	fmt.Println(strings.Repeat("=", 50))
}

func ShowLessons() {
	var lessons []models.Lesson
	database.DB.Order("\"order\"").Find(&lessons)

	fmt.Println("\n" + strings.Repeat("=", 50))
	color.Cyan("📚 AVAILABLE LESSONS")
	fmt.Println(strings.Repeat("=", 50))

	for _, lesson := range lessons {
		status := "🔒"
		statusColor := color.RedString
		
		if lesson.IsCompleted {
			status = "✅"
			statusColor = color.GreenString
		} else if isLessonUnlocked(lesson) {
			status = "🔓"
			statusColor = color.YellowString
		}

		fmt.Printf("%s %s Level %d: %s\n", 
			status, 
			statusColor("Lesson %d", lesson.Order), 
			lesson.Level, 
			lesson.Title,
		)
		color.White("   📝 %s", lesson.Description)
		color.Yellow("   💰 XP Reward: %d", lesson.XPReward)
		fmt.Println()
	}
	
	fmt.Println(strings.Repeat("=", 50))
}

func ShowVocabulary(category string) {
	var vocab []models.Vocabulary
	query := database.DB.Order("difficulty, german")
	
	if category != "" {
		query = query.Where("category = ?", category)
	}
	
	query.Find(&vocab)

	fmt.Println("\n" + strings.Repeat("=", 50))
	if category != "" {
		color.Cyan("📖 VOCABULARY - %s", strings.ToUpper(category))
	} else {
		color.Cyan("📖 ALL VOCABULARY")
	}
	fmt.Println(strings.Repeat("=", 50))

	currentCategory := ""
	for _, word := range vocab {
		if word.Category != currentCategory {
			currentCategory = word.Category
			color.Blue("\n🏷️  %s", strings.ToUpper(currentCategory))
			fmt.Println(strings.Repeat("-", 30))
		}
		
		difficulty := strings.Repeat("⭐", word.Difficulty)
		color.White("🇩🇪 %-15s 🇺🇸 %-15s %s", word.German, word.English, difficulty)
		
		if word.Example != "" {
			color.Yellow("   💬 %s", word.Example)
			if word.Translation != "" {
				color.White("   📝 %s", word.Translation)
			}
		}
		fmt.Println()
	}
	
	fmt.Println(strings.Repeat("=", 50))
}

func ShowStats(userID uint) {
	var user models.User
	database.DB.First(&user, userID)

	// Get progress data
	var totalExercises int64
	var correctExercises int64
	
	database.DB.Model(&models.Progress{}).Where("user_id = ?", userID).Count(&totalExercises)
	database.DB.Model(&models.Progress{}).Where("user_id = ? AND is_correct = ?", userID, true).Count(&correctExercises)
	
	accuracy := float64(0)
	if totalExercises > 0 {
		accuracy = float64(correctExercises) / float64(totalExercises) * 100
	}

	// Get recent activity (last 7 days)
	var recentProgress []models.Progress
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	database.DB.Where("user_id = ? AND completed_at > ?", userID, sevenDaysAgo).Find(&recentProgress)

	fmt.Println("\n" + strings.Repeat("=", 50))
	color.Cyan("📊 LEARNING STATISTICS")
	fmt.Println(strings.Repeat("=", 50))
	
	color.White("Total Exercises Completed: %d", totalExercises)
	color.Green("Correct Answers: %d", correctExercises)
	color.Yellow("Accuracy: %.1f%%", accuracy)
	color.Blue("Recent Activity (7 days): %d exercises", len(recentProgress))
	
	// Show accuracy bar
	accuracyBar := createProgressBar(int(accuracy), 30)
	color.White("Accuracy: %s %.1f%%", accuracyBar, accuracy)
	
	fmt.Println(strings.Repeat("=", 50))
}

func createProgressBar(percentage, width int) string {
	if percentage > 100 {
		percentage = 100
	}
	if percentage < 0 {
		percentage = 0
	}
	
	filled := int(float64(percentage) / 100.0 * float64(width))
	empty := width - filled
	
	bar := color.GreenString(strings.Repeat("█", filled)) + 
		   color.WhiteString(strings.Repeat("░", empty))
	
	return fmt.Sprintf("[%s]", bar)
}

func calculateStreak(user models.User) int {
	// Simple streak calculation - this could be more sophisticated
	now := time.Now()
	lastSeen := user.LastSeen
	
	daysDiff := int(now.Sub(lastSeen).Hours() / 24)
	
	if daysDiff <= 1 {
		return user.Streak + 1
	}
	
	return 0
}

func isLessonUnlocked(lesson models.Lesson) bool {
	if lesson.Order == 1 {
		return true // First lesson is always unlocked
	}
	
	// Check if previous lesson is completed
	var prevLesson models.Lesson
	err := database.DB.Where("\"order\" = ?", lesson.Order-1).First(&prevLesson).Error
	if err != nil {
		return false
	}
	
	return prevLesson.IsCompleted
}