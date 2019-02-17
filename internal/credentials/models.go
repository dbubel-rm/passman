package credentials

type Add struct {
	ServiceName string `json:"serviceName" validate:"required"`
	Username    string `json:"username" validate:"required"`
	Password    string `json:"password" validate:"required"`
}
