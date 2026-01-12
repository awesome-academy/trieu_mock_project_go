package dtos

type LoginRequest struct {
	User struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	} `json:"user" binding:"required"`
}

type LoginResponse struct {
	User struct {
		ID          uint   `json:"id"`
		Name        string `json:"name"`
		Email       string `json:"email"`
		AccessToken string `json:"access_token"`
	} `json:"user"`
}

type WSTicketResponse struct {
	Ticket string `json:"ticket"`
}
