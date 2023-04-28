package convert

import (
	"errors"
	"fmt"
	"strings"
)

var months = map[string]string{
	"января":   "01",
	"февраля":  "02",
	"марта":    "03",
	"апреля":   "04",
	"мая":      "05",
	"июня":     "06",
	"июля":     "07",
	"августа":  "08",
	"сентября": "09",
	"октября":  "10",
	"ноября":   "11",
	"декабря":  "12",
}

var days = map[string]string{
	"первое":             "01",
	"второе":             "02",
	"третье":             "03",
	"четвертое":          "04",
	"пятое":              "05",
	"шестое":             "06",
	"седьмое":            "07",
	"восьмое":            "08",
	"девятое":            "09",
	"десятое":            "10",
	"одиннадцатое":       "11",
	"двенадцатое":        "12",
	"тринадцатое":        "13",
	"четырнадцатое":      "14",
	"пятнадцатое":        "15",
	"шестнадцатое":       "16",
	"семнадцатое":        "17",
	"восемнадцатое":      "18",
	"девятнадцатое":      "19",
	"двадцатое":          "20",
	"двадцать первое":    "21",
	"двадцать второе":    "22",
	"двадцать третье":    "23",
	"двадцать четвертое": "24",
	"двадцать пятое":     "25",
	"двадцать шестое":    "26",
	"двадцать седьмое":   "27",
	"двадцать восьмое":   "28",
	"двадцать девятое":   "29",
	"тридцатое":          "30",
	"тридцать первое":    "31",
}

func ConverterToNumericDate(str string) (day, month string, err error) {
	split := strings.Split(str, " ")
	var dayStr, monthStr string
	switch len(split) {
	case 3:
		dayStr, monthStr = split[0]+" "+split[1], split[2]
	case 2:
		dayStr, monthStr = split[0], split[1]
	default:
		return "", "", errors.New("неправильный формат даты")
	}

	day, ok := days[dayStr]
	if !ok {
		return "", "", fmt.Errorf("неправильный формат дня: %s", dayStr)
	}
	month, ok = months[monthStr]
	if !ok {
		return "", "", fmt.Errorf("неправильный формат месяца: %s", monthStr)
	}

	return day, month, nil
}

func ConverterCarbonDate(dm string) (string, string) {
	date := strings.Split(dm, ".")
	d := date[0]
	m := date[1]
	return d, m
}
