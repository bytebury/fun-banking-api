package models

type User struct {
	AuditModel
	Username  string `json:"username" gorm:"unique;not null;size:15"`
	Email     string `json:"email" gorm:"unique;not null"`
	FirstName string `json:"first_name" gorm:"not null;size:15"`
	LastName  string `json:"last_name" gorm:"not null;size:20"`
	Password  string `json:"-"`
}

type UserRequest struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
