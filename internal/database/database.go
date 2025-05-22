package database

import (
	"duocli/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() error {
	var err error
	DB, err = gorm.Open(sqlite.Open("duocli.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return err
	}

	// Auto migrate the schema
	err = DB.AutoMigrate(
		&models.User{},
		&models.Lesson{},
		&models.Exercise{},
		&models.Progress{},
		&models.Vocabulary{},
	)
	if err != nil {
		return err
	}

	// Seed initial data
	return seedData()
}

func seedData() error {
	// Check if data already exists
	var lessonCount int64
	DB.Model(&models.Lesson{}).Count(&lessonCount)
	if lessonCount > 0 {
		return nil // Data already seeded
	}

	// Seed vocabulary
	vocabularies := []models.Vocabulary{
		{German: "Hallo", English: "Hello", Category: "greetings", Difficulty: 1, Example: "Hallo, wie geht es dir?", Translation: "Hello, how are you?"},
		{German: "Tschüss", English: "Goodbye", Category: "greetings", Difficulty: 1, Example: "Tschüss, bis später!", Translation: "Goodbye, see you later!"},
		{German: "Danke", English: "Thank you", Category: "greetings", Difficulty: 1, Example: "Danke für deine Hilfe.", Translation: "Thank you for your help."},
		{German: "Bitte", English: "Please/You're welcome", Category: "greetings", Difficulty: 1, Example: "Bitte schön!", Translation: "You're welcome!"},
		{German: "Entschuldigung", English: "Excuse me/Sorry", Category: "greetings", Difficulty: 1, Example: "Entschuldigung, wo ist der Bahnhof?", Translation: "Excuse me, where is the train station?"},
		{German: "Ja", English: "Yes", Category: "basic", Difficulty: 1, Example: "Ja, das ist richtig.", Translation: "Yes, that is correct."},
		{German: "Nein", English: "No", Category: "basic", Difficulty: 1, Example: "Nein, das ist falsch.", Translation: "No, that is wrong."},
		{German: "ich", English: "I", Category: "pronouns", Difficulty: 1, Example: "Ich bin Student.", Translation: "I am a student."},
		{German: "du", English: "you (informal)", Category: "pronouns", Difficulty: 1, Example: "Du bist nett.", Translation: "You are nice."},
		{German: "er", English: "he", Category: "pronouns", Difficulty: 1, Example: "Er kommt aus Deutschland.", Translation: "He comes from Germany."},
		{German: "sie", English: "she/they", Category: "pronouns", Difficulty: 1, Example: "Sie ist Lehrerin.", Translation: "She is a teacher."},
		{German: "wir", English: "we", Category: "pronouns", Difficulty: 1, Example: "Wir lernen Deutsch.", Translation: "We are learning German."},
		{German: "der Hund", English: "the dog", Category: "animals", Difficulty: 1, Example: "Der Hund ist süß.", Translation: "The dog is cute."},
		{German: "die Katze", English: "the cat", Category: "animals", Difficulty: 1, Example: "Die Katze schläft.", Translation: "The cat is sleeping."},
		{German: "das Haus", English: "the house", Category: "objects", Difficulty: 1, Example: "Das Haus ist groß.", Translation: "The house is big."},
		{German: "das Auto", English: "the car", Category: "transport", Difficulty: 1, Example: "Das Auto ist rot.", Translation: "The car is red."},
		{German: "essen", English: "to eat", Category: "verbs", Difficulty: 2, Example: "Ich esse einen Apfel.", Translation: "I eat an apple."},
		{German: "trinken", English: "to drink", Category: "verbs", Difficulty: 2, Example: "Wir trinken Wasser.", Translation: "We drink water."},
		{German: "gehen", English: "to go", Category: "verbs", Difficulty: 2, Example: "Sie geht nach Hause.", Translation: "She goes home."},
		{German: "kommen", English: "to come", Category: "verbs", Difficulty: 2, Example: "Er kommt morgen.", Translation: "He comes tomorrow."},
	}

	for _, vocab := range vocabularies {
		DB.Create(&vocab)
	}

	// Seed lessons
	lessons := []models.Lesson{
		{Title: "Basic Greetings", Description: "Learn essential German greetings", Level: 1, Order: 1, XPReward: 15},
		{Title: "Pronouns", Description: "Master German personal pronouns", Level: 1, Order: 2, XPReward: 15},
		{Title: "Articles and Nouns", Description: "Learn der, die, das and basic nouns", Level: 1, Order: 3, XPReward: 20},
		{Title: "Basic Verbs", Description: "Essential German verbs and conjugation", Level: 2, Order: 4, XPReward: 25},
		{Title: "Family Members", Description: "Learn words for family relationships", Level: 2, Order: 5, XPReward: 20},
	}

	for i, lesson := range lessons {
		DB.Create(&lesson)
		lessons[i] = lesson // Update with ID
	}

	// Seed exercises for each lesson
	exercises := map[string][]models.Exercise{
		"Basic Greetings": {
			{Type: "translation", Question: "How do you say 'Hello' in German?", Answer: "Hallo", Hint: "It sounds similar to English", Order: 1, Difficulty: 1},
			{Type: "multiple_choice", Question: "What does 'Danke' mean?", Answer: "Thank you", Options: `["Thank you", "Goodbye", "Please", "Hello"]`, Order: 2, Difficulty: 1},
			{Type: "translation", Question: "Translate: Goodbye", Answer: "Tschüss", Hint: "Pronounced 'choos'", Order: 3, Difficulty: 1},
			{Type: "fill_blank", Question: "_____, wie geht es dir?", Answer: "Hallo", Hint: "Greeting", Order: 4, Difficulty: 1},
		},
		"Pronouns": {
			{Type: "translation", Question: "How do you say 'I' in German?", Answer: "ich", Hint: "Lowercase in German", Order: 1, Difficulty: 1},
			{Type: "multiple_choice", Question: "What does 'du' mean?", Answer: "you (informal)", Options: `["I", "you (informal)", "he", "we"]`, Order: 2, Difficulty: 1},
			{Type: "translation", Question: "Translate: we", Answer: "wir", Hint: "Sounds like 'veer'", Order: 3, Difficulty: 1},
			{Type: "fill_blank", Question: "_____ bist nett.", Answer: "Du", Hint: "You (informal)", Order: 4, Difficulty: 1},
		},
		"Articles and Nouns": {
			{Type: "multiple_choice", Question: "What is the article for 'Hund' (dog)?", Answer: "der", Options: `["der", "die", "das", "den"]`, Order: 1, Difficulty: 1},
			{Type: "translation", Question: "Translate: the cat", Answer: "die Katze", Hint: "Feminine article", Order: 2, Difficulty: 1},
			{Type: "multiple_choice", Question: "What is the article for 'Haus' (house)?", Answer: "das", Options: `["der", "die", "das", "den"]`, Order: 3, Difficulty: 1},
			{Type: "fill_blank", Question: "_____ Auto ist rot.", Answer: "Das", Hint: "Neuter article", Order: 4, Difficulty: 1},
		},
		"Basic Verbs": {
			{Type: "translation", Question: "How do you say 'to eat' in German?", Answer: "essen", Hint: "Similar to English", Order: 1, Difficulty: 2},
			{Type: "multiple_choice", Question: "What does 'trinken' mean?", Answer: "to drink", Options: `["to eat", "to drink", "to go", "to come"]`, Order: 2, Difficulty: 2},
			{Type: "translation", Question: "Translate: to go", Answer: "gehen", Hint: "Think of 'go' sounds", Order: 3, Difficulty: 2},
			{Type: "fill_blank", Question: "Ich _____ einen Apfel.", Answer: "esse", Hint: "I eat", Order: 4, Difficulty: 2},
		},
	}

	for _, lesson := range lessons {
		if exerciseList, exists := exercises[lesson.Title]; exists {
			for _, exercise := range exerciseList {
				exercise.LessonID = lesson.ID
				DB.Create(&exercise)
			}
		}
	}

	return nil
}