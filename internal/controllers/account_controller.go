package controllers

import (
	"html/template"
	"log"
	"net/http"
	"regexp"
	"rentjoy/configs"
	"rentjoy/internal/dto/account"
	interfaces "rentjoy/internal/interfaces/services"
	"rentjoy/internal/middleware"
	"rentjoy/internal/services"
)

type AccountController struct {
	BaseController
	accountService interfaces.AccountService
}

var emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

func NewAccountController(accountService interfaces.AccountService, template map[string]*template.Template) *AccountController {
	return &AccountController{
		BaseController: NewBaseController(template),
		accountService: accountService,
	}
}

// 註冊頁面
func (c *AccountController) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	c.RenderTemplate(w, r, "signup", nil)
	return
}

// 註冊作業
func (c *AccountController) SignUpDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		c.RenderTemplate(w, r, "signup_detail", nil)
		return
	} else if r.Method == http.MethodPost {
		// 解析表單
		if err := r.ParseForm(); err != nil {
			http.Error(w, "無法解析表單數據", http.StatusBadRequest)
			return
		}

		// 取得 dto
		request := account.RegisterRequest{
			Account:         r.FormValue("Account"),
			LastName:        r.FormValue("LastName"),
			FirstName:       r.FormValue("FirstName"),
			Phone:           r.FormValue("Phone"),
			Password:        r.FormValue("Password"),
			PasswordConfirm: r.FormValue("PasswordConfirm"),
		}

		regex := regexp.MustCompile(emailRegex)
		if !regex.MatchString(request.Account) {
			vm := map[string]interface{}{
				"AccountErrorMessage": "Email格式不正確",
			}
			c.RenderTemplate(w, r, "signup_detail", vm)
			return
		}

		// 輸入判斷
		if request.Password != request.PasswordConfirm {
			vm := map[string]interface{}{
				"AccountErrorMessage": "密碼不一致",
			}
			c.RenderTemplate(w, r, "signup_detail", vm)
			return
		}

		// 送往 service 處理商務邏輯
		token, err := c.accountService.RegisterAccount(request)
		if err != nil {
			switch err {
			case services.ErrEmailExists:
				vm := map[string]interface{}{
					"AccountErrorMessage": "此信箱已被註冊",
				}
				c.RenderTemplate(w, r, "signup_detail", vm)
			default:
				http.Error(w, "註冊失敗", http.StatusInternalServerError)
			}
			return
		}

		// 處理登入流程
		// 新增 cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			Path:     "/",
			MaxAge:   1 * 60 * 60,
			HttpOnly: true,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// 登入頁面
func (c *AccountController) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// 獲取 returnUrl
		returnUrl := r.URL.Query().Get("returnUrl")

		vm := account.LoginResponse{
			ReturnUrl: returnUrl,
		}

		c.RenderTemplate(w, r, "login", vm)
		return
	} else if r.Method == http.MethodPost {
		// 解析表單
		if err := r.ParseForm(); err != nil {
			http.Error(w, "無法解析表單數據", http.StatusBadRequest)
			return
		}

		// 取得 dto
		request := account.LoginRequest{
			Account:    r.FormValue("Account"),
			Password:   r.FormValue("Password"),
			IsRemember: r.FormValue("IsRemember") == "on",
		}

		// 檢查帳號格式
		_, err := regexp.MatchString(emailRegex, request.Account)
		if err != nil {
			vm := map[string]interface{}{
				"LoginErrorMessage": "Email格式錯誤",
			}
			c.RenderTemplate(w, r, "login", vm)
			return
		}

		// 送往 service 處理商務邏輯
		token, err := c.accountService.Login(request)
		if err != nil {
			vm := map[string]interface{}{
				"LoginErrorMessage": err,
			}
			c.RenderTemplate(w, r, "login", vm)
			return
		}

		maxAge := 3600
		if request.IsRemember {
			maxAge = 24 * 3600
		}
		// 新增 cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			Path:     "/",
			MaxAge:   maxAge,
			HttpOnly: true,
		})

		returnUrl := r.FormValue("returnUrl")
		if returnUrl != "" {
			http.Redirect(w, r, "returnUrl", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

		return
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// 導向 facebook 第三方介面
func (c *AccountController) FacebookLogin(w http.ResponseWriter, r *http.Request) {
	// 生成 URL 並重定向
	facebookConfig := configs.GetFacebookOAuthConfig()
	url := facebookConfig.AuthCodeURL("state")

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// facebook 取得資料後進行本網站的登入流程
func (c *AccountController) FacebookCallback(w http.ResponseWriter, r *http.Request) {
	// 獲取 code
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	// 處理 Facebook 登入
	token, err := c.accountService.FacebookLogin(code)
	if err != nil {
		http.Error(w, "Facebook login failed", http.StatusInternalServerError)
		return
	}

	// 設置 Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   24 * 60 * 60,
	})

	// 重定向到首頁或帳戶頁面
	http.Redirect(w, r, "/Account", http.StatusSeeOther)
}

// 導向 Google 第三方介面
func (c *AccountController) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	googleConfig := configs.GetGoogleOAuthConfig()
	url := googleConfig.AuthCodeURL("state")

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// google 取得資料後進行本網站的登入流程
func (c *AccountController) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
	}

	// 處理 Google 登入
	token, err := c.accountService.GoogleLogin(code)
	if err != nil {
		http.Error(w, "Google login failed", http.StatusInternalServerError)
		return
	}

	// 設置 Cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   24 * 60 * 60,
	})

	// 重定向到首頁或帳戶頁面
	http.Redirect(w, r, "/Account", http.StatusSeeOther)
}

