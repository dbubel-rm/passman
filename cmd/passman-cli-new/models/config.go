package models

type Config struct {
	Password string `json:"masterPassword"`
	Username string `json:"username"`
	Backend  string `json:"backend"`
}
