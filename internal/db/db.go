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

