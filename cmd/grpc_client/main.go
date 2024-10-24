package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/vbulash/auth/config"
	"github.com/vbulash/auth/database/operations"

	"github.com/brianvoe/gofakeit"
	"github.com/golang/protobuf/ptypes/wrappers"
	user "github.com/vbulash/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const address = "localhost:50051"

func closeConnection(conn *grpc.ClientConn) {
	err := conn.Close()
	if err != nil {
		log.Fatalf("Фатальная ошибка закрытия коннекта к серверу: %v", err)
	}
}

func main() {
	// Week 2
	config.Config = config.LoadConfig()

	db, err := operations.InitDb()
	if err != nil {
		log.Fatalf("Фатальная ошибка коннекта к базе данных: %v", err)
	}

	operations.Seed(db)

	users, err := operations.Get(db)
	if err != nil {
		log.Fatalf("Фатальная ошибка получения данных из базы данных: %v", err)
	}
	log.Println("Получены записи из таблицы users")
	for index := range *users {
		log.Printf("%#v\n", (*users)[index])
	}

	// Week 1
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Фатальная ошибка коннекта к серверу: %v", err)
	}
	defer closeConnection(conn)

	client := user.NewAuthV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Create
	fmt.Println("Клиент: создание пользователя")
	password := gofakeit.Password(false, false, false, true, false, 32)
	newRecord := &user.CreateRequest{
		Name:            gofakeit.Name(),
		Email:           gofakeit.Email(),
		Password:        password,
		PasswordConfirm: password,
		Role:            user.Role_USER,
	}
	fmt.Printf("Клиент: создаем нового пользователя: %+v\n", newRecord)
	response1, err := client.Create(ctx, newRecord)
	if err != nil {
		log.Fatalf("Клиент: фатальная ошибка создания записи пользователя: %v", err)
	}
	fmt.Printf("Клиент: создан новый пользователь ID = %d\n", response1.Id)
	id := response1.Id // Сквозной ID по всем эндпойнтам

	// Get
	fmt.Println()
	fmt.Println("Клиент: получение пользователя")
	fmt.Printf("Клиент: получаем информацию пользователя ID = %d\n", id)
	response2, err := client.Get(ctx, &user.GetRequest{Id: id})
	if err != nil {
		log.Fatalf("Клиент: фатальная ошибка получения записи пользователя ID = %d: %v", id, err)
	}
	fmt.Printf("Клиент: получен пользователь %+v\n", response2)

	// Update
	fmt.Println()
	fmt.Println("Клиент: обновление пользователя")
	name := wrappers.StringValue{Value: gofakeit.Name()}
	email := wrappers.StringValue{Value: gofakeit.Email()}
	record := &user.UpdateRequest{
		Id:    id,
		Name:  &name,
		Email: &email,
		Role:  user.Role_USER,
	}
	// Отображение записи из-за ссылок / StringValue очень некрасивое, но в целом понятное
	fmt.Printf("Клиент: обновляем информацию пользователя ID = %d: %+v\n", id, record)
	_, err = client.Update(ctx, record)
	if err != nil {
		log.Fatalf("Клиент: фатальная ошибка обновления записи пользователя ID = %d: %v", id, err)
	}
	fmt.Printf("Клиент: обновлена запись пользователя ID = %d\n", id)

	// Delete
	fmt.Println()
	fmt.Println("Клиент: удаление пользователя")
	fmt.Printf("Клиент: удаляем запись пользователя ID = %d\n", id)
	_, err = client.Delete(ctx, &user.DeleteRequest{Id: id})
	if err != nil {
		log.Fatalf("Клиент: фатальная ошибка удаления записи пользователя ID = %d: %v", id, err)
	}
	fmt.Printf("Клиент: запись пользователя ID = %d удалена\n", id)
}
