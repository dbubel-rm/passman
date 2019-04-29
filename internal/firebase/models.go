package firebase

type Signin struct {
	Email             string `json:"email" validate:"required"`
	Password          string `json:"password" validate:"required"`
	ReturnSecureToken bool   `json:"returnSecureToken" validate:"required"`
}

type Create struct {
	Email             string `json:"email" validate:"required"`
	Password          string `json:"password" validate:"required"`
	ReturnSecureToken bool   `json:"returnSecureToken" validate:"required"`
}

type Delete struct {
	IdToken string `json:"idToken" validate:"required"`
}

type UpdatePassword struct {
	IdToken           string `json:"idToken" validate:"required"`
	Password          string `json:"password" validate:"required"`
	returnSecureToken bool   `json:"returnSecureToken" validate:"required"`
}

type Verify struct {
	RequestType string `json:"requestType" validate:"required"`
	IdToken     string `json:"idToken" validate:"required"`
}
