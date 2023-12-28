package main

import (
	"fmt"
	"log"
	"net/http"

	"app/configs"
	"app/login"
	"app/user"
)

func main() {
	db, err := configs.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	loginRepository := login.NewLoginRepository(db)
	userRepository := user.NewUserRepository(db)
	loginService := login.NewLoginService(loginRepository, userRepository)
	loginUseCase := login.NewLoginUseCase(loginService, loginRepository, userRepository)
	loginController := login.NewLoginController(loginUseCase)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		loginController.Login(w, r)
	})

	log.Fatal(http.ListenAndServe(":8000", nil))
}
