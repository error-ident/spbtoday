package utils

import (
	"github.com/dlclark/regexp2"
	"math/rand"
	"strings"
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
	var sb strings.Builder
	words := strings.Fields(s)
	for i := range words {
		if sb.Len()+len(words[i])+1 > 1024 {
			chunks = append(chunks, sb.String()+" ...")
			sb.Reset()
		}
		if sb.Len() > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(words[i])
	}
	if sb.Len() > 0 {
		chunks = append(chunks, sb.String())
	}
	return chunks
}

func CheckCommand(command string) bool {
	match, _ := datePattern.MatchString(command)
	return match
}
