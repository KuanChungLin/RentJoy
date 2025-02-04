package controllers

import (
	"html/template"
	"log"
	"net/http"
	interfaces "rentjoy/internal/interfaces/services"
	"rentjoy/pkg/helper"
)

type VenuePageController struct {
	BaseController
	venueService interfaces.VenuePageService
}

func NewVenuePageController(venueService interfaces.VenuePageService, templates map[string]*template.Template) *VenuePageController {
	return &VenuePageController{
		BaseController: NewBaseController(templates),
		venueService:   venueService,
	}
}

func (c *VenuePageController) VenuePage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "無法解析表單數據", http.StatusBadRequest)
		return
	}

	venueID, err := helper.StrToInt(r.FormValue("VenueId"))
	if err != nil {
		log.Printf("無法解析 Error: %s", err)
		// TODO
		// http.Redirect(w, r, "/error", http.StatusSeeOther)
		return
	}

	// 刪除 Cookie
	cookie := &http.Cookie{
		Name:   "TimeDetailCookie",
		Path:   "/Venue",
		MaxAge: -1, // 設為-1表示立即刪除
	}
	http.SetCookie(w, cookie)

	vm := c.venueService.GetVenuePage(venueID)

	c.RenderTemplate(w, r, "venuepage", vm)
}
