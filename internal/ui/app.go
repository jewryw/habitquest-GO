package ui

import (
    "fmt"
    "habitquest/internal/db"
    "strconv"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
)

func BuildUI(window fyne.Window) fyne.CanvasObject {
    // === ВКЛАДКА: ПРИВЫЧКИ ===
    habitsList := widget.NewList(
        func() int {
            habits, _ := db.GetAllHabits()
            return len(habits)
        },
        func() fyne.CanvasObject {
            return widget.NewLabel("Привычка")
        },
        func(id widget.ListItemID, obj fyne.CanvasObject) {
            habits, _ := db.GetAllHabits()
            h := habits[id]
            obj.(*widget.Label).SetText(fmt.Sprintf("%s | +%d очков | 🔥 %d дней",
                h.Title, h.Points, h.Streak))
        },
    )

    titleEntry := widget.NewEntry()
    titleEntry.SetPlaceHolder("Название (например: Читать 30 мин)")
    descEntry := widget.NewEntry()
    descEntry.SetPlaceHolder("Описание")
    pointsEntry := widget.NewEntry()
    pointsEntry.SetPlaceHolder("Очки (по умолчанию 10)")

    addBtn := widget.NewButton("➕ Добавить привычку", func() {
        title := titleEntry.Text
        if title == "" {
            return
        }
        points, _ := strconv.Atoi(pointsEntry.Text)
        if points == 0 {
            points = 10
        }
        db.CreateHabit(title, descEntry.Text, points)
        titleEntry.SetText("")
        descEntry.SetText("")
        pointsEntry.SetText("")
        habitsList.Refresh()
    })

    completeBtn := widget.NewButton("✅ Выполнено сегодня", func() {
        habits, _ := db.GetAllHabits()
        if len(habits) > 0 {
            db.CompleteHabit(habits[0].ID, habits[0].Points)
            habitsList.Refresh()
        }
    })

    deleteBtn := widget.NewButton("🗑 Удалить последнюю", func() {
        habits, _ := db.GetAllHabits()
        if len(habits) > 0 {
            db.DeleteHabit(habits[len(habits)-1].ID)
            habitsList.Refresh()
        }
    })

    habitsTab := container.NewVBox(
        titleEntry, descEntry, pointsEntry,
        addBtn, completeBtn, deleteBtn,
        habitsList,
    )

    // === ВКЛАДКА: СТАТИСТИКА ===
    refreshStats := func() string {
        points := db.GetTotalPoints()
        level := db.GetLevel()
        monthly := db.GetMonthlyCompletions()
        avgStreak := db.GetAverageStreak()
        return fmt.Sprintf(
            "🏆 Всего очков: %d\n"+
                "⭐ Уровень: %d\n"+
                "📅 Выполнено в этом месяце: %d\n"+
                "🔥 Средний streak: %.1f дней",
            points, level, monthly, avgStreak,
        )
    }

    statsLabel := widget.NewLabel(refreshStats())
    statsLabel.Wrapping = fyne.TextWrapWord

    statsBtn := widget.NewButton("🔄 Обновить", func() {
        statsLabel.SetText(refreshStats())
    })

    statsTab := container.NewVBox(statsLabel, statsBtn)

    // === ВКЛАДКИ ===
    tabs := container.NewAppTabs(
        container.NewTabItem("🎯 Привычки", habitsTab),
        container.NewTabItem("📊 Статистика", statsTab),
    )

    return tabs
}