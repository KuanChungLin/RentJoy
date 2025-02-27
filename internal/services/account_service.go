package services

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"regexp"
	"rentjoy/configs"
	"rentjoy/internal/dto/account"
	repoInterfaces "rentjoy/internal/interfaces/repositories"
	serviceInterfaces "rentjoy/internal/interfaces/services"
	"rentjoy/internal/models"
	"rentjoy/internal/repositories"
	"rentjoy/pkg/auth"
	"rentjoy/pkg/helper"
	"time"

	"gorm.io/gorm"
)

var (
	ErrEmailExists    = errors.New("email already exists")
	isThirdPartyLogin = false
)

type AccountService struct {
	memberRepo   repoInterfaces.MemberRepository
	facebookRepo repoInterfaces.FacebookRepository
	googleRepo   repoInterfaces.GoogleRepository
}

func NewAccountService(db *gorm.DB) serviceInterfaces.AccountService {
	return &AccountService{
		memberRepo:   repositories.NewMemberRepository(db),
		facebookRepo: repositories.NewFacebookRepository(db),
		googleRepo:   repositories.NewGoogleRepository(db),
	}
}

// 註冊帳號
func (s *AccountService) RegisterAccount(r account.RegisterRequest) (token string, err error) {
	member, err := s.memberRepo.FindByAccount(r.Account)
	if err != nil {
		log.Printf("Find By Account Error:", err)
		return "", err
	}

	if member != nil {
		return "", ErrEmailExists
	}

	newMember := models.Member{
		Account:   r.Account,
		Password:  helper.GetSHA256(r.Password),
		FirstName: r.FirstName,
		LastName:  r.LastName,
		Email:     r.Account,
		Phone:     r.Phone,
		CreatedAt: time.Now(),
		UpdatedAt: nil,
	}

	err = s.memberRepo.Create(newMember)
	if err != nil {
		log.Printf("Create Member Error:", err)
		return "", err
	}

	member, err = s.memberRepo.FindByAccount(r.Account)
	if err != nil {
		log.Printf("Find By Account Error:", err)
		return "", err
	}

	token, err = auth.GenerateToken(member.ID, member.Email, false)
	if err != nil {
		return "", err
	}

	return token, nil
}

// 登入帳號
func (s *AccountService) Login(r account.LoginRequest) (code string, err error) {
	member, err := s.memberRepo.FindByAccount(r.Account)
	if err != nil {
		log.Printf("Find By Account Error: %s", err)
		return "", err
	}

	if member == nil {
		return "", errors.New("帳號或密碼輸入錯誤")
	}

	if !helper.CheckPasswordHash(r.Password, member.Password) {
		return "", errors.New("帳號或密碼輸入錯誤")
	}

	token, err := auth.GenerateToken(member.ID, member.Email, r.IsRemember)
	if err != nil {
		return "", err
	}

	return token, nil
}