// 登出處理
func (c *AccountController) Logout(w http.ResponseWriter, r *http.Request) {
	// 清除 cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// 帳號資訊頁面
func (c *AccountController) Account(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	profile, err := c.accountService.GetProfile(userID)
	if err != nil {
		http.Error(w, "獲取用戶資料失敗", http.StatusInternalServerError)
		return
	}

	c.RenderTemplate(w, r, "account", profile)
}

// 更新姓名
func (c *AccountController) UpdateName(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "無法解析表單數據", http.StatusBadRequest)
		return
	}

	request := account.UpdateNameRequest{
		FirstName: r.FormValue("FirstName"),
		LastName:  r.FormValue("LastName"),
	}

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	response, err := c.accountService.UpdateName(request, userID)
	if err != nil {
		log.Printf("Update Name Error: %s", err)
		c.RenderTemplate(w, r, "account", response)
		return
	}

	http.Redirect(w, r, "/Account", http.StatusSeeOther)
}

// 更新Email
func (c *AccountController) UpdateEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "無法解析表單數據", http.StatusBadRequest)
		return
	}

	request := account.UpdateEmailRequest{
		Email: r.FormValue("Email"),
	}

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	response, err := c.accountService.UpdateEmail(request, userID)
	if err != nil {
		log.Printf("Update Email Error: %s", err)
		c.RenderTemplate(w, r, "account", response)
		return
	}

	http.Redirect(w, r, "/Account", http.StatusSeeOther)
}

// 更新手機
func (c *AccountController) UpdatePhone(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "無法解析表單數據", http.StatusBadRequest)
		return
	}

	request := account.UpdatePhoneRequest{
		Phone: r.FormValue("Phone"),
	}

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	response, err := c.accountService.UpdatePhone(request, userID)
	if err != nil {
		log.Printf("Update Phone Error: %s", err)
		c.RenderTemplate(w, r, "account", response)
		return
	}

	http.Redirect(w, r, "/Account", http.StatusSeeOther)
}

// 更新密碼
func (c *AccountController) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "無法解析表單數據", http.StatusBadRequest)
		return
	}

	request := account.UpdatePasswordRequest{
		CurrentPassword: r.FormValue("CurrentPassword"),
		NewPassword:     r.FormValue("NewPassword"),
	}

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	response, err := c.accountService.UpdatePassword(request, userID)
	if err != nil {
		log.Printf("Update Password Error: %s", err)
		c.RenderTemplate(w, r, "account", response)
		return
	}

	http.Redirect(w, r, "/Account", http.StatusSeeOther)
}
