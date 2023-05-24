package main

import (
	"encoding/json"
	"fmt"
	"github.com/SevereCloud/vksdk/v2/marusia"
	"github.com/golang-module/carbon/v2"
	"log"
	"net/http"
	"spbtoday/convert"
	"spbtoday/sender"
	"spbtoday/utils"
)

type myPayload struct {
	Text string
	marusia.DefaultPayload
}

const (
	dm = "d.m"
)

var (
	incidents    = []string{" произошло вот что...\n", " было вот что...\n", " случилось такое событие...\n"}
	otherPreDate = "Про какое событие вы хотите узнать подробнее? Скажите: «про первое» или «про второе», или назовите другую дату для продолжения."
	goodBye      = []string{"До новых встреч!", "До свидания!", "Всего доброго!", "Всего доброго! Хорошего дня!"}
	pen, mayak   = &sender.Penevin{}, &sender.Mayakovsky{}
	day, mouth   = "", ""
	wh           = marusia.NewWebhook()
	err          error
)

func main() {
	wh.EnableDebuging()
	wh.OnEvent(func(r marusia.Request) (resp marusia.Response) {
		if r.Session.New {
			resp.TextArray = []string{`Я могу рассказывать о примечательных датах, которые произошли в истории великого города.
Достаточно назвать интересующую дату в формате "первое апреля", и я расскажу, что было в истории Санкт-Петербурга в этот день!
Если Вы уже знаете все события, можно сказать "Хватит".`, `Какая дата Вас интересует?`}
			resp.AddButton("Сегодня", myPayload{
				Text: "Сегодня",
			})
			resp.AddButton("Завтра", myPayload{
				Text: "Завтра",
			})
			resp.AddButton("Вчера", myPayload{
				Text: "Вчера",
			})
			return
		}

		switch r.Request.Type {
		case marusia.SimpleUtterance:
			switch r.Request.Command {
			case "хватит":
				answer := utils.RandAnswer(goodBye)
				resp.Text = answer
				resp.TTS = answer
				resp.EndSession = true
			case marusia.OnInterrupt:
				answer := utils.RandAnswer(goodBye)
				resp.Text = answer
				resp.TTS = answer
				resp.EndSession = true
				return
			default:
				if utils.CheckCommand(r.Request.Command) {
					day, mouth, err = convert.ConverterToNumericDate(r.Request.Command)
					if err != nil {
						resp.Text = err.Error()
						resp.TTS = err.Error()
						return
					}

					pen, mayak, err = sender.SendRequests(day, mouth)
					if err != nil {
						resp.Text = err.Error()
						resp.TTS = err.Error()
						return
					}

					inc := fmt.Sprintf("%s %s \n%s\n \n Источник:  %s \n\n %s\n\n Источник: %s\n\n", convert.ConverterToRussianTextDate(pen.Date), utils.RandAnswer(incidents), mayak.Title, mayak.Source, pen.Title, pen.Source)
					response := inc + otherPreDate

					resp.TextArray = []string{response}

					resp.AddButton("Первый", myPayload{
						Text: "первый",
					})
					resp.AddButton("Второй", myPayload{
						Text: "второй",
					})
					return
				} else {
					switch r.Request.Command {
					case "сегодня":
						date := carbon.Now().Format(dm)
						day, mouth = convert.ConverterCarbonDate(date)
						pen, mayak, err = sender.SendRequests(day, mouth)
						if err != nil {
							resp.Text = err.Error()
							resp.TTS = err.Error()
							return
						}
						inc := fmt.Sprintf("%s %s \n%s\n \n Источник:  %s \n\n %s\n\n Источник: %s\n\n", convert.ConverterToRussianTextDate(pen.Date), utils.RandAnswer(incidents), mayak.Title, mayak.Source, pen.Title, pen.Source)
						response := inc + otherPreDate

						resp.TextArray = []string{response}

						resp.AddButton("Первый", myPayload{
							Text: "первый",
						})
						resp.AddButton("Второй", myPayload{
							Text: "второй",
						})

						resp.AddButton("Хватит", myPayload{
							Text: "Хватит",
						})
						return
					case "завтра":
						date := carbon.Tomorrow().Format(dm)
						day, mouth = convert.ConverterCarbonDate(date)
						pen, mayak, err = sender.SendRequests(day, mouth)
						if err != nil {
							resp.Text = err.Error()
							resp.TTS = err.Error()
						}
						inc := fmt.Sprintf("%s %s \n%s\n \n Источник:  %s \n\n %s\n\n Источник: %s\n\n", convert.ConverterToRussianTextDate(pen.Date), utils.RandAnswer(incidents), mayak.Title, mayak.Source, pen.Title, pen.Source)
						response := inc + otherPreDate

						resp.TextArray = []string{response}

						resp.AddButton("Первый", myPayload{
							Text: "первый",
						})
						resp.AddButton("Второй", myPayload{
							Text: "второй",
						})
						return
					case "вчера":
						date := carbon.Yesterday().Format(dm)
						day, mouth = convert.ConverterCarbonDate(date)
						pen, mayak, err = sender.SendRequests(day, mouth)
						if err != nil {
							resp.Text = err.Error()
							resp.TTS = err.Error()
						}
						inc := fmt.Sprintf("%s %s \n%s\n \n Источник:  %s \n\n %s\n\n Источник: %s\n\n", convert.ConverterToRussianTextDate(pen.Date), utils.RandAnswer(incidents), mayak.Title, mayak.Source, pen.Title, pen.Source)

						response := inc + otherPreDate

						resp.TextArray = []string{response}

						resp.AddButton("Первый", myPayload{
							Text: "первый",
						})
						resp.AddButton("Второй", myPayload{
							Text: "второй",
						})
						return
					default:
						if pen.Title != "" || mayak.Title != "" {
							switch r.Request.Command {
							case "первое", "1", "первый", "один", "про первое":
								inc := fmt.Sprintf("%s\n%s\n \n Источник: %s\n\n", mayak.Title, mayak.Description, mayak.Source)
								result := utils.DivideString(inc)
								result = append(result, "Узнайте подробнее о втором событии, сказав: \"второе\", или скажите другую дату")
								resp.TextArray = result

								resp.AddButton("Завтра", myPayload{
									Text: "Завтра",
								})

								resp.AddButton("Хватит", myPayload{
									Text: "Хватит",
								})
								return

							case "второе", "2", "два", "про второе", "второй":
								inc := fmt.Sprintf("%s\n%s\n \n Источник: %s\n\n", pen.Title, pen.Text, pen.Source)
								result := utils.DivideString(inc)
								result = append(result, "Узнайте подробнее о первом событии, сказав: \"первое\", или скажите другую дату")
								resp.TextArray = result

								resp.AddButton("Завтра", myPayload{
									Text: "Завтра",
								})

								resp.AddButton("Хватит", myPayload{
									Text: "Хватит",
								})
								return
							}
						} else {
							resp.Text = "Мне непонятно, извините. Скажите или напишите мне день и месяц в формате \"двадцать второе февраля\", и я расскажу, что было в истории Санкт-Петербурга в этот день"
							resp.TTS = "Мне непонятно, извините. Скажите или напишите мне день и месяц в формате \"двадцать второе февраля\", и я расскажу, что было в истории Санкт-Петербурга в этот день"
							return
						}
					}
				}
			}
		case marusia.ButtonPressed:
			var p myPayload

			err = json.Unmarshal(r.Request.Payload, &p)
			if err != nil {
				resp.Text = err.Error()
				return
			}

			var date string

			switch p.Text {
			case "Сегодня":
				date = carbon.Now().Format(dm)
				day, mouth = convert.ConverterCarbonDate(date)
				pen, mayak, err = sender.SendRequests(day, mouth)
				if err != nil {
					resp.Text = err.Error()
					resp.TTS = err.Error()
					return
				}
				inc := fmt.Sprintf("%s %s \n%s\n \n Источник:  %s \n\n %s\n\n Источник: %s\n\n", convert.ConverterToRussianTextDate(pen.Date), utils.RandAnswer(incidents), mayak.Title, mayak.Source, pen.Title, pen.Source)
				response := inc + otherPreDate

				resp.TextArray = []string{response}

				resp.AddButton("Первый", myPayload{
					Text: "первый",
				})
				resp.AddButton("Второй", myPayload{
					Text: "второй",
				})

				resp.AddButton("Хватит", myPayload{
					Text: "хватит",
				})
				return
			case "Завтра":
				date = carbon.Tomorrow().Format(dm)
				day, mouth = convert.ConverterCarbonDate(date)
				pen, mayak, err = sender.SendRequests(day, mouth)
				if err != nil {
					resp.Text = err.Error()
					resp.TTS = err.Error()
				}
				inc := fmt.Sprintf("%s %s \n%s\n \n Источник:  %s \n\n %s\n\n Источник: %s\n\n", convert.ConverterToRussianTextDate(pen.Date), utils.RandAnswer(incidents), mayak.Title, mayak.Source, pen.Title, pen.Source)
				response := inc + otherPreDate

				resp.TextArray = []string{response}

				resp.AddButton("Первый", myPayload{
					Text: "первый",
				})
				resp.AddButton("Второй", myPayload{
					Text: "второй",
				})
				return
			case "Вчера":
				date = carbon.Yesterday().Format(dm)
				day, mouth = convert.ConverterCarbonDate(date)
				pen, mayak, err = sender.SendRequests(day, mouth)
				if err != nil {
					resp.Text = err.Error()
					resp.TTS = err.Error()
				}
				inc := fmt.Sprintf("%s %s \n%s\n \n Источник:  %s \n\n %s\n\n Источник: %s\n\n", convert.ConverterToRussianTextDate(pen.Date), utils.RandAnswer(incidents), mayak.Title, mayak.Source, pen.Title, pen.Source)
				response := inc + otherPreDate

				resp.TextArray = []string{response}

				resp.AddButton("Первый", myPayload{
					Text: "первый",
				})
				resp.AddButton("Второй", myPayload{
					Text: "второй",
				})
				return
			case "первое", "1", "первый", "один", "про первое":
				inc := fmt.Sprintf("%s\n%s\n \n Источник: %s\n\n", mayak.Title, mayak.Description, mayak.Source)
				result := utils.DivideString(inc)
				result = append(result, "Если Вас интересуют второе событие этого дня, скажите: \"второе\" или назовите другую дату.")
				resp.TextArray = result

				resp.AddButton("Завтра", myPayload{
					Text: "Завтра",
				})

				resp.AddButton("Хватит", myPayload{
					Text: "Хватит",
				})
				return
			case "второе", "2", "два", "про второе", "второй":
				inc := fmt.Sprintf("%s\n%s\n \n Источник: %s\n\n", pen.Title, pen.Text, pen.Source)
				result := utils.DivideString(inc)
				result = append(result, "Если Вас интересуют пеовое событие этого дня, скажите: \"первое\" или назовите другую дату.")
				resp.TextArray = result

				resp.AddButton("Завтра", myPayload{
					Text: "Завтра",
				})

				resp.AddButton("Хватит", myPayload{
					Text: "Хватит",
				})
				return

			case "хватит":
				answer := utils.RandAnswer(goodBye)
				resp.Text = answer
				resp.TTS = answer
				resp.EndSession = true
				return
			}
		}
		return
	})

	http.HandleFunc("/spbtoday", wh.HandleFunc)

	if err := http.ListenAndServe(":8026", nil); err != nil {
		log.Fatal(err)
	}
}