// FaceBook 登入
func (s *AccountService) FacebookLogin(code string) (jwtToken string, err error) {
	// 使用 code 換取 token
	token, err := configs.GetFacebookOAuthConfig().Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Facebook Code Exchange To Token Error: %s", err)
		return "", err
	}

	// 獲取用戶資訊
	client := configs.GetFacebookOAuthConfig().Client(context.Background(), token)
	response, err := client.Get("https://graph.facebook.com/me?fields=id,email,first_name,last_name")
	if err != nil {
		log.Printf("Facebook Client Get Response Error: %s", err)
		return "", err
	}
	defer response.Body.Close()

	var fbUser account.FacebookUserInfo
	if err := json.NewDecoder(response.Body).Decode(&fbUser); err != nil {
		log.Printf("FaceBook JSON Decode Error: %s", err)
		return "", err
	}

	member, err := s.memberRepo.FindByFacebookID(fbUser.ID)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Printf("Find By Account Error: %s", err)
		return "", nil
	}

	if member == nil {
		// 創建用戶
		member = &models.Member{
			Account:   "您以FaceBook帳號登入RentJoy",
			Password:  helper.GetSHA256(helper.GenerateRandomPassword()),
			Email:     fbUser.Email,
			FirstName: fbUser.FirstName,
			LastName:  fbUser.LastName,
			CreatedAt: time.Now(),
			UpdatedAt: nil,
			FacebookLogins: []models.FacebookThirdPartyLogin{
				{
					FacebookThirdPartyID: fbUser.ID,
				},
			},
		}

		if err := s.memberRepo.Create(*member); err != nil {
			log.Printf("Create Facebook Member Error: %s", err)
			return "", err
		}

		member, err = s.memberRepo.FindByFacebookID(fbUser.ID)
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Printf("Find By Account Error: %s", err)
			return "", nil
		}

	} else {
		member.FirstName = fbUser.FirstName
		member.LastName = fbUser.LastName
		member.Email = fbUser.Email

		if err := s.memberRepo.Update(*member); err != nil {
			return "", err
		}
	}
	jwtToken, err = auth.GenerateToken(member.ID, member.Email, false)
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

// Google 登入
func (s *AccountService) GoogleLogin(code string) (jwtToken string, err error) {
	// 使用 code 換取 token
	token, err := configs.GetGoogleOAuthConfig().Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Facebook Code Exchange To Token Error: %s", err)
		return "", err
	}

	// 獲取用戶資訊
	client := configs.GetFacebookOAuthConfig().Client(context.Background(), token)
	response, err := client.Get("https://www.googleapis.com/oauth2/v1/userinfo?alt=json")
	if err != nil {
		log.Printf("Facebook Client Get Response Error: %s", err)
		return "", err
	}
	defer response.Body.Close()

	var gUser account.GoogleUserInfo
	if err := json.NewDecoder(response.Body).Decode(&gUser); err != nil {
		log.Printf("FaceBook JSON Decode Error: %s", err)
		return "", err
	}

	member, err := s.memberRepo.FindByGoogleID(gUser.ID)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Printf("Find By Account Error: %s", err)
		return "", nil
	}

	if member == nil {
		// 創建用戶
		member = &models.Member{
			Account:   "您以Google帳號登入RentJoy",
			Password:  helper.GetSHA256(helper.GenerateRandomPassword()),
			Email:     gUser.Email,
			FirstName: gUser.FirstName,
			LastName:  gUser.LastName,
			CreatedAt: time.Now(),
			UpdatedAt: nil,
			GoogleLogins: []models.GoogleThirdPartyLogin{
				{
					GoogleThirdPartyID: gUser.ID,
				},
			},
		}

		if err := s.memberRepo.Create(*member); err != nil {
			log.Printf("Create Google Member Error: %s", err)
			return "", err
		}

		member, err = s.memberRepo.FindByGoogleID(gUser.ID)
		if err != nil && err != gorm.ErrRecordNotFound {
			log.Printf("Find By Account Error: %s", err)
			return "", nil
		}

	} else {
		member.FirstName = gUser.FirstName
		member.LastName = gUser.LastName
		member.Email = gUser.Email

		if err := s.memberRepo.Update(*member); err != nil {
			return "", err
		}
	}
	jwtToken, err = auth.GenerateToken(member.ID, member.Email, false)
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

// 取得帳號資訊
func (s *AccountService) GetProfile(userID uint) (*account.ProfileResponse, error) {
	member, err := s.memberRepo.FindByID(userID)
	if err != nil {
		log.Printf("Find By ID Error: %s", err)
		return nil, err
	}

	if len(member.FacebookLogins) > 0 || len(member.GoogleLogins) > 0 {
		isThirdPartyLogin = true
	}

	r := account.ProfileResponse{
		Account:           member.Account,
		FirstName:         member.FirstName,
		LastName:          member.LastName,
		Email:             member.Email,
		Phone:             member.Phone,
		IsThirdPartyLogin: isThirdPartyLogin,
	}

	return &r, nil
}

