package models

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type UserResponse struct {
	User         User   `json:"user"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type UserRegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Name     string `json:"name" validate:"required"`
	Role     string `json:"role" validate:"required,oneof=tracked caregiver"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func NewUserRegisterRequest(email, password, name, role string) *UserRegisterRequest {
	return &UserRegisterRequest{
		Email:    email,
		Password: password,
		Name:     name,
		Role:     role,
	}
}

func NewUserLoginRequest(email, password string) *UserLoginRequest {
	return &UserLoginRequest{
		Email:    email,
		Password: password,
	}
}

func NewUser(id, email, name, role, createdAt, updatedAt string) *User {
	return &User{
		ID:        id,
		Email:     email,
		Name:      name,
		Role:      role,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

func NewUserResponse(user *User, token, refreshToken string) *UserResponse {
	return &UserResponse{
		User:         *user,
		Token:        token,
		RefreshToken: refreshToken,
	}
}
