package users

type UserAuth interface {
	Login(request LoginRequest) (string, User, error)
}

type WelcomeMailer interface {
	SendEmail(recipient string, user User) error
}

type LoginRequest struct {
	UsernameOrEmail string `json:"username_or_email"`
	Password        string `json:"password"`
}

type UserService interface {
	FindByID(id string) (User, error)
	FindByUsernameOrEmail(usernameOrEmail string) (User, error)
	Update(id string, user *User) error
	Login(usernameOrEmail, password string) (string, User, error)
	Create(request *NewUserRequest) (User, error)
}

type userService struct {
	authService    UserAuth
	userRepository UserRepository
	welcomeMailer  WelcomeMailer
}

func NewUserService(userRepository UserRepository, authService UserAuth, welcomeMailer WelcomeMailer) UserService {
	return userService{
		userRepository: userRepository,
		authService:    authService,
		welcomeMailer:  welcomeMailer,
	}
}

func (s userService) FindByID(id string) (User, error) {
	var user User
	err := s.userRepository.FindByID(id, &user)
	return user, err
}

func (s userService) FindByUsernameOrEmail(usernameOrEmail string) (User, error) {
	var user User
	err := s.userRepository.FindByUsernameOrEmail(usernameOrEmail, &user)
	return user, err
}

func (s userService) Update(id string, user *User) error {
	return s.userRepository.Update(id, user)
}

func (s userService) Create(request *NewUserRequest) (User, error) {
	user := User{
		Username:  request.Username,
		Email:     request.Email,
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Password:  request.Password,
		Role:      0,
		About:     "",
		Avatar:    "https://www.gravatar.com/avatar/2533c61da0bd2b79b63fd599cd045a31?default=https%3A%2F%2Fcloud.digitalocean.com%2Favatars%2Fdefault30.png&secure=true",
	}

	if err := s.userRepository.Create(&user); err != nil {
		return User{}, err
	}

	s.welcomeMailer.SendEmail(user.Email, user)

	return user, nil
}

func (s userService) Login(usernameOrEmail, password string) (string, User, error) {
	request := LoginRequest{
		UsernameOrEmail: usernameOrEmail,
		Password:        password,
	}

	return s.authService.Login(request)
}
