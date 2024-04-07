package users

import (
	"funbanking/internal/infrastructure/pagination"
	"funbanking/internal/infrastructure/persistence"
	"strings"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetCurrentUser(user *User) error
	FindByID(userID string, user *User) error
	FindByUsernameOrEmail(usernameOrEmail string, user *User) error
	FindAll(itemsPerPage, pageNumber int, params map[string]string) (pagination.PaginatedResponse[User], error)
	Update(userID string, user *User) error
	Create(user *User) error
	AddVisitor(visitor *Visitor) error
	UpdateEmail(user *User) error
	Verify(user *User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository() UserRepository {
	return userRepository{db: persistence.DB}
}

func (r userRepository) GetCurrentUser(user *User) error {
	return r.db.First(&user, "username = ?", "marcello").Error
}

func (r userRepository) FindByID(userID string, user *User) error {
	return r.db.First(&user, "id = ?", userID).Error
}

func (r userRepository) FindByUsernameOrEmail(usernameOrEmail string, user *User) error {
	usernameOrEmail = strings.TrimSpace(strings.ToLower(usernameOrEmail))
	return r.db.First(&user, "username = ? or email = ?", usernameOrEmail, usernameOrEmail).Error
}

func (r userRepository) FindAll(itemsPerPage, pageNumber int, params map[string]string) (pagination.PaginatedResponse[User], error) {
	query := r.db.Find(&User{}).Order("created_at DESC")

	if params["ID"] != "" {
		query = query.Where("id = ?", params["ID"])
	}

	if params["Username"] != "" {
		query = query.Where("username ILIKE ?", "%"+params["Username"]+"%")
	}

	if params["Email"] != "" {
		query = query.Where("email ILIKE ?", "%"+params["Email"]+"%")
	}

	if params["LastSeen"] != "" {
		query = query.Where("last_seen > ?", params["LastSeen"])
	}

	return pagination.Find[User](query, pageNumber, itemsPerPage)
}

func (r userRepository) Update(userID string, user *User) error {
	var foundUser User
	if err := r.FindByID(userID, &foundUser); err != nil {
		return err
	}

	if user.Username == "" {
		user.Username = foundUser.Username
	}

	if user.FirstName == "" {
		user.FirstName = foundUser.FirstName
	}

	if user.LastName == "" {
		user.LastName = foundUser.LastName
	}

	if user.Avatar == "" {
		user.Avatar = foundUser.Avatar
	}

	if user.About == "" {
		user.About = foundUser.About
	}

	if user.LastSeen.IsZero() {
		user.LastName = foundUser.LastName
	}

	r.normalize(user)

	return r.db.Model(&foundUser).Select("Username", "FirstName", "LastName", "Avatar", "About", "LastSeen").Updates(&user).Error
}

func (r userRepository) UpdateEmail(user *User) error {
	user.Verified = false
	return r.db.Save(user).Error
}

func (r userRepository) Verify(user *User) error {
	return r.db.Save(user).Error
}

func (r userRepository) Create(user *User) error {
	r.normalize(user)

	passwordHash, err := HashString(user.Password)

	if err != nil {
		return err
	}

	user.Password = passwordHash

	return r.db.Create(&user).Error
}

func (r userRepository) AddVisitor(visitor *Visitor) error {
	return r.db.Create(&visitor).Error
}

func (r userRepository) normalize(user *User) error {
	user.Username = strings.ToLower(user.Username)
	user.Email = strings.ToLower(user.Email)
	user.FirstName = strings.ToLower(user.FirstName)
	user.LastName = strings.ToLower(user.LastName)

	return nil
}
