package pkg

import (
	"fmt"
	"time"

	"github.com/niiharamegumu/ChronoWork/db"
	"github.com/niiharamegumu/ChronoWork/models"
)

func FormatTime(seconds int) string {
	if seconds >= 3600 {
		return fmt.Sprintf("%02d:%02d:%02d", seconds/3600, (seconds%3600)/60, seconds%60)
	} else if seconds >= 60 {
		return fmt.Sprintf("00:%02d:%02d", seconds/60, seconds%60)
	} else {
		return fmt.Sprintf("00:00:%02d", seconds)
	}
}

func TodayEndTime() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local)
}

func RelativeStartTime() time.Time {
	now := time.Now()
	startTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	var setting models.Setting
	if err := setting.GetSetting(db.DB); err != nil {
		return startTime
	}

	return startTime.AddDate(0, 0, -int(setting.RelativeDate))
}
