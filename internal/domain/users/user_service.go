package users

type UserAuth interface {
	Login(request LoginRequest) (string, User, error)
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
	Create(user *User) error
}

type userService struct {
	authService    UserAuth
	userRepository UserRepository
}

func NewUserService(userRepository UserRepository, authService UserAuth) UserService {
	return userService{
		userRepository: userRepository,
		authService:    authService,
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

func (s userService) Create(user *User) error {
	// TODO this will need to map a user to a new user request
	return s.userRepository.Create(user)
}

func (s userService) Login(usernameOrEmail, password string) (string, User, error) {
	request := LoginRequest{
		UsernameOrEmail: usernameOrEmail,
		Password:        password,
	}

	return s.authService.Login(request)
}
