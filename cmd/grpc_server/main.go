package main

import (
	"context"
	"github.com/vbulash/auth/internal/app"
	"log"
)

func main() {
	ctx := context.Background()

	app, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("Фатальная ошибка инициализации приложения: %s", err.Error())
	}

	err = app.Run()
	if err != nil {
		log.Fatalf("Фатальная ошибка запуска приложения: %s", err.Error())
	}
}
