package main

import (
	db2 "authService/db"
	"authService/internal/config"
	"authService/internal/mail"
	"authService/internal/models"
	"authService/internal/models/tokens"
	"log"
	"net/http"
)

type application struct {
	Users    models.UserModel
	Sessions models.SessionModel
	JWT      *tokens.JWT
	Mailer   mail.Mailer
}

func main() {

	cfg := config.New()

	// /tokens/get/{id} -  выдает пару Access, Refresh токенов для пользователя с идентификатором (GUID) указанным в параметре запроса
	// /tokens/refresh - выполняет Refresh операцию на пару Access, Refresh токенов

	db, err := db2.OpenDB(cfg.DSN)
	if err != nil {
		log.Fatal(err)
	}

	app := &application{
		Users:    models.UserModel{DB: db},
		Sessions: models.SessionModel{DB: db},
		JWT:      tokens.NewJWT(cfg.SecretKey),
		Mailer:   mail.New(cfg.STMP.Host, cfg.STMP.Port, cfg.STMP.Username, cfg.STMP.Password, cfg.STMP.Sender),
	}

	server := http.Server{
		Addr:    cfg.Port,
		Handler: app.routes(),
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
