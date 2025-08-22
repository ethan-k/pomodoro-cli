package goals

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ethan-k/pomodoro-cli/internal/audio"
	"github.com/ethan-k/pomodoro-cli/internal/notify"
)

// Achievement represents different types of goal achievements
type Achievement int

const (
	AchievementDailyGoal Achievement = iota
	AchievementWeeklyGoal
	AchievementMonthlyGoal
	AchievementStreakMilestone
	AchievementPersonalBest
	AchievementOverachiever
)

// AchievementInfo contains details about an achievement
type AchievementInfo struct {
	Type        Achievement `json:"type"`
	Title       string      `json:"title"`
	Message     string      `json:"message"`
	Icon        string      `json:"icon"`
	SoundFile   string      `json:"sound_file,omitempty"`
	Celebration string      `json:"celebration"`
}

// NotificationManager handles goal achievement notifications
type NotificationManager struct {
	audioPlayer *audio.Player
	notifier    *notify.Notifier
}

// NewNotificationManager creates a new notification manager
func NewNotificationManager(audioPlayer *audio.Player, notifier *notify.Notifier) *NotificationManager {
	return &NotificationManager{
		audioPlayer: audioPlayer,
		notifier:    notifier,
	}
}

// GetAchievementInfo returns information about an achievement
func (nm *NotificationManager) GetAchievementInfo(achievementType Achievement, details map[string]interface{}) AchievementInfo {
	switch achievementType {
	case AchievementDailyGoal:
		return AchievementInfo{
			Type:        AchievementDailyGoal,
			Title:       "Daily Goal Complete! 🎉",
			Message:     fmt.Sprintf("Congratulations! You've completed your daily goal of %v pomodoros!", details["target"]),
			Icon:        "🎯",
			SoundFile:   "goal_complete.wav",
			Celebration: "🎉 🍅 🎉 🍅 🎉",
		}

	case AchievementWeeklyGoal:
		return AchievementInfo{
			Type:        AchievementWeeklyGoal,
			Title:       "Weekly Goal Achieved! 🏆",
			Message:     fmt.Sprintf("Amazing work! You've reached your weekly target of %v pomodoros!", details["target"]),
			Icon:        "🏆",
			SoundFile:   "weekly_goal_complete.wav",
			Celebration: "🏆 🌟 🏆 🌟 🏆",
		}

	case AchievementMonthlyGoal:
		return AchievementInfo{
			Type:        AchievementMonthlyGoal,
			Title:       "Monthly Goal Crushed! 🚀",
			Message:     fmt.Sprintf("Incredible dedication! You've completed %v pomodoros this month!", details["target"]),
			Icon:        "🚀",
			SoundFile:   "monthly_goal_complete.wav",
			Celebration: "🚀 ⭐ 🚀 ⭐ 🚀",
		}

	case AchievementStreakMilestone:
		streak := details["streak"].(int)
		var title, celebration string
		switch {
		case streak >= 30:
			title = "Legendary Streak! 🔥"
			celebration = "🔥 👑 🔥 👑 🔥"
		case streak >= 14:
			title = "Two Week Streak! 🔥"
			celebration = "🔥 💪 🔥 💪 🔥"
		case streak >= 7:
			title = "Week Streak! 🔥"
			celebration = "🔥 ⚡ 🔥 ⚡ 🔥"
		case streak >= 3:
			title = "Streak Building! 🔥"
			celebration = "🔥 🎯 🔥 🎯 🔥"
		default:
			title = "Streak Started! 🔥"
			celebration = "🔥 🌱 🔥 🌱 🔥"
		}

		return AchievementInfo{
			Type:        AchievementStreakMilestone,
			Title:       title,
			Message:     fmt.Sprintf("You're on fire! %d days in a row of completing your goals!", streak),
			Icon:        "🔥",
			SoundFile:   "streak_milestone.wav",
			Celebration: celebration,
		}

	case AchievementPersonalBest:
		count := details["count"].(int)
		return AchievementInfo{
			Type:        AchievementPersonalBest,
			Title:       "Personal Best! 🌟",
			Message:     fmt.Sprintf("New record! You've completed %d pomodoros in a single day!", count),
			Icon:        "🌟",
			SoundFile:   "personal_best.wav",
			Celebration: "🌟 🎊 🌟 🎊 🌟",
		}

	case AchievementOverachiever:
		extra := details["extra"].(int)
		return AchievementInfo{
			Type:        AchievementOverachiever,
			Title:       "Overachiever! 💎",
			Message:     fmt.Sprintf("Going above and beyond! You've exceeded your goal by %d extra pomodoros!", extra),
			Icon:        "💎",
			SoundFile:   "overachiever.wav",
			Celebration: "💎 ✨ 💎 ✨ 💎",
		}

	default:
		return AchievementInfo{
			Type:        achievementType,
			Title:       "Goal Achievement! 🎉",
			Message:     "Congratulations on your achievement!",
			Icon:        "🎉",
			Celebration: "🎉 🎉 🎉 🎉 🎉",
		}
	}
}

