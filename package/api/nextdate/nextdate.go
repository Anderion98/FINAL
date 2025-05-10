package nextdate

import (
	"errors"
	"fmt"
	"gofer/package/other"
	"net/http"
	"strconv"
	"strings"
	"time"
)

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
	for date.Before(now) {
		date = date.AddDate(0, 0, day)
	}
	return date.Format(other.TimeFormat), nil
}

func repeatY(now time.Time, date time.Time, repeatParams []string) (string, error) {
	if len(repeatParams) != 1 {
		return "", errors.New("некорректное правило повторения")
	}
	date = date.AddDate(1, 0, 0)
	for date.Before(now) {
		date = date.AddDate(1, 0, 0)
	}
	return date.Format(other.TimeFormat), nil
}

func repeatWM() (string, error) {
	return "", fmt.Errorf("неподдерживаемый формат")
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("правило повтора не указано")
	}

	date, err := time.Parse(other.TimeFormat, dstart)
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
		return "", errors.New("формат правила повторения не соблюдён")
	}
}
func NextDayHandler(w http.ResponseWriter, r *http.Request) {
	now, err := time.Parse(other.TimeFormat, r.FormValue("now"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	date := r.FormValue("date")
	repeat := r.FormValue("repeat")
	res, err := NextDate(now, date, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, res)
}
