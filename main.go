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
	otherPreDate = []string{"Если Вас интересуют другие события этого дня, скажите: \"ещё события\".\n", "Если Вас интересуют другие события этого дня, скажите: \"ещё события\", либо назовите другую дату."}
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

					inc := fmt.Sprintf("%s %s \n%s\n \n Источник:  %s\n\n", convert.ConverterToRussianTextDate(pen.Date), utils.RandAnswer(incidents), mayak.Title, mayak.Source)
					anDate := utils.RandAnswer(otherPreDate)
					response := inc + anDate

					resp.TextArray = []string{response, "Чтобы узнать подробности источника, скажите: Маяковский"}

					resp.AddButton("Ёще", myPayload{
						Text: "еще",
					})
					resp.AddButton("Маяковский", myPayload{
						Text: "маяковский",
					})
					return
				} else {
					switch r.Request.Command {
					case "еще", "ещё", "еще событие", "ещё событие", "еще события", "ещё события":
						inc := fmt.Sprintf("%s %s \n%s\n \n Источник:  %s\n\n", convert.ConverterToRussianTextDate(pen.Date), utils.RandAnswer(incidents), pen.Title, pen.Source)
						fulldata := "Чтобы узнать подробности источника, скажите: Панёвин"
						response := inc
						resp.TextArray = []string{response, fulldata}

						resp.AddButton("Панёвин", myPayload{
							Text: "панёвин",
						})
						return

					case "сегодня":
						date := carbon.Now().Format(dm)
						day, mouth = convert.ConverterCarbonDate(date)
						pen, mayak, err = sender.SendRequests(day, mouth)
						if err != nil {
							resp.Text = err.Error()
							resp.TTS = err.Error()
							return
						}
						inc := fmt.Sprintf("%s %s \n%s\n \n Источник:  %s\n\n", convert.ConverterToRussianTextDate(pen.Date), utils.RandAnswer(incidents), mayak.Title, mayak.Source)
						anDate := utils.RandAnswer(otherPreDate)
						response := inc + anDate

						resp.TextArray = []string{response, "Чтобы узнать подробности источника, скажите: Маяковский"}

						resp.AddButton("Ёще", myPayload{
							Text: "еще",
						})
						resp.AddButton("Маяковский", myPayload{
							Text: "маяковский",
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
						inc := fmt.Sprintf("%s %s \n%s\n \n Источник:  %s\n\n", convert.ConverterToRussianTextDate(pen.Date), utils.RandAnswer(incidents), mayak.Title, mayak.Source)
						anDate := utils.RandAnswer(otherPreDate)
						response := inc + anDate

						resp.TextArray = []string{response, "Чтобы узнать подробности источника, скажите: Маяковский"}

						resp.AddButton("Ёще", myPayload{
							Text: "еще",
						})
						resp.AddButton("Маяковский", myPayload{
							Text: "маяковский",
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
						inc := fmt.Sprintf("%s %s \n%s\n \n Источник:  %s\n\n", convert.ConverterToRussianTextDate(pen.Date), utils.RandAnswer(incidents), mayak.Title, mayak.Source)

						anDate := utils.RandAnswer(otherPreDate)
						response := inc + anDate

						resp.TextArray = []string{response, "Чтобы узнать подробности источника, скажите: Маяковский"}

						resp.AddButton("Ёще", myPayload{
							Text: "еще",
						})
						resp.AddButton("Маяковский", myPayload{
							Text: "маяковский",
						})
						return
					default:
						if pen.Title != "" || mayak.Title != "" {
							switch r.Request.Command {
							case "паневин", "панёвин":
								inc := fmt.Sprintf("%s\n%s\n \n Источник: %s\n\n", pen.Title, pen.Text, pen.Source)
								result := utils.DivideString(inc)
								result = append(result, "Чтобы узнать другие события, назовите другую дату или скажите \"ещё события\"")
								resp.TextArray = result

								resp.AddButton("Завтра", myPayload{
									Text: "Завтра",
								})

								resp.AddButton("Хватит", myPayload{
									Text: "Хватит",
								})
								return
							case "маяковский":
								inc := fmt.Sprintf("%s\n%s\n \n Источник: %s\n\n", mayak.Title, mayak.Description, mayak.Source)
								result := utils.DivideString(inc)
								result = append(result, "Чтобы узнать другие события, назовите другую дату")
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
			case "еще":
				inc := fmt.Sprintf("%s %s \n%s\n \n Источник:  %s\n\n", convert.ConverterToRussianTextDate(pen.Date), utils.RandAnswer(incidents), pen.Title, pen.Source)
				anDate := `Если интересуют другие события, назовите другую дату`
				response := inc + anDate
				resp.TextArray = []string{response, "Чтобы узнать подробности источника, скажите: Панёвин"}

				resp.AddButton("Панёвин", myPayload{
					Text: "панёвин",
				})

			case "Сегодня":
				date = carbon.Now().Format(dm)
				day, mouth = convert.ConverterCarbonDate(date)
				pen, mayak, err = sender.SendRequests(day, mouth)
				if err != nil {
					resp.Text = err.Error()
					resp.TTS = err.Error()
					return
				}
				inc := fmt.Sprintf("%s %s \n%s\n \n Источник:  %s\n\n", convert.ConverterToRussianTextDate(pen.Date), utils.RandAnswer(incidents), mayak.Title, mayak.Source)
				anDate := utils.RandAnswer(otherPreDate)
				response := inc + anDate

				resp.TextArray = []string{response, "Чтобы узнать подробности источника, скажите: Маяковский"}

				resp.AddButton("Ёще", myPayload{
					Text: "еще",
				})
				resp.AddButton("Маяковский", myPayload{
					Text: "маяковский",
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
				inc := fmt.Sprintf("%s %s \n%s\n \n Источник:  %s\n\n", convert.ConverterToRussianTextDate(pen.Date), utils.RandAnswer(incidents), mayak.Title, mayak.Source)
				anDate := utils.RandAnswer(otherPreDate)
				response := inc + anDate

				resp.TextArray = []string{response, "Чтобы узнать подробности источника, скажите: Маяковский"}

				resp.AddButton("Ёще", myPayload{
					Text: "еще",
				})
				resp.AddButton("Маяковский", myPayload{
					Text: "маяковский",
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
				inc := fmt.Sprintf("%s %s \n%s\n \n Источник:  %s\n\n", convert.ConverterToRussianTextDate(pen.Date), utils.RandAnswer(incidents), mayak.Title, mayak.Source)
				anDate := utils.RandAnswer(otherPreDate)
				response := inc + anDate

				resp.TextArray = []string{response, "Чтобы узнать подробности источника, скажите: Маяковский"}

				resp.AddButton("Ёще", myPayload{
					Text: "еще",
				})
				resp.AddButton("Маяковский", myPayload{
					Text: "маяковский",
				})
				return
			case "панёвин":
				inc := fmt.Sprintf("%s\n%s\n \n Источник: %s\n\n", pen.Title, pen.Text, pen.Source)
				result := utils.DivideString(inc)
				resp.TextArray = result

				resp.AddButton("Завтра", myPayload{
					Text: "Завтра",
				})

				resp.AddButton("Хватит", myPayload{
					Text: "Хватит",
				})
				return
			case "маяковский":
				inc := fmt.Sprintf("%s\n%s\n \n Источник: %s\n\n", mayak.Title, mayak.Description, mayak.Source)
				result := utils.DivideString(inc)
				result = append(result, "Если Вас интересуют другие события этого дня, скажите: \"ещё события\".")
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

	if err := http.ListenAndServe(":1337", nil); err != nil {
		log.Fatal(err)
	}
}
