package config

import (
	"flag"
	"log"
)

const minSecretKeySize = 32

type Config struct {
	SecretKey string
	Port      string
	DSN       string
	STMP      struct {
		Host     string
		Port     int
		Username string
		Password string
		Sender   string
	}
}

func New() *Config {

	var cfg = new(Config)

	flag.StringVar(&cfg.SecretKey, "secret_key", "2e169fe90c625a384c9cad38b3afd3b83ad2176e16fb435923a64666c8410b33", "секретный ключ для подписи")
	if len(cfg.SecretKey) < minSecretKeySize {
		log.Fatalf("Длина ключа должна быть хотя бы %d символа", minSecretKeySize)
	}

	flag.StringVar(&cfg.DSN, "db-dsn", "postgres://postgres:123@localhost:5432/auth_service?sslmode=disable", "PostgreSQL DSN")

	flag.StringVar(&cfg.Port, "port", ":8000", "server port")

	flag.StringVar(&cfg.STMP.Host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.STMP.Port, "smtp-port", 2525, "SMTP port")
	flag.StringVar(&cfg.STMP.Username, "smtp-username", "4e6131523774d9", "SMTP username")
	flag.StringVar(&cfg.STMP.Password, "smtp-password", "e656e78aaac168", "STMP password")
	flag.StringVar(&cfg.STMP.Sender, "smtp-sender", "Todo <no-reply@todo.goserv.ru>", "SMTP sender")

	return cfg
}
