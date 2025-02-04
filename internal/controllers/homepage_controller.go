package controllers

import (
	"html/template"
	"net/http"
	interfaces "rentjoy/internal/interfaces/services"
)

type HomePageController struct {
	BaseController
	homeService interfaces.HomeService
}

func NewHomePageController(homeService interfaces.HomeService, templates map[string]*template.Template) *HomePageController {
	return &HomePageController{
		BaseController: NewBaseController(templates),
		homeService:    homeService,
	}
}

// 首頁
func (c *HomePageController) HomePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vm := c.homeService.GetHomePage()

	c.RenderTemplate(w, r, "homepage", vm)
}

// 隱私權保護條款頁面
func (c *HomePageController) Privacy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	c.RenderTemplate(w, r, "privacy_policy", nil)
}

// 常見問答頁面
func (c *HomePageController) Faq(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	c.RenderTemplate(w, r, "faq", nil)
}

// 使用者條款頁面
func (c *HomePageController) UserRules(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	c.RenderTemplate(w, r, "user_rules", nil)
}

// 線上預訂條約頁面
func (c *HomePageController) OrderRules(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	c.RenderTemplate(w, r, "order_rules", nil)
}

// 成為場地主頁面
func (c *HomePageController) BecomeVenueOwner(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	c.RenderTemplate(w, r, "become_venue_owner", nil)
}
