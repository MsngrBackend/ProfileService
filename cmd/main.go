package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"

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

	profileRepo := repository.NewProfilePostgres(db)
	contactsRepo := repository.NewContactsPostgres(db)
	privacyRepo := repository.NewPrivacyPostgres(db)
	// favRepo := repository.NewFavoritePostgres(db)
	// notifRepo := repository.NewNotificationPostgres(db)
	avatarStore := repository.NewMinIOStorage(
		os.Getenv("MINIO_ENDPOINT"),
		os.Getenv("MINIO_ACCESS_KEY"),
		os.Getenv("MINIO_SECRET_KEY"),
	)

	profileUC := usecase.NewProfileUsecase(profileRepo, avatarStore)
	contactsUC := usecase.NewContactsUsecase(contactsRepo)
	privacyUC := usecase.NewPrivacyUsecase(privacyRepo)
	// favUC := usecase.NewFavoriteUsecase(favRepo)
	// notifUC := usecase.NewNotificationUsecase(notifRepo)

	// h := delivery.NewHandler(profileUC, contactUC, privacyUC, favUC, notifUC)
	h := delivery.NewHandler(profileUC, contactsUC, privacyUC)
	router := h.NewRouter()

	log.Println("Profile service running on :8082")
	log.Fatal(http.ListenAndServe(":8082", router))
}
