package cmd

import (
	"bufio"
	"duocli/internal/database"
	"duocli/internal/exercises"
	"duocli/internal/models"
	"duocli/internal/ui"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var currentUser *models.User

var rootCmd = &cobra.Command{
	Use:   "duocli",
	Short: "Learn German in your CLI",
	Long: `DuoCLI is a complete German language learning application for the command line.
Learn vocabulary, complete lessons, and track your progress - all from your terminal!`,
	Run: func(cmd *cobra.Command, args []string) {
		runInteractiveMode()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(profileCmd)
	rootCmd.AddCommand(lessonsCmd)
	rootCmd.AddCommand(vocabCmd)
	rootCmd.AddCommand(statsCmd)
	rootCmd.AddCommand(resetCmd)
}

func runInteractiveMode() {
	ui.ShowWelcome()
	
	// Ensure user exists
	ensureUser()
	
	scanner := bufio.NewScanner(os.Stdin)
	
	for {
		showMainMenu()
		fmt.Print("Choose an option: ")
		
		if !scanner.Scan() {
			break
		}
		
		choice := strings.TrimSpace(scanner.Text())
		
		switch choice {
		case "1":
			startLearning()
		case "2":
			ui.ShowUserProfile(currentUser.ID)
		case "3":
			ui.ShowLessons()
		case "4":
			showVocabMenu()
		case "5":
			ui.ShowStats(currentUser.ID)
		case "6":
			color.Cyan("👋 Auf Wiedersehen! (Goodbye!)")
			return
		default:
			color.Red("❌ Invalid option. Please try again.")
		}
		
		fmt.Print("\nPress Enter to continue...")
		scanner.Scan()
	}
}

func showMainMenu() {
	fmt.Println("\n" + strings.Repeat("=", 50))
	color.Cyan("🏠 MAIN MENU")
	fmt.Println(strings.Repeat("=", 50))
	
	color.White("1. 🎓 Start Learning")
	color.White("2. 👤 Profile")
	color.White("3. 📚 View Lessons")
	color.White("4. 📖 Vocabulary")
	color.White("5. 📊 Statistics")
	color.White("6. 🚪 Exit")
	
	fmt.Println(strings.Repeat("=", 50))
}

func startLearning() {
	var lessons []models.Lesson
	database.DB.Order("\"order\"").Find(&lessons)
	
	fmt.Println("\n" + strings.Repeat("=", 40))
	color.Cyan("🎓 SELECT A LESSON")
	fmt.Println(strings.Repeat("=", 40))
	
	availableLessons := []models.Lesson{}
	
	for _, lesson := range lessons {
		status := "🔒 Locked"
		available := false
		
		if lesson.IsCompleted {
			status = "✅ Completed"
			available = true
		} else if isLessonUnlocked(lesson) {
			status = "🔓 Available"
			available = true
		}
		
		if available {
			availableLessons = append(availableLessons, lesson)
		}
		
		color.White("%d. %s - %s", lesson.Order, lesson.Title, status)
		color.Yellow("   %s (XP: %d)", lesson.Description, lesson.XPReward)
	}
	
	if len(availableLessons) == 0 {
		color.Red("❌ No lessons available!")
		return
	}
	
	fmt.Print("\nSelect lesson number (0 to cancel): ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	choice := strings.TrimSpace(scanner.Text())
	
	if choice == "0" {
		return
	}
	
	lessonNum, err := strconv.Atoi(choice)
	if err != nil || lessonNum < 1 || lessonNum > len(lessons) {
		color.Red("❌ Invalid lesson number!")
		return
	}
	
	selectedLesson := lessons[lessonNum-1]
	
	// Check if lesson is unlocked
	if !selectedLesson.IsCompleted && !isLessonUnlocked(selectedLesson) {
		color.Red("🔒 This lesson is locked! Complete previous lessons first.")
		return
	}
	
	// Start the lesson
	err = exercises.StartLesson(currentUser.ID, selectedLesson.ID)
	if err != nil {
		color.Red("❌ Error starting lesson: %v", err)
	}
}

func showVocabMenu() {
	fmt.Println("\n" + strings.Repeat("=", 40))
	color.Cyan("📖 VOCABULARY MENU")
	fmt.Println(strings.Repeat("=", 40))
	
	color.White("1. 👋 Greetings")
	color.White("2. 👁️  Pronouns")
	color.White("3. 🐕 Animals")
	color.White("4. 🏠 Objects")
	color.White("5. 🚗 Transport")
	color.White("6. 🔧 Verbs")
	color.White("7. 📚 All Vocabulary")
	color.White("0. 🔙 Back to Main Menu")
	
	fmt.Print("\nChoose category: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	choice := strings.TrimSpace(scanner.Text())
	
	categories := map[string]string{
		"1": "greetings",
		"2": "pronouns",
		"3": "animals",
		"4": "objects",
		"5": "transport",
		"6": "verbs",
		"7": "",
	}
	
	if choice == "0" {
		return
	}
	
	if category, exists := categories[choice]; exists {
		ui.ShowVocabulary(category)
	} else {
		color.Red("❌ Invalid option!")
	}
}

func ensureUser() {
	var user models.User
	err := database.DB.First(&user).Error
	
	if err != nil {
		// Create new user
		color.Yellow("👋 Welcome to DuoCLI! Let's set up your profile.")
		fmt.Print("Enter your name: ")
		
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		name := strings.TrimSpace(scanner.Text())
		
		if name == "" {
			name = "Learner"
		}
		
		user = models.User{
			Name:     name,
			Level:    1,
			XP:       0,
			Streak:   0,
		}
		
		database.DB.Create(&user)
		color.Green("✅ Profile created! Welcome, %s!", user.Name)
	}
	
	currentUser = &user
}

func isLessonUnlocked(lesson models.Lesson) bool {
	if lesson.Order == 1 {
		return true
	}
	
	var prevLesson models.Lesson
	err := database.DB.Where("\"order\" = ?", lesson.Order-1).First(&prevLesson).Error
	if err != nil {
		return false
	}
	
	return prevLesson.IsCompleted
}