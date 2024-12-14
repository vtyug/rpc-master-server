package utils

import (
	"time"
)

// Now 返回当前时间
func Now() time.Time {
	return time.Now()
}

// FormatTime 格式化时间为字符串
func FormatTime(t time.Time, layout string) string {
	return t.Format(layout)
}

// ParseTime 从字符串解析时间
func ParseTime(value, layout string) (time.Time, error) {
	return time.Parse(layout, value)
}

// AddDays 给时间增加指定天数
func AddDays(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}

// StartOfDay 返回一天的开始时间
func StartOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// EndOfDay 返回一天的结束时间
func EndOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, int(time.Second-time.Nanosecond), t.Location())
}

// DaysBetween 计算两个时间之间的天数
func DaysBetween(start, end time.Time) int {
	return int(end.Sub(start).Hours() / 24)
}
