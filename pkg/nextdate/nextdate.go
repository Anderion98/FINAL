package nextdate

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const TimeFormat = "20060102"

func afterNow(date, now time.Time) bool {
	nowRounding := now.Truncate(24 * time.Hour)
	dateRounding := date.Truncate(24 * time.Hour)
	return dateRounding.After(nowRounding) || nowRounding.Equal(dateRounding)
}

// правило для d
func repeatD(now time.Time, date time.Time, repeatParams []string) (string, error) {
	if len(repeatParams) != 2 {
		return "", errors.New("некорректное правило повторения")
	}
	day, err := strconv.Atoi(repeatParams[1])
	if err != nil {
		return "", errors.New("некорректное правило повторения")
	}
	if day < 1 || day > 400 {
		return "", errors.New("превышен максимально допустимый интервал")
	}
	date = date.AddDate(0, 0, day)
	for afterNow(now, date) {
		date = date.AddDate(0, 0, day)
	}
	return date.Format(TimeFormat), nil
}

// правило для y
func repeatY(now time.Time, date time.Time, repeatParams []string) (string, error) {
	if len(repeatParams) != 1 {
		return "", errors.New("некорректное правило повторения")
	}
	date = date.AddDate(1, 0, 0)
	for afterNow(now, date) {
		date = date.AddDate(1, 0, 0)
	}
	return date.Format(TimeFormat), nil
}

// правило для w и m
func repeatWM() (string, error) {
	return "", fmt.Errorf("неподдерживаемый формат")
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("правило повтора не указано")
	}
	date, err := time.Parse(TimeFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("недопустимый формат даты")
	}
	component := strings.Split(repeat, " ")
	switch component[0] {
	case "d":
		return repeatD(now, date, component)
	case "y":
		return repeatY(now, date, component)
	case "w":
		return repeatWM()

	case "m":
		return repeatWM()
	default:
		return "", errors.New("некорректное правило повторения")
	}
}
