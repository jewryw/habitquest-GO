package db

import (
    "database/sql"
    "habitquest/internal/models"
    "time"

    _ "modernc.org/sqlite"
)

var DB *sql.DB

func Init() error {
    var err error
    DB, err = sql.Open("sqlite", "habitquest.db")
    if err != nil {
        return err
    }

    // Создаём таблицы
    _, err = DB.Exec(`
        CREATE TABLE IF NOT EXISTS habits (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL,
            description TEXT,
            points INTEGER DEFAULT 10,
            streak INTEGER DEFAULT 0,
            best_streak INTEGER DEFAULT 0,
            created_at DATETIME DEFAULT CURRENT_TIMESTAMP
        );
        CREATE TABLE IF NOT EXISTS completions (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            habit_id INTEGER,
            completed_at DATE,
            points_earned INTEGER,
            FOREIGN KEY(habit_id) REFERENCES habits(id)
        );
        CREATE TABLE IF NOT EXISTS achievements (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT,
            description TEXT,
            unlocked_at DATETIME DEFAULT CURRENT_TIMESTAMP
        );
    `)
    return err
}

// === HABITS ===
func CreateHabit(title, desc string, points int) error {
    _, err := DB.Exec(`INSERT INTO habits (title, description, points) VALUES (?, ?, ?)`,
        title, desc, points)
    return err
}

func GetAllHabits() ([]models.Habit, error) {
    rows, err := DB.Query(`SELECT id, title, description, points, streak, best_streak, created_at FROM habits`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var habits []models.Habit
    for rows.Next() {
        var h models.Habit
        rows.Scan(&h.ID, &h.Title, &h.Description, &h.Points, &h.Streak, &h.BestStreak, &h.CreatedAt)
        habits = append(habits, h)
    }
    return habits, nil
}

func DeleteHabit(id int) error {
    _, err := DB.Exec(`DELETE FROM habits WHERE id = ?`, id)
    return err
}

// === COMPLETIONS ===
func CompleteHabit(habitID, points int) error {
    today := time.Now().Format("2006-01-02")
    yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

    // Проверяем, выполнялась ли привычка сегодня
    var count int
    DB.QueryRow(`SELECT COUNT(*) FROM completions WHERE habit_id = ? AND completed_at = ?`,
        habitID, today).Scan(&count)
    if count > 0 {
        return nil // уже выполнено сегодня
    }

    // Проверяем, выполнялась ли вчера (для streak)
    var yesterdayCount int
    DB.QueryRow(`SELECT COUNT(*) FROM completions WHERE habit_id = ? AND completed_at = ?`,
        habitID, yesterday).Scan(&yesterdayCount)

    newStreak := 1
    if yesterdayCount > 0 {
        // Получаем текущий streak
        var currentStreak int
        DB.QueryRow(`SELECT streak FROM habits WHERE id = ?`, habitID).Scan(&currentStreak)
        newStreak = currentStreak + 1
    }

    // Обновляем streak и best_streak
    DB.Exec(`UPDATE habits SET streak = ?, best_streak = MAX(best_streak, ?) WHERE id = ?`,
        newStreak, newStreak, habitID)

    // Записываем выполнение
    _, err := DB.Exec(`INSERT INTO completions (habit_id, completed_at, points_earned) VALUES (?, ?, ?)`,
        habitID, today, points)
    return err
}

// === СТАТИСТИКА ===
func GetTotalPoints() int {
    var total int
    DB.QueryRow(`SELECT COALESCE(SUM(points_earned), 0) FROM completions`).Scan(&total)
    return total
}

func GetLevel() int {
    return GetTotalPoints()/100 + 1 // каждые 100 очков = новый уровень
}

func GetMonthlyCompletions() int {
    month := time.Now().Format("2006-01")
    var count int
    DB.QueryRow(`SELECT COUNT(*) FROM completions WHERE strftime('%Y-%m', completed_at) = ?`, month).Scan(&count)
    return count
}

func GetAverageStreak() float64 {
    var avg float64
    DB.QueryRow(`SELECT COALESCE(AVG(streak), 0) FROM habits`).Scan(&avg)
    return avg
}