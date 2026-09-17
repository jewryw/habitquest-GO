package main

import (
    "habitquest/internal/db"
    "habitquest/internal/ui"
    "log"
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
)

func main() {
    // Инициализация БД
    if err := db.Init(); err != nil {
        log.Fatal("Ошибка инициализации БД:", err)
    }

    // Создание окна
    a := app.New()
    w := a.NewWindow("HabitQuest — Трекер привычек")
    w.Resize(fyne.NewSize(600, 500))

    w.SetContent(ui.BuildUI(w))
    w.ShowAndRun()
}