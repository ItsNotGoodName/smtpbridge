package models

type Config struct {
	RetentionPolicy RetentionPolicy
	AuthSMTP        Auth
	AuthHTTP        Auth
}
