package sender

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

type Penevin struct {
	Title   string `json:"title"`
	Date    string `json:"date"`
	DateTxt string `json:"date_txt"`
	Text    string `json:"text"`
	URL     string `json:"url"`
	Source  string `json:"source"`
}

type Mayakovsky struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Date        string `json:"date"`
	Description string `json:"description"`
	Source      string `json:"source"`
}

var sources = []struct {
	id     int
	source string
	url    string
}{
	{
		id:     0,
		source: "Петербургский календарь Панёвина",
		url:    "http://api.panevin.ru/v1/?date=",
	},
	{
		id:     1,
		source: "Библиотека Маяковского",
		url:    "http://84.201.142.84:8005/memorable_dates/date/day/{day}/month/{month}",
	},
}

func SendRequests(day, month string) (*Penevin, *Mayakovsky, error) {
	var wg sync.WaitGroup

	// Создаем каналы для результатов запросов
	result0 := make(chan []byte, 1)
	result1 := make(chan []byte, 1)

	wg.Add(len(sources))

	for _, s := range sources {
		go func(s struct {
			id     int
			source string
			url    string
		}) {
			defer wg.Done()

			var url string

			switch s.id {
			case 0:
				url = fmt.Sprintf("%s%s.%s", s.url, day, month)
			case 1:
				newURL := strings.Replace(s.url, "{day}", day, -1)
				url = strings.Replace(newURL, "{month}", month, -1)
			}

			// Отправляем запрос и обрабатываем ошибки
			resp, err := http.Get(url)
			if err != nil {
				log.Printf("%s: %v\n", s.source, err)
				return
			}
			defer resp.Body.Close()

			// Читаем ответ и отправляем в соответствующий канал
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("%s: %v\n", s.source, err)
				return
			}

			switch s.id {
			case 0:
				result0 <- body
			case 1:
				result1 <- body
			}
		}(s)
	}

	go func() {
		wg.Wait()
		close(result0)
		close(result1)
	}()

	penevin, mayakovsky, err := ProcessingRequestData(result0, result1)
	if err != nil {
		return nil, nil, err
	}

	return penevin, mayakovsky, nil

}

func ProcessingRequestData(result0, result1 chan []byte) (*Penevin, *Mayakovsky, error) {
	var p Penevin
	mlist := []Mayakovsky{}

	//fmt.Println(string(<-result1), string(<-result0))
	select {
	case res0, ok := <-result0:
		if ok {
			regex := regexp.MustCompile(`<[^>]*>|\\n`)
			result := regex.ReplaceAll(res0, []byte{})
			if err := json.Unmarshal(result, &p); err != nil {
				return nil, nil, fmt.Errorf("Error parsing result0: %s\n", err.Error())
			}
		}
	case <-time.After(5 * time.Second):
		return nil, nil, fmt.Errorf("timeout while reading result0")
	}

	select {
	case res1, ok := <-result1:
		if ok {
			regex := regexp.MustCompile(`<[^>]*>|\\n`)
			result := regex.ReplaceAll(res1, []byte{})
			if err := json.Unmarshal(result, &mlist); err != nil {
				return nil, nil, fmt.Errorf("Error parsing result1: %s\n", err.Error())
			}
		}
	case <-time.After(5 * time.Second):
		return nil, nil, fmt.Errorf("timeout while reading result1")
	}

	p.Source = "Петербургский календарь Панёвина"
	mlist[0].Source = "Библиотека Маяковского"

	return &p, &mlist[0], nil

}
