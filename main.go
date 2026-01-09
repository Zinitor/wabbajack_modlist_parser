package main

import "wabbajackModlistParser/parser"

func main() {
	// план - ускорить все это чтобы оно занимало меньше 1й секунды в идеале и меньше 3х в допустимом виде
	// сейчас - 5 секунд
	// необходимо уменьшить потребление памяти потому что сейчас оно сьедает 1.5 гига на чтение всех модов
	// поменять везде где используется io.ReadAll на json.NewDecoder(response.Body).Decode(&data)
	// функции приемщики долны читать из io.Reader
	includeGames := []string{"skyrimspecialedition"}
	parser.MainParse(includeGames)
}
