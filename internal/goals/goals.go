package goals

import (
	"fmt"
	"math"
	"time"

	"github.com/ethan-k/pomodoro-cli/internal/config"
	"github.com/ethan-k/pomodoro-cli/internal/db"
)

// GoalType represents the type of goal
type GoalType string

const (
	GoalTypeDaily   GoalType = "daily"
	GoalTypeWeekly  GoalType = "weekly"
	GoalTypeMonthly GoalType = "monthly"
)

// GoalProgress represents progress towards a goal
type GoalProgress struct {
	Type             GoalType  `json:"type"`
	Target           int       `json:"target"`
	Current          int       `json:"current"`
	Percentage       float64   `json:"percentage"`
	Remaining        int       `json:"remaining"`
	IsComplete       bool      `json:"is_complete"`
	IsOverAchieved   bool      `json:"is_over_achieved"`
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
	DaysRemaining    int       `json:"days_remaining"`
	AveragePerDay    float64   `json:"average_per_day"`
	RequiredPerDay   float64   `json:"required_per_day"`
}

// StreakInfo represents streak tracking information
type StreakInfo struct {
	Current    int       `json:"current"`
	Best       int       `json:"best"`
	LastActive time.Time `json:"last_active"`
	IsActive   bool      `json:"is_active"`
}

// GoalManager handles goal tracking functionality
type GoalManager struct {
	db     db.DB
	config *config.Config
}

// NewGoalManager creates a new goal manager
func NewGoalManager(database db.DB, conf *config.Config) *GoalManager {
	return &GoalManager{
		db:     database,
		config: conf,
	}
}

// GetDailyGoalProgress returns daily goal progress
func (gm *GoalManager) GetDailyGoalProgress() (*GoalProgress, error) {
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)
	
	sessions, err := gm.db.GetSessionsByDateRange(today, tomorrow)
	if err != nil {
		return nil, fmt.Errorf("error getting today's sessions: %w", err)
	}

	completed := 0
	for _, session := range sessions {
		if !session.WasBreak {
			completed++
		}
	}

	target := gm.config.Goals.DailyCount
	percentage := float64(completed) / float64(target) * 100
	if percentage > 100 {
		percentage = 100
	}

	return &GoalProgress{
		Type:           GoalTypeDaily,
		Target:         target,
		Current:        completed,
		Percentage:     percentage,
		Remaining:      int(math.Max(0, float64(target-completed))),
		IsComplete:     completed >= target,
		IsOverAchieved: completed > target,
		StartDate:      today,
		EndDate:        tomorrow.Add(-time.Nanosecond),
		DaysRemaining:  1,
		AveragePerDay:  float64(completed),
		RequiredPerDay: math.Max(0, float64(target-completed)),
	}, nil
}

// GetWeeklyGoalProgress returns weekly goal progress
func (gm *GoalManager) GetWeeklyGoalProgress() (*GoalProgress, error) {
	now := time.Now()
	
	// Start from the beginning of the week (Monday)
	daysToMonday := int(now.Weekday())
	if daysToMonday == 0 { // Sunday
		daysToMonday = 6
	} else {
		daysToMonday--
	}
	weekStart := time.Date(now.Year(), now.Month(), now.Day()-daysToMonday, 0, 0, 0, 0, now.Location())
	weekEnd := weekStart.Add(7 * 24 * time.Hour)
	
	sessions, err := gm.db.GetSessionsByDateRange(weekStart, now)
	if err != nil {
		return nil, fmt.Errorf("error getting week's sessions: %w", err)
	}

	completed := 0
	for _, session := range sessions {
		if !session.WasBreak {
			completed++
		}
	}

	target := gm.config.Goals.WeeklyCount
	percentage := float64(completed) / float64(target) * 100
	if percentage > 100 {
		percentage = 100
	}

	daysRemaining := int(math.Ceil(weekEnd.Sub(now).Hours() / 24))
	if daysRemaining < 0 {
		daysRemaining = 0
	}

	averagePerDay := float64(completed) / math.Max(1, float64(7-daysRemaining))
	requiredPerDay := 0.0
	if daysRemaining > 0 && completed < target {
		requiredPerDay = float64(target-completed) / float64(daysRemaining)
	}

	return &GoalProgress{
		Type:           GoalTypeWeekly,
		Target:         target,
		Current:        completed,
		Percentage:     percentage,
		Remaining:      int(math.Max(0, float64(target-completed))),
		IsComplete:     completed >= target,
		IsOverAchieved: completed > target,
		StartDate:      weekStart,
		EndDate:        weekEnd.Add(-time.Nanosecond),
		DaysRemaining:  daysRemaining,
		AveragePerDay:  averagePerDay,
		RequiredPerDay: requiredPerDay,
	}, nil
}

