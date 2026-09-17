package models

import "time"

type Habit struct {
    ID          int
    Title       string
    Description string
    Points      int       // очки за выполнение
    Streak      int       // текущая серия дней
    BestStreak  int       // лучшая серия
    CreatedAt   time.Time
}

type Completion struct {
    ID          int
    HabitID     int
    CompletedAt time.Time
    PointsEarned int
}

type Achievement struct {
    ID          int
    Name        string
    Description string
    UnlockedAt  time.Time
}