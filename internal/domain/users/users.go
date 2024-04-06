package users

import (
	"funbanking/internal/domain"
	"funbanking/internal/infrastructure/persistence"
	"time"
)

type User struct {
	domain.AuditModel
	Username         string    `json:"username" gorm:"unique;not null;size:15"`
	Email            string    `json:"email" gorm:"unique;not null"`
	FirstName        string    `json:"first_name" gorm:"not null;size:20"`
	LastName         string    `json:"last_name" gorm:"not null;size:20"`
	Role             int       `json:"role" gorm:"not null; default:0"`
	About            string    `json:"about" gorm:"type:text"`
	Avatar           string    `json:"avatar" gorm:"not null;type:text;default:https://www.gravatar.com/avatar/2533c61da0bd2b79b63fd599cd045a31?default=https%3A%2F%2Fcloud.digitalocean.com%2Favatars%2Fdefault30.png&secure=true"`
	LastSeen         time.Time `json:"last_seen"`
	SubscriptionTier int16     `json:"subscription_tier" gorm:"not null;default:0"`
	Password         string    `json:"-"`
}

type NewUserRequest struct {
	Username             string `json:"username"`
	Email                string `json:"email"`
	FirstName            string `json:"first_name"`
	LastName             string `json:"last_name"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

type UpdateUserRequest struct {
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar"`
	About     string `json:"about"`
}

type Visitor struct {
	domain.AuditModel
	UserID     *uint  `gorm:"user_id"`
	CustomerID *uint  `gorm:"customer_id"`
	IPAddress  string `gotm:"ip_address"`
}

func RunMigrations() {
	persistence.DB.AutoMigrate(&User{})
	persistence.DB.AutoMigrate(&Visitor{})
}
