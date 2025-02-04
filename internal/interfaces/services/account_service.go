package interfaces

import (
	"rentjoy/internal/dto/account"
)

type AccountService interface {
	RegisterAccount(r account.RegisterRequest) (string, error)
	Login(r account.LoginRequest) (string, error)
	GetProfile(userID uint) (*account.ProfileResponse, error)
	UpdateName(r account.UpdateNameRequest, userID uint) (*account.ProfileResponse, error)
	UpdateEmail(r account.UpdateEmailRequest, userID uint) (*account.ProfileResponse, error)
	UpdatePhone(r account.UpdatePhoneRequest, userID uint) (*account.ProfileResponse, error)
	UpdatePassword(r account.UpdatePasswordRequest, userID uint) (*account.ProfileResponse, error)
	FacebookLogin(code string) (string, error)
	GoogleLogin(code string) (string, error)
}
