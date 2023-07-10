package main

import (
	"fmt"
	"time"
)

func getTimeValue(inputTime string) float64 {
	parsedTime, err := time.Parse("2006-01-02 15:04:05 -0700 MST", inputTime)
	if err != nil {
		fmt.Println("Ошибка разбора времени:", err)
		return -1
	}

	hour := parsedTime.Hour()
	minute := parsedTime.Minute()

	if (hour >= 8 && hour < 10) || (hour == 17 && minute >= 30) || (hour > 17 && hour < 21) {
		return 1.2
	}

	return 1
}

func main() {
	//inputTime := "2023-07-09 18:34:56 +0000 UTC"
	//value := getTimeValue(inputTime)
	//fmt.Println("Коэффициент цены:", value)
}