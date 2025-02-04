package account

// 登入相關
type LoginRequest struct {
	Account    string `json:"account" validate:"required"`
	Password   string `json:"password" validate:"required"`
	IsRemember bool   `json:"isRemember"`
}

type LoginResponse struct {
	Token             string `json:"token"`
	IsSuccess         bool   `json:"isSuccess"`
	LoginErrorMessage string `json:"loginErrorMessage"`
	ThirdPartyError   string `json:"thirdPartyError"`
	ReturnUrl         string `json:"returnUrl"`
}

// 註冊相關
type RegisterRequest struct {
	FirstName       string `json:"firstName" validate:"required"`
	LastName        string `json:"lastName" validate:"required"`
	Account         string `json:"account" validate:"required"`
	Password        string `json:"password" validate:"required"`
	Phone           string `json:"phone"`
	PasswordConfirm string `json:"passwordConfirm"`
}

type RegisterResponse struct {
	Success             bool   `json:"success"`
	AccountErrorMessage string `json:"accountErrorMessage"`
}

// 帳號資訊相關
// 更新姓名請求
type UpdateNameRequest struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
}

// 更新Email請求
type UpdateEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

// 更新電話請求
type UpdatePhoneRequest struct {
	Phone string `json:"phone" validate:"required"`
}

// 更新密碼請求
type UpdatePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" validate:"required"`
	NewPassword     string `json:"newPassword" validate:"required"`
}

type ProfileResponse struct {
	Account               string `json:"account"`
	FirstName             string `json:"firstName"`
	LastName              string `json:"lastName"`
	Email                 string `json:"email"`
	Phone                 string `json:"phone"`
	IsThirdPartyLogin     bool   `json:"isThirdPartyLogin"`
	NameErrorMsg          string `json:"nameErrorMsg"`
	EmailErrorMsg         string `json:"emailErrorMsg"`
	PhoneErrorMsg         string `json:"phoneErrorMsg"`
	PasswordErrorMsg      string `json:"passwordErrorMsg"`
	PasswordUpdateSuccess bool   `json:"passwordUpdateSuccess"`
}

type FacebookUserInfo struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type GoogleUserInfo struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"given_name"`
	LastName  string `json:"family_name"`
}
