package controllers

import (
	"html/template"
	"log"
	"net/http"
	"reflect"
	"rentjoy/internal/dto/venuepage"
	"rentjoy/pkg/auth"
)

type BaseController struct {
	Templates map[string]*template.Template
}

type BaseViewModel struct {
	IsOrderPending bool
	IsLoggedIn     bool
	UserEmail      string
	PageData       interface{}
}

func NewBaseController(templates map[string]*template.Template) BaseController {
	return BaseController{
		Templates: templates,
	}
}

func (c *BaseController) RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, data interface{}) {
	template, ok := c.Templates[tmpl]
	if !ok {
		log.Printf("Template '%s' not found", tmpl)
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	// 創建基礎 layout 結構
	vm := BaseViewModel{
		IsOrderPending: false,
		IsLoggedIn:     false,
		UserEmail:      "",
		PageData:       data,
	}

	// 從 Cookie 讀取 token 判斷登入狀態
	cookie, err := r.Cookie("token")
	if err == nil {
		// 驗證 token
		if claims, err := auth.ValidateToken(cookie.Value); err == nil {
			vm.IsLoggedIn = true
			vm.UserEmail = claims.Email
		}
	}

	reflectType := reflect.TypeOf(venuepage.OrderPending{})
	pageDataType := reflect.TypeOf(vm.PageData)

	// 檢查是否為指標類型
	if pageDataType != nil && pageDataType.Kind() == reflect.Ptr {
		// 獲取指標指向的類型
		pageDataType = pageDataType.Elem()
		log.Println("pageDataType (after Elem()):", pageDataType)
	}

	// 依照不同的 PageData 顯示 layout
	if pageDataType == reflectType {
		vm.IsOrderPending = true
	}

	err = template.ExecuteTemplate(w, "layout", vm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Template execution error: %v", err)
	}
}
