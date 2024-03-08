package auth

type Request struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `password:"password"`
}