// 更新姓名
func (s *AccountService) UpdateName(r account.UpdateNameRequest, userID uint) (*account.ProfileResponse, error) {
	member, err := s.memberRepo.FindByID(userID)
	if err != nil {
		log.Printf("Find By ID Error: %s", err)
		return nil, err
	}

	if len(member.FacebookLogins) > 0 || len(member.GoogleLogins) > 0 {
		isThirdPartyLogin = true
	}

	// 檢查姓名是否都有填寫
	if r.FirstName == "" || r.LastName == "" {
		return &account.ProfileResponse{
			Account:           member.Account,
			FirstName:         member.FirstName, // 保持原始資料
			LastName:          member.LastName,  // 保持原始資料
			Email:             member.Email,
			Phone:             member.Phone,
			IsThirdPartyLogin: isThirdPartyLogin,
			NameErrorMsg:      "姓名欄位不可為空",
		}, errors.New("name fields cannot be empty")
	}

	// 暫存原始資料（用於更新失敗時返回）
	originalFirstName := member.FirstName
	originalLastName := member.LastName

	// 嘗試更新資料
	member.FirstName = r.FirstName
	member.LastName = r.LastName
	member.UpdatedAt = helper.TimePtr(time.Now())

	err = s.memberRepo.Update(*member)

	if err != nil {
		return &account.ProfileResponse{
			Account:           member.Account,
			FirstName:         originalFirstName, // 返回原始資料
			LastName:          originalLastName,  // 返回原始資料
			Email:             member.Email,
			Phone:             member.Phone,
			IsThirdPartyLogin: isThirdPartyLogin,
			NameErrorMsg:      "姓名更新失敗",
		}, err
	}

	// 更新成功才返回新資料
	return &account.ProfileResponse{
		Account:           member.Account,
		FirstName:         member.FirstName,
		LastName:          member.LastName,
		Email:             member.Email,
		Phone:             member.Phone,
		IsThirdPartyLogin: isThirdPartyLogin,
	}, nil
}

// 更新 Email
func (s *AccountService) UpdateEmail(r account.UpdateEmailRequest, userID uint) (*account.ProfileResponse, error) {
	// 取得 Token 的使用者資料
	member, err := s.memberRepo.FindByID(userID)
	if err != nil {
		log.Printf("Find By ID Error: %s", err)
		return nil, err
	}

	// 判斷帳號是否為第三方登入
	if len(member.FacebookLogins) > 0 || len(member.GoogleLogins) > 0 {
		isThirdPartyLogin = true
	}

	// 保存原始資料
	originalEmail := member.Email

	// 驗證 Email 格式
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(r.Email) {
		return &account.ProfileResponse{
			Account:           member.Account,
			FirstName:         member.FirstName,
			LastName:          member.LastName,
			Email:             originalEmail,
			Phone:             member.Phone,
			IsThirdPartyLogin: isThirdPartyLogin,
			EmailErrorMsg:     "Email 格式不正確",
		}, errors.New("invalid email format")
	}

	member.Email = r.Email
	member.UpdatedAt = helper.TimePtr(time.Now())

	// 將修改後資料更新至資料庫內
	err = s.memberRepo.Update(*member)
	if err != nil {
		return &account.ProfileResponse{
			Account:           member.Account,
			FirstName:         member.FirstName,
			LastName:          member.LastName,
			Email:             originalEmail,
			Phone:             member.Phone,
			IsThirdPartyLogin: isThirdPartyLogin,
			EmailErrorMsg:     "Email 更新失敗",
		}, err
	}

	return &account.ProfileResponse{
		Account:           member.Account,
		FirstName:         member.FirstName,
		LastName:          member.LastName,
		Email:             member.Email,
		Phone:             member.Phone,
		IsThirdPartyLogin: isThirdPartyLogin,
	}, nil
}

