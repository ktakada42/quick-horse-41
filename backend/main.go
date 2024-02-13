package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"app/configs"
	"app/login/controller"
	"app/login/repository"
	"app/login/service"
	"app/login/usecase"
	ur "app/user/repository"

	bookController "app/book/controller"
	bookRepository "app/book/repository"
	bookUseCase "app/book/usecase"
)

func initBookController(db *sql.DB) bookController.ControllerInterface {
	bookRepository := bookRepository.NewBookRepository(db)
	bookUseCase := bookUseCase.NewBookUseCase(bookRepository)
	bookController := bookController.NewBookController(bookUseCase)
	return bookController
}

func main() {
	db, err := configs.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	loginRepository := repository.NewLoginRepository(db)
	userRepository := ur.NewUserRepository(db)
	loginService := service.NewLoginService(loginRepository, userRepository)
	loginUseCase := usecase.NewLoginUseCase(loginService, loginRepository, userRepository)
	loginController := controller.NewLoginController(loginUseCase)

	bookController := initBookController(db)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		loginController.Login(w, r)
	})

	http.HandleFunc("/books/offices", func(w http.ResponseWriter, r *http.Request) {
		bookController.GetOfficeBooksController(w, r)
	})

	log.Fatal(http.ListenAndServe(":8000", nil))
}
