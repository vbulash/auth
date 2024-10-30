package main

import (
	"context"
	"log"

	"github.com/vbulash/auth/internal/app"
)

func main() {
	ctx := context.Background()

	application, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("Фатальная ошибка инициализации приложения: %s", err.Error())
	}

	err = application.Run()
	if err != nil {
		log.Fatalf("Фатальная ошибка запуска приложения: %s", err.Error())
	}
}
