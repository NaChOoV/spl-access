package helpers

import (
	"fmt"
	"spl-access/src/model"
	"strconv"
	"strings"
	"time"
)

// Scheduler defines the working hours for each day of the week
// Monday=1, Tuesday=2, ..., Sunday=7
var scheduler = map[int][2]string{
	1: {"06:30", "23:00"}, // Monday
	2: {"06:30", "23:00"}, // Tuesday
	3: {"06:30", "23:00"}, // Wednesday
	4: {"06:30", "23:00"}, // Thursday
	5: {"06:30", "23:00"}, // Friday
	6: {"09:00", "20:00"}, // Saturday
	7: {"09:00", "14:00"}, // Sunday
}

func maskString(str string, keepStart, keepEnd int) string {
	runes := []rune(str)
	length := len(runes)

	if length <= keepStart+keepEnd {
		return str
	}

	masked := make([]rune, length)
	for i := 0; i < length; i++ {
		if i < keepStart || i >= length-keepEnd {
			masked[i] = runes[i]
		} else {
			masked[i] = '*'
		}
	}
	return string(masked)
}

func maskRun(run string) string {
	parts := strings.Split(run, "-")
	if len(parts) != 2 {
		return run
	}
	mainPart := maskString(parts[0], 2, 0)
	return mainPart + "-" + parts[1]
}

func maskFullName(fullName string) string {
	words := strings.Fields(fullName)
	maskedWords := make([]string, len(words))

	for i, word := range words {
		maskedWords[i] = maskString(word, 1, 1)
	}

	return strings.Join(maskedWords, " ")
}

func MaskAccessData(accesses []*model.Access) []*model.Access {
	maskedAccesses := make([]*model.Access, len(accesses))

	for i, access := range accesses {
		maskedAccess := *access
		maskedAccess.Run = maskRun(access.Run)
		maskedAccess.FullName = maskFullName(access.FullName)
		maskedAccesses[i] = &maskedAccess
	}

	return maskedAccesses
}

func IsChileSleepTime(utcTime time.Time, zone string) bool {
	var location *time.Location
	switch zone {
	case "GMT-3":
		location = time.FixedZone("GMT-3", -3*60*60)
	case "GMT-4":
		location = time.FixedZone("GMT-4", -4*60*60)
	}

	chileTime := utcTime.In(location)
	wd, hr, min := chileTime.Weekday(), chileTime.Hour(), chileTime.Minute()

	switch wd {
	case time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday:
		if ((hr > 6) || (hr == 6 && min >= 30)) && (hr < 23) {
			return false
		}
	case time.Saturday:
		if hr >= 9 && hr < 20 {
			return false
		}
	case time.Sunday:
		if hr >= 9 && hr < 14 {
			return false
		}
	}

	return true
}

func CheckSleepTime(timeZone string) {
	sleepSeconds, err := GetSleepSeconds(timeZone)
	if err != nil {
		fmt.Printf("Error getting sleep seconds: %v\n", err)
		return
	}

	if sleepSeconds > 0 {
		hours := sleepSeconds / 3600
		fmt.Printf("[CRON] Sleeping for %d seconds (%d hours)\n", sleepSeconds, hours)
		time.Sleep(time.Duration(sleepSeconds) * time.Second)
	}
}

func GetSleepSeconds(timeZone string) (int64, error) {
	now := time.Now()

	var chileOffset int
	switch timeZone {
	case "GMT-3":
		chileOffset = -3
	case "GMT-4":
		chileOffset = -4
	default:
		return 0, fmt.Errorf("unsupported timezone: %s", timeZone)
	}

	// Create a time in the Chile timezone using the timezone offset
	location := time.FixedZone("Chile", chileOffset*60*60)
	chileTime := now.In(location)

	// Convert Go weekday (Sunday=0) to our scheduler format (Monday=1, Sunday=7)
	day := int(chileTime.Weekday())
	if day == 0 {
		day = 7
	}

	schedule, exists := scheduler[day]
	if !exists {
		return 0, fmt.Errorf("no schedule found for day %d", day)
	}

	startStr, endStr := schedule[0], schedule[1]

	// Parse start time
	startParts := strings.Split(startStr, ":")
	sh, err := strconv.Atoi(startParts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid start hour: %w", err)
	}
	sm, err := strconv.Atoi(startParts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid start minute: %w", err)
	}

	// Parse end time
	endParts := strings.Split(endStr, ":")
	eh, err := strconv.Atoi(endParts[0])
	if err != nil {
		return 0, fmt.Errorf("invalid end hour: %w", err)
	}
	em, err := strconv.Atoi(endParts[1])
	if err != nil {
		return 0, fmt.Errorf("invalid end minute: %w", err)
	}

	// Create start and end times for today
	start := time.Date(chileTime.Year(), chileTime.Month(), chileTime.Day(), sh, sm, 0, 0, location)
	end := time.Date(chileTime.Year(), chileTime.Month(), chileTime.Day(), eh, em, 0, 0, location)

	// Check if we're currently in the active period
	if chileTime.After(start) && chileTime.Before(end) {
		return 0, nil // Currently active, no sleep needed
	} else if chileTime.Before(start) {
		// Before start time today
		return int64(start.Sub(chileTime).Seconds()), nil
	} else {
		// After end time today, calculate next day's start
		nextDay := (day % 7) + 1
		nextSchedule, exists := scheduler[nextDay]
		if !exists {
			return 0, fmt.Errorf("no schedule found for next day %d", nextDay)
		}

		nextStartStr := nextSchedule[0]
		nextStartParts := strings.Split(nextStartStr, ":")
		nsh, err := strconv.Atoi(nextStartParts[0])
		if err != nil {
			return 0, fmt.Errorf("invalid next start hour: %w", err)
		}
		nsm, err := strconv.Atoi(nextStartParts[1])
		if err != nil {
			return 0, fmt.Errorf("invalid next start minute: %w", err)
		}

		// Create next day's start time
		nextStart := time.Date(chileTime.Year(), chileTime.Month(), chileTime.Day()+1, nsh, nsm, 0, 0, location)
		return int64(nextStart.Sub(chileTime).Seconds()), nil
	}
}

// KeyExtractorFunc defines a function type for extracting a key from any type
type KeyExtractorFunc[T any] func(item T) string

// RemoveDuplicatesGeneric is a generic function to remove duplicates from any slice
// based on a key extraction function
func RemoveDuplicatesGeneric[T any](items []T, keyExtractor KeyExtractorFunc[T]) []T {
	seen := make(map[string]bool)
	var result []T

	for _, item := range items {
		key := keyExtractor(item)
		if !seen[key] {
			seen[key] = true
			result = append(result, item)
		}
	}

	return result
}