// SendAchievementNotification sends a notification for a goal achievement
func (nm *NotificationManager) SendAchievementNotification(achievementType Achievement, details map[string]interface{}) error {
	info := nm.GetAchievementInfo(achievementType, details)

	// Send desktop notification
	if nm.notifier != nil {
		if err := nm.notifier.Send(info.Title, info.Message, info.Icon); err != nil {
			// Log error but don't fail - notification is not critical
			fmt.Printf("Warning: Failed to send desktop notification: %v\n", err)
		}
	}

	// Play achievement sound
	if nm.audioPlayer != nil && info.SoundFile != "" {
		soundPath := filepath.Join("sounds", info.SoundFile)
		if err := nm.audioPlayer.PlaySound(soundPath); err != nil {
			// Log error but don't fail - sound is not critical
			fmt.Printf("Warning: Failed to play achievement sound: %v\n", err)
		}
	}

	// Display console celebration
	fmt.Println("\n" + info.Celebration)
	fmt.Printf("🎊 %s 🎊\n", info.Title)
	fmt.Println(info.Message)
	fmt.Println(info.Celebration + "\n")

	return nil
}

// CheckForAchievements checks if any achievements should be triggered
func (gm *GoalManager) CheckForAchievements(notificationManager *NotificationManager) error {
	// Check daily goal achievement
	daily, err := gm.GetDailyGoalProgress()
	if err != nil {
		return fmt.Errorf("error checking daily progress: %w", err)
	}

	if daily.IsComplete {
		details := map[string]interface{}{
			"target": daily.Target,
			"current": daily.Current,
		}

		if daily.IsOverAchieved {
			// Send overachiever notification
			extraDetails := map[string]interface{}{
				"extra": daily.Current - daily.Target,
			}
			notificationManager.SendAchievementNotification(AchievementOverachiever, extraDetails)
		} else {
			// Send daily goal completion notification
			notificationManager.SendAchievementNotification(AchievementDailyGoal, details)
		}
	}

	// Check weekly goal achievement
	weekly, err := gm.GetWeeklyGoalProgress()
	if err != nil {
		return fmt.Errorf("error checking weekly progress: %w", err)
	}

	if weekly.IsComplete {
		details := map[string]interface{}{
			"target": weekly.Target,
			"current": weekly.Current,
		}
		notificationManager.SendAchievementNotification(AchievementWeeklyGoal, details)
	}

	// Check monthly goal achievement
	monthly, err := gm.GetMonthlyGoalProgress()
	if err != nil {
		return fmt.Errorf("error checking monthly progress: %w", err)
	}

	if monthly.IsComplete {
		details := map[string]interface{}{
			"target": monthly.Target,
			"current": monthly.Current,
		}
		notificationManager.SendAchievementNotification(AchievementMonthlyGoal, details)
	}

	// Check streak milestones
	streak, err := gm.GetStreak()
	if err != nil {
		return fmt.Errorf("error checking streak: %w", err)
	}

	if streak.Current > 0 && shouldCelebrateStreak(streak.Current) {
		details := map[string]interface{}{
			"streak": streak.Current,
		}
		notificationManager.SendAchievementNotification(AchievementStreakMilestone, details)
	}

	// Check for personal best (simplified - would need historical tracking)
	if daily.Current > 0 && isPersonalBest(daily.Current) {
		details := map[string]interface{}{
			"count": daily.Current,
		}
		notificationManager.SendAchievementNotification(AchievementPersonalBest, details)
	}

	return nil
}

// shouldCelebrateStreak determines if a streak milestone should be celebrated
func shouldCelebrateStreak(streak int) bool {
	// Celebrate at specific milestones
	milestones := []int{3, 7, 14, 21, 30, 60, 90, 180, 365}
	for _, milestone := range milestones {
		if streak == milestone {
			return true
		}
	}
	return false
}

// isPersonalBest checks if the current count is a personal best (simplified)
func isPersonalBest(count int) bool {
	// This is a simplified check - in a real implementation,
	// you'd want to track historical daily maximums
	return count >= 15 // Arbitrary threshold for demonstration
}

// DisplayGoalCelebration shows a celebration message in the terminal
func DisplayGoalCelebration(achievementType Achievement, details map[string]interface{}) {
	info := getBasicAchievementInfo(achievementType, details)
	
	fmt.Println("\n" + strings.Repeat("🎉", 20))
	fmt.Printf("🏆 %s 🏆\n", info.Title)
	fmt.Println(info.Message)
	fmt.Println(info.Celebration)
	fmt.Println(strings.Repeat("🎉", 20) + "\n")
}

func getBasicAchievementInfo(achievementType Achievement, details map[string]interface{}) AchievementInfo {
	// This is a simplified version that doesn't require the full NotificationManager
	switch achievementType {
	case AchievementDailyGoal:
		return AchievementInfo{
			Title:       "Daily Goal Complete!",
			Message:     "🎯 Congratulations! You've reached your daily target!",
			Celebration: "🎉 🍅 🎉 🍅 🎉",
		}
	case AchievementWeeklyGoal:
		return AchievementInfo{
			Title:       "Weekly Goal Achieved!",
			Message:     "🏆 Amazing work! You've completed your weekly goal!",
			Celebration: "🏆 🌟 🏆 🌟 🏆",
		}
	default:
		return AchievementInfo{
			Title:       "Achievement Unlocked!",
			Message:     "🎉 Great job on reaching your goal!",
			Celebration: "🎉 🎉 🎉 🎉 🎉",
		}
	}
}