package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"app/configs"
	"app/login"
	reviewController "app/review/controller"
	reviewRepository "app/review/repository"
	reviewUseCase "app/review/usecase"
	"app/user"
)

// ReviewControllerの初期化
func initReviewController(db *sql.DB) reviewController.ControllerInterface {
	reviewRepository :=reviewRepository.NewReviewRepository(db)
	reviewUseCase := reviewUseCase.NewReviewUseCase(reviewRepository)
	return reviewController.NewReviewController(reviewUseCase)
}

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

	reviewController := initReviewController(db)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World")
	})
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		loginController.Login(w, r)
	})

	http.HandleFunc("/api/v1/reviews", func(w http.ResponseWriter, r *http.Request) {
		reviewController.GetReviews(w, r)
	})

	log.Fatal(http.ListenAndServe(":8000", nil))
}
