package main

import (
	"log/slog"
	"net/http"
	"os"
	"wabbajackModlistParser/api"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// план - ускорить все это чтобы оно занимало меньше 1й секунды в идеале и меньше 3х в допустимом виде
	// сейчас - 5 секунд
	// необходимо уменьшить потребление памяти потому что сейчас оно сьедает 1.5 гига на чтение всех модов
	// поменять везде где используется io.ReadAll на json.NewDecoder(response.Body).Decode(&data)
	// функции приемщики долны читать из io.Reader
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{})
	sLogger := slog.New(handler)

	myHandler := api.NewHandler(sLogger)

	r.Route("/", myHandler.RegisterRoutes)

	http.ListenAndServe(":3000", r)
}
