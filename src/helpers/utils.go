package helpers

import (
	"spl-access/src/model"
	"strings"
	"time"
)

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

func MaskAccessData(accesses *[]model.Access) *[]model.Access {
	maskedAccesses := make([]model.Access, len(*accesses))

	for i, access := range *accesses {
		maskedAccesses[i] = access // Copy the struct
		maskedAccesses[i].Run = maskRun(access.Run)
		maskedAccesses[i].FullName = maskFullName(access.FullName)
	}

	return &maskedAccesses
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
