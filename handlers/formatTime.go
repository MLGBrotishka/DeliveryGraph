package main

import (
	"fmt"
	"time"
)

func formatTime(inputTime string) (string, error) {
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

func main() {
//	inputTime := "12:45:15"
//	formattedTime, err := formatTime(inputTime)
//	if err != nil {
//		fmt.Println("Ошибка:", err)
		return
	}

//	fmt.Println(formattedTime)
}
