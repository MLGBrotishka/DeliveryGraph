package handlers

import (
	"fmt"
	"time"
)

func GetTimeValue(inputTime string) float64 {
	parsedTime, err := time.Parse("2006-01-02 15:04:05 -0700 MST", inputTime)
	if err != nil {
		fmt.Println("Ошибка разбора времени:", err)
		return -1
	}

	hour := parsedTime.Hour()

	if (hour >= 8 && hour < 10) || (hour > 17 && hour < 21) {
		return 1.2
	}

	return 1
}

func FormatTime(inputTime string) (string, error) {
	parsedTime, err := time.Parse("15:04:05", inputTime)
	if err != nil {
		return "", err
	}

	currentTime := time.Now()
	formattedTime := time.Date(
		currentTime.Year(), currentTime.Month(), currentTime.Day(),
		parsedTime.Hour(), parsedTime.Minute(), parsedTime.Second(), parsedTime.Nanosecond(),
		currentTime.Location(),
	).Format("2006-01-02 15:04:05 -0700 MST")

	return formattedTime, nil
}
