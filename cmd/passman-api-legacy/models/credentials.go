package models

type Credentials struct {
	ServiceName string `json:"serviceName" db:"service_name"`
	Username    string `json:"username" db:"username"`
	Password    string `json:"password" db:"password"`
}
