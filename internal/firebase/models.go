package firebase

type Signin struct {
	Email             string `json:"email" validate:"required"`
	Password          string `json:"password" validate:"required"`
	ReturnSecureToken bool   `json:"returnSecureToken" validate:"required"`
}

type Create struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Delete struct {
	IdToken string `json:"idToken" validate:"required"`
}