// GetMonthlyGoalProgress returns monthly goal progress
func (gm *GoalManager) GetMonthlyGoalProgress() (*GoalProgress, error) {
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	monthEnd := monthStart.AddDate(0, 1, 0)
	
	sessions, err := gm.db.GetSessionsByDateRange(monthStart, now)
	if err != nil {
		return nil, fmt.Errorf("error getting month's sessions: %w", err)
	}

	completed := 0
	for _, session := range sessions {
		if !session.WasBreak {
			completed++
		}
	}

	// Monthly target is weekly target * 4
	target := gm.config.Goals.WeeklyCount * 4
	percentage := float64(completed) / float64(target) * 100
	if percentage > 100 {
		percentage = 100
	}

	daysRemaining := int(math.Ceil(monthEnd.Sub(now).Hours() / 24))
	if daysRemaining < 0 {
		daysRemaining = 0
	}

	daysInMonth := int(monthEnd.Sub(monthStart).Hours() / 24)
	daysPassed := daysInMonth - daysRemaining
	averagePerDay := float64(completed) / math.Max(1, float64(daysPassed))
	
	requiredPerDay := 0.0
	if daysRemaining > 0 && completed < target {
		requiredPerDay = float64(target-completed) / float64(daysRemaining)
	}

	return &GoalProgress{
		Type:           GoalTypeMonthly,
		Target:         target,
		Current:        completed,
		Percentage:     percentage,
		Remaining:      int(math.Max(0, float64(target-completed))),
		IsComplete:     completed >= target,
		IsOverAchieved: completed > target,
		StartDate:      monthStart,
		EndDate:        monthEnd.Add(-time.Nanosecond),
		DaysRemaining:  daysRemaining,
		AveragePerDay:  averagePerDay,
		RequiredPerDay: requiredPerDay,
	}, nil
}

// GetStreak calculates the current and best streak
func (gm *GoalManager) GetStreak() (*StreakInfo, error) {
	// Get sessions from the last 30 days for streak calculation
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)
	
	sessions, err := gm.db.GetSessionsByDateRange(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error getting sessions for streak: %w", err)
	}

	// Group sessions by date
	dailySessions := make(map[string]int)
	for _, session := range sessions {
		if !session.WasBreak {
			dateKey := session.StartTime.Format("2006-01-02")
			dailySessions[dateKey]++
		}
	}

	// Calculate current streak
	currentStreak := 0
	today := time.Now()
	for i := 0; i < 30; i++ {
		checkDate := today.AddDate(0, 0, -i)
		dateKey := checkDate.Format("2006-01-02")
		
		if count, exists := dailySessions[dateKey]; exists && count > 0 {
			currentStreak++
		} else {
			break
		}
	}

	// Calculate best streak (simplified - would need more historical data for accuracy)
	bestStreak := currentStreak
	tempStreak := 0
	for i := 0; i < 30; i++ {
		checkDate := today.AddDate(0, 0, -i)
		dateKey := checkDate.Format("2006-01-02")
		
		if count, exists := dailySessions[dateKey]; exists && count > 0 {
			tempStreak++
			if tempStreak > bestStreak {
				bestStreak = tempStreak
			}
		} else {
			tempStreak = 0
		}
	}

	lastActive := time.Time{}
	if len(sessions) > 0 {
		for _, session := range sessions {
			if !session.WasBreak && session.StartTime.After(lastActive) {
				lastActive = session.StartTime
			}
		}
	}

	isActive := false
	if !lastActive.IsZero() {
		todayStart := time.Now().Truncate(24 * time.Hour)
		isActive = lastActive.After(todayStart)
	}

	return &StreakInfo{
		Current:    currentStreak,
		Best:       bestStreak,
		LastActive: lastActive,
		IsActive:   isActive,
	}, nil
}

// UpdateGoalTargets updates the goal targets in config
func (gm *GoalManager) UpdateGoalTargets(dailyTarget, weeklyTarget int) error {
	gm.config.Goals.DailyCount = dailyTarget
	gm.config.Goals.WeeklyCount = weeklyTarget
	return config.SaveConfig(gm.config)
}

// GetGoalHistory returns historical goal performance
func (gm *GoalManager) GetGoalHistory(days int) ([]DailyGoalResult, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)
	
	sessions, err := gm.db.GetSessionsByDateRange(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("error getting historical sessions: %w", err)
	}

	// Group sessions by date
	dailySessions := make(map[string][]db.PomodoroSession)
	for _, session := range sessions {
		dateKey := session.StartTime.Format("2006-01-02")
		dailySessions[dateKey] = append(dailySessions[dateKey], session)
	}

	var results []DailyGoalResult
	for i := days - 1; i >= 0; i-- {
		date := endDate.AddDate(0, 0, -i)
		dateKey := date.Format("2006-01-02")
		
		pomodoroCount := 0
		breakCount := 0
		totalDuration := 0
		
		if sessions, exists := dailySessions[dateKey]; exists {
			for _, session := range sessions {
				if session.WasBreak {
					breakCount++
				} else {
					pomodoroCount++
				}
				totalDuration += int(session.DurationSec)
			}
		}

		results = append(results, DailyGoalResult{
			Date:           date,
			PomodoroCount:  pomodoroCount,
			BreakCount:     breakCount,
			TotalDuration:  totalDuration,
			GoalMet:        pomodoroCount >= gm.config.Goals.DailyCount,
			GoalTarget:     gm.config.Goals.DailyCount,
		})
	}

	return results, nil
}

// DailyGoalResult represents a single day's goal performance
type DailyGoalResult struct {
	Date           time.Time `json:"date"`
	PomodoroCount  int       `json:"pomodoro_count"`
	BreakCount     int       `json:"break_count"`
	TotalDuration  int       `json:"total_duration"`
	GoalMet        bool      `json:"goal_met"`
	GoalTarget     int       `json:"goal_target"`
}