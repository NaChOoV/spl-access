package helpers

import (
	"spl-access/src/model"
	"strings"
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
