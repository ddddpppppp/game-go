package utils

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
)

// tokens计算
type TimeHelperUtil struct {
}

var TimeHelper = &TimeHelperUtil{}

func (u *TimeHelperUtil) ParseAnyTimeToGTime(input string) (*gtime.Time, error) {
	// 尝试解析完整日期格式 "2025年4月6日 03:34"
	if t, err := u.parseFullDateFormat(input); err == nil {
		return t, nil
	}

	// 尝试解析带星期的格式 "周六10:01"
	if t, err := u.parseWeekdayTimeFormat(input); err == nil {
		return t, nil
	}

	// 尝试解析纯时间格式 "10:01"
	if t, err := u.parseSimpleTimeFormat(input); err == nil {
		return t, nil
	}

	return nil, fmt.Errorf("unrecognized time format: %s", input)
}

// 解析完整日期格式 "2025年4月6日 03:34"
func (u *TimeHelperUtil) parseFullDateFormat(input string) (*gtime.Time, error) {
	// 匹配 "2025年4月6日 03:34" 或 "2025年4月6日03:34"
	re := regexp.MustCompile(`^(\d{4})年(\d{1,2})月(\d{1,2})日\s*(\d{1,2}):(\d{1,2})$`)
	matches := re.FindStringSubmatch(input)
	if len(matches) != 6 {
		return nil, fmt.Errorf("not a full date format")
	}

	// 构造时间字符串 "2025-04-06 03:34"
	timeStr := fmt.Sprintf("%s-%02s-%02s %s:%s",
		matches[1], matches[2], matches[3], matches[4], matches[5])

	t, err := gtime.StrToTime(timeStr, "Y-m-d H:i")
	if err != nil {
		return nil, fmt.Errorf("invalid full date format: %v", err)
	}

	return t, nil
}

// 解析带星期的格式 "周六10:01"
func (u *TimeHelperUtil) parseWeekdayTimeFormat(input string) (*gtime.Time, error) {
	weekdayMap := map[string]time.Weekday{
		"周日": time.Sunday,
		"周一": time.Monday,
		"周二": time.Tuesday,
		"周三": time.Wednesday,
		"周四": time.Thursday,
		"周五": time.Friday,
		"周六": time.Saturday,
	}

	for k := range weekdayMap {
		if strings.HasPrefix(input, k) {
			timeStr := input[len(k):]
			parsedTime, err := gtime.StrToTime(timeStr, "H:i")
			if err != nil {
				return nil, fmt.Errorf("invalid time format: %v", err)
			}

			now := gtime.Now()
			targetWeekday := weekdayMap[k]
			currentWeekday := now.Weekday()

			// 计算上一个匹配的星期X
			daysToSubtract := (int(currentWeekday) - int(targetWeekday) + 7) % 7
			if daysToSubtract == 0 {
				// 如果是今天，检查时间是否已过
				if parsedTime.Hour() > now.Hour() ||
					(parsedTime.Hour() == now.Hour() && parsedTime.Minute() > now.Minute()) {
					daysToSubtract = 7 // 如果时间未到，用上周的
				}
			}

			targetDate := now.AddDate(0, 0, -daysToSubtract)
			return gtime.NewFromTime(time.Date(
				targetDate.Year(),
				time.Month(targetDate.Month()),
				targetDate.Day(),
				parsedTime.Hour(),
				parsedTime.Minute(),
				0, 0,
				time.Local,
			)), nil
		}
	}

	return nil, fmt.Errorf("not a weekday format")
}

// 解析纯时间格式 "10:01"
func (u *TimeHelperUtil) parseSimpleTimeFormat(input string) (*gtime.Time, error) {
	// 检查是否是纯时间格式
	re := regexp.MustCompile(`^\d{1,2}:\d{1,2}$`)
	if !re.MatchString(input) {
		return nil, fmt.Errorf("not a simple time format")
	}

	parsedTime, err := gtime.StrToTime(input, "H:i")
	if err != nil {
		return nil, fmt.Errorf("invalid time format: %v", err)
	}

	now := gtime.Now()
	return gtime.NewFromTime(time.Date(
		now.Year(),
		time.Month(now.Month()),
		now.Day(),
		parsedTime.Hour(),
		parsedTime.Minute(),
		0, 0,
		time.Local,
	)), nil
}
