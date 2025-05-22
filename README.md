# DuoCLI - German Language Learning in Your Terminal

🇩🇪 **Learn German directly from your command line!**

DuoCLI is a complete German language learning application built for the terminal. It features interactive lessons, vocabulary practice, progress tracking, and gamification elements - all without leaving your CLI.

## ✨ Features

- **📚 Interactive Lessons**: Complete structured German lessons with various exercise types
- **🎯 Multiple Exercise Types**: Translation, multiple choice, fill-in-the-blank
- **📈 Progress Tracking**: XP system, levels, streaks, and detailed statistics
- **📖 Vocabulary Database**: Comprehensive German vocabulary with categories
- **🏆 Gamification**: Level up, earn XP, maintain streaks
- **💾 Persistent Progress**: SQLite database saves your learning journey
- **🎨 Beautiful CLI Interface**: Colorful, intuitive terminal UI

## 🚀 Quick Start

### Installation

```bash
# Clone the repository
git clone git@github.com:marcusziade/duocli.git
cd duocli

# Install dependencies
go mod tidy

# Build the application
go build -o duocli

# Run DuoCLI
./duocli
```

### First Time Setup

When you first run DuoCLI, you'll be prompted to create a profile:

```bash
./duocli
```

Follow the interactive prompts to set up your learning profile.

## 🎮 Usage

### Interactive Mode

Simply run `./duocli` to enter interactive mode with a full menu system:

```bash
./duocli
```

### Command Line Interface

DuoCLI also supports direct commands:

```bash
# Start a specific lesson
./duocli start 1

# View your profile
./duocli profile

# List all lessons
./duocli lessons

# View vocabulary by category
./duocli vocab greetings

# Show learning statistics
./duocli stats

# Reset all progress (careful!)
./duocli reset
```

## 📚 Lesson Structure

### Available Lessons

1. **Basic Greetings** (Level 1)
   - Learn essential German greetings
   - XP Reward: 15

2. **Pronouns** (Level 1)
   - Master German personal pronouns
   - XP Reward: 15

3. **Articles and Nouns** (Level 1)
   - Learn der, die, das and basic nouns
   - XP Reward: 20

4. **Basic Verbs** (Level 2)
   - Essential German verbs and conjugation
   - XP Reward: 25

5. **Family Members** (Level 2)
   - Learn words for family relationships
   - XP Reward: 20

### Exercise Types

- **Translation**: Translate between English and German
- **Multiple Choice**: Choose the correct answer from options
- **Fill in the Blank**: Complete sentences with missing words

## 📖 Vocabulary Categories

- **Greetings**: Hallo, Tschüss, Danke, Bitte, etc.
- **Pronouns**: ich, du, er, sie, wir, etc.
- **Animals**: der Hund, die Katze, etc.
- **Objects**: das Haus, das Auto, etc.
- **Transport**: Vehicle-related vocabulary
- **Verbs**: essen, trinken, gehen, kommen, etc.

## 🏆 Gamification System

### XP and Levels
- Earn XP by completing exercises correctly
- Bonus XP for high lesson completion rates
- Level up every 100 XP

### Streaks
- Maintain daily learning streaks
- Visual streak counter in your profile

### Progress Tracking
- Completion percentage for lessons
- Accuracy statistics
- Recent activity tracking

## 🗄️ Data Structure

DuoCLI uses SQLite for persistent storage:

- **Users**: Profile, XP, level, streak
- **Lessons**: Structured learning content
- **Exercises**: Individual practice items
- **Progress**: Detailed completion tracking
- **Vocabulary**: German-English word pairs

## 🛠️ Development

### Project Structure

```
duocli/
├── main.go                 # Application entry point
├── cmd/                    # CLI commands and interactive mode
│   ├── root.go
│   └── commands.go
├── internal/
│   ├── database/           # Database setup and seeding
│   ├── models/             # Data models
│   ├── exercises/          # Exercise logic and session management
│   └── ui/                 # User interface components
└── data/                   # Data files (if needed)
```

### Key Components

- **Models**: GORM-based data models for all entities
- **Database**: SQLite with auto-migration and seeding
- **Exercises**: Interactive lesson engine with multiple question types
- **UI**: Rich terminal interface with colors and formatting
- **Commands**: Cobra-based CLI with both interactive and direct modes

### Dependencies

- **GORM**: ORM for database operations
- **Cobra**: CLI framework
- **Fatih/Color**: Terminal colors
- **SQLite**: Embedded database

## 🎯 Lesson Completion

- **70%+ accuracy**: Lesson marked as completed
- **80%+ accuracy**: Bonus XP awarded
- **90%+ accuracy**: Perfect score recognition

## 🔄 Progress System

### Unlocking Lessons
- Lesson 1 is always available
- Subsequent lessons unlock when previous lesson is completed
- Visual indicators show lesson status (🔒 🔓 ✅)

### Statistics Tracking
- Total exercises completed
- Accuracy percentage
- Recent activity (7-day window)
- Visual progress bars

## 🎨 UI Features

- **Colorful Interface**: Different colors for different types of information
- **Progress Bars**: Visual representation of completion and accuracy
- **Status Icons**: Emojis for lesson status, exercise types, and achievements
- **Formatted Output**: Clean, structured display of information

## 🚀 Production Ready Features

- **Error Handling**: Comprehensive error handling throughout
- **Data Validation**: Input validation and sanitization
- **Database Migrations**: Automatic schema setup
- **Graceful Degradation**: Handles missing data gracefully
- **Cross-Platform**: Works on Linux, macOS, and Windows

## 🔮 Future Enhancements

Potential areas for expansion:

- **Audio Support**: Pronunciation practice
- **Spaced Repetition**: Smart review scheduling
- **More Languages**: Extend beyond German
- **Online Sync**: Cloud progress backup
- **Community Features**: Shared vocabulary sets
- **Advanced Grammar**: Complex grammar lessons
- **Speaking Exercises**: Voice recognition integration

## 📄 License

This project is open source and available under the MIT License.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit pull requests or open issues for bugs and feature requests.

---

**Happy Learning! Viel Erfolg beim Deutschlernen! 🇩🇪**
