package models

type User struct {
	Id       int
	Name     string `validate:"required,min=3"`
	Email    string `validate:"required,email"`
	Username string `validate:"required"`
	Password string `validate:"required,min=6"`
	Phone    string `validate:"required,len=10,numeric"`
}

type Post struct {
	Id      int    `json:"id"`
	User_Id int    `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
type ContactDetails struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
