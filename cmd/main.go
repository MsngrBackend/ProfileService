package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/jmoiron/sqlx"
	_ "github.com/jackc/pgx/v5/stdlib"

	delivery "github.com/MsngrBackend/ProfileService/internal/delivery/http"
	"github.com/MsngrBackend/ProfileService/internal/repository"
	"github.com/MsngrBackend/ProfileService/internal/usecase"
)

func main() {

	if err := godotenv.Load(); err != nil {
    log.Println("no .env file found")
  }

	db := sqlx.MustConnect("pgx", os.Getenv("DATABASE_URL"))
	defer db.Close()

	// Repositories
	profileRepo := repository.NewProfilePostgres(db)
	// contactRepo  := repository.NewContactPostgres(db)
	privacyRepo  := repository.NewPrivacyPostgres(db)
	// favRepo      := repository.NewFavoritePostgres(db)
	// notifRepo    := repository.NewNotificationPostgres(db)
	avatarStore  := repository.NewMinIOStorage(
		os.Getenv("MINIO_ENDPOINT"),
		os.Getenv("MINIO_ACCESS_KEY"),
		os.Getenv("MINIO_SECRET_KEY"),
	)

	// Usecases
	profileUC := usecase.NewProfileUsecase(profileRepo, avatarStore)
	// contactUC  := usecase.NewContactUsecase(contactRepo)
	privacyUC  := usecase.NewPrivacyUsecase(privacyRepo)
	// favUC      := usecase.NewFavoriteUsecase(favRepo)
	// notifUC    := usecase.NewNotificationUsecase(notifRepo)

	// Handler + Router
	h := delivery.NewHandler(profileUC, privacyUC, os.Getenv("JWT_SECRET"))
	router := h.NewRouter()

	log.Println("Profile service running on :8082")
	log.Fatal(http.ListenAndServe(":8082", router))
}
