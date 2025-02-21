package model

type User struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type DB struct {
	Password    []byte
	IsUserValid bool
}

type SignInResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken,omitempty"`
}
