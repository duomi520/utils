package utils

import (
	"fmt"
	"time"
)

// StartOfMonth  获取月初
func StartOfMonth(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
}

// EndOfMonth 获取月底
func EndOfMonth(date time.Time) time.Time {
	firstDayOfNextMonth := StartOfMonth(date).AddDate(0, 1, 0)
	return firstDayOfNextMonth.Add(-time.Second)
}

// StartOfDayOfWeek 获取每周的开始日
func StartOfDayOfWeek(date time.Time) time.Time {
	daysSinceSunday := int(date.Weekday())
	return date.AddDate(0, 0, -daysSinceSunday+1)
}

// EndOfDayOfWeek 获取每周的结束日
func EndOfDayOfWeek(date time.Time) time.Time {
	daysUntilSaturday := 7 - int(date.Weekday())
	return date.AddDate(0, 0, daysUntilSaturday)
}

// StartOfYear 获取新年伊始
func StartOfYear(date time.Time) time.Time {
	return time.Date(date.Year(), time.January, 1, 0, 0, 0, 0, date.Location())
}

// EndOfYear 获取年底
func EndOfYear(date time.Time) time.Time {
	startOfNextYear := StartOfYear(date).AddDate(1, 0, 0)
	return startOfNextYear.Add(-time.Second)
}

// StartOfQuarter 获取季度初数据
func StartOfQuarter(date time.Time) time.Time {
	// you can directly use 0, 1, 2, 3 quarter
	quarter := (int(date.Month()) - 1) / 3
	startMonth := time.Month(quarter*3 + 1)
	return time.Date(date.Year(), startMonth, 1, 0, 0, 0, 0, date.Location())
}

// EndOfQuarter 获取季度末
func EndOfQuarter(date time.Time) time.Time {
	startOfNextQuarter := StartOfQuarter(date).AddDate(0, 3, 0)
	return startOfNextQuarter.Add(-time.Second)
}

// FormatDuration 将持续时间格式化字符串
func FormatDuration(duration time.Duration) string {
	days := int(duration.Hours() / 24)
	hours := int(duration.Hours()) % 24
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60
	return fmt.Sprintf("%d天 %02d小时 %02d分 %02d秒", days, hours, minutes, seconds)
}
