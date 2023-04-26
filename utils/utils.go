package utils

import (
	"github.com/dlclark/regexp2"
	"math/rand"
	"time"
)

var (
	datePattern = regexp2.MustCompile(`^(?!.*\bещ[ёе]\b)[а-яА-Я]{2,}( [а-яА-Я]{2,}){1,2}$`, 0)
)

func RandAnswer(answer []string) string {
	rand.Seed(time.Now().UnixNano())
	return answer[rand.Intn(len(answer))]
}

func DivideString(s string) []string {
	var chunks []string
	runes := []rune(s)
	for i := 0; i < len(runes); i += 1024 {
		end := i + 1024
		if end > len(runes) {
			end = len(runes)
		}
		if end == len(runes) {
			chunks = append(chunks, string(runes[i:end])+"\nХотите узнать ещё? Назовите другую дату...")
		} else {
			chunks = append(chunks, string(runes[i:end])+"...")
		}

	}
	return chunks
}

func CheckCommand(command string) bool {
	match, _ := datePattern.MatchString(command)
	return match
}