// 更新電話
func (s *AccountService) UpdatePhone(r account.UpdatePhoneRequest, userID uint) (*account.ProfileResponse, error) {
	// 取得 Token 的使用者資料
	member, err := s.memberRepo.FindByID(userID)
	if err != nil {
		log.Printf("Find By ID Error: %s", err)
		return nil, err
	}

	// 判斷帳號是否為第三方登入
	if len(member.FacebookLogins) > 0 || len(member.GoogleLogins) > 0 {
		isThirdPartyLogin = true
	}

	originalPhone := member.Phone

	// 驗證手機格式
	if !helper.ValidatePhoneNumber(r.Phone) {
		return &account.ProfileResponse{
			Account:           member.Account,
			FirstName:         member.FirstName,
			LastName:          member.LastName,
			Email:             member.Email,
			Phone:             originalPhone,
			IsThirdPartyLogin: isThirdPartyLogin,
			PhoneErrorMsg:     "手機號碼格式不正確",
		}, errors.New("invalid phone format")
	}

	member.Phone = r.Phone
	member.UpdatedAt = helper.TimePtr(time.Now())

	// 將修改後資料更新至資料庫內
	err = s.memberRepo.Update(*member)
	if err != nil {
		return &account.ProfileResponse{
			Account:           member.Account,
			FirstName:         member.FirstName,
			LastName:          member.LastName,
			Email:             member.Email,
			Phone:             originalPhone,
			IsThirdPartyLogin: isThirdPartyLogin,
			PhoneErrorMsg:     "手機號碼更新失敗",
		}, err
	}

	return &account.ProfileResponse{
		Account:           member.Account,
		FirstName:         member.FirstName,
		LastName:          member.LastName,
		Email:             member.Email,
		Phone:             r.Phone,
		IsThirdPartyLogin: isThirdPartyLogin,
	}, nil
}

// 更新密碼
func (s *AccountService) UpdatePassword(r account.UpdatePasswordRequest, userID uint) (*account.ProfileResponse, error) {
	// 取得 Token 的使用者資料
	member, err := s.memberRepo.FindByID(userID)
	if err != nil {
		log.Printf("Find By ID Error: %s", err)
		return nil, err
	}

	// 判斷帳號是否為第三方登入
	if len(member.FacebookLogins) > 0 || len(member.GoogleLogins) > 0 {
		isThirdPartyLogin = true
	}

	// 驗證當前密碼
	if !helper.CheckPasswordHash(r.CurrentPassword, member.Password) {
		return &account.ProfileResponse{
			Account:           member.Account,
			FirstName:         member.FirstName,
			LastName:          member.LastName,
			Email:             member.Email,
			Phone:             member.Phone,
			IsThirdPartyLogin: isThirdPartyLogin,
			PasswordErrorMsg:  "當前密碼不正確",
		}, errors.New("current password incorrect")
	}

	// 檢查新舊密碼是否相同
	if r.CurrentPassword == r.NewPassword {
		return &account.ProfileResponse{
			Account:           member.Account,
			FirstName:         member.FirstName,
			LastName:          member.LastName,
			Email:             member.Email,
			Phone:             member.Phone,
			IsThirdPartyLogin: isThirdPartyLogin,
			PasswordErrorMsg:  "新密碼不能與當前密碼相同",
		}, errors.New("new password same as current")
	}

	// 更新密碼
	member.Password = helper.GetSHA256(r.NewPassword)
	member.UpdatedAt = helper.TimePtr(time.Now())

	err = s.memberRepo.Update(*member)
	response := account.ProfileResponse{
		Account:               member.Account,
		FirstName:             member.FirstName,
		LastName:              member.LastName,
		Email:                 member.Email,
		Phone:                 member.Phone,
		IsThirdPartyLogin:     isThirdPartyLogin,
		PasswordUpdateSuccess: true,
	}

	if err != nil {
		response.PasswordErrorMsg = "密碼更新失敗"
		response.PasswordUpdateSuccess = false
		return &response, err
	}

	return &response, nil
}
