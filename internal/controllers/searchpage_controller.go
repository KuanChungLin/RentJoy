package controllers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"rentjoy/internal/dto/searchpage"
	interfaces "rentjoy/internal/interfaces/services"
	"rentjoy/pkg/helper"
	"strconv"
)

type SearchPageController struct {
	BaseController
	searchService interfaces.SearchPageService
}

func NewSearchPageController(searchService interfaces.SearchPageService, templates map[string]*template.Template) *SearchPageController {
	return &SearchPageController{
		BaseController: NewBaseController(templates),
		searchService:  searchService,
	}
}

// 搜尋頁
func (c *SearchPageController) SearchPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 解析表單
	if err := r.ParseForm(); err != nil {
		http.Error(w, "無法解析表單數據", http.StatusBadRequest)
		return
	}

	activityID, err := helper.StrToUint(r.FormValue("ActivityId"))
	if r.FormValue("ActivityId") != "" && err != nil {
		log.Printf("無法解析 ActivityId Error: %s", err)
		return
	}

	maxPrice, err := strconv.Atoi(r.FormValue("MaxPrice"))
	if r.FormValue("MaxPrice") != "" && err != nil {
		log.Printf("無法解析 maxPrice Error: %s", err)
		return
	}

	minPrice, err := strconv.Atoi(r.FormValue("MinPrice"))
	if r.FormValue("MinPrice") != "" && err != nil {
		log.Printf("無法解析 minPrice Error: %s", err)
		return
	}

	request := searchpage.VenueFilter{
		ActivityID:     activityID,
		NumberOfPeople: r.FormValue("NumberOfPeople"),
		City:           r.FormValue("City"),
		District:       r.FormValue("District"),
		MaxPrice:       maxPrice,
		MinPrice:       minPrice,
		VenueName:      r.FormValue("VenueName"),
		DayType:        r.FormValue("DayType"),
		RentTime:       r.FormValue("RentTime"),
	}

	vm := c.searchService.GetSearchPage(request)

	c.RenderTemplate(w, r, "searchpage", vm)
}

// 滾軸下拉刷新資料
func (c *SearchPageController) SearchPageLoading(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request URL: %s", r.URL.String())
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 解析表單
	if err := r.ParseForm(); err != nil {
		http.Error(w, "無法解析表單數據", http.StatusBadRequest)
		return
	}

	activityID, err := helper.StrToUint(r.FormValue("ActivityId"))
	if r.FormValue("ActivityId") != "" && err != nil {
		log.Printf("無法解析 ActivityId Error: %s", err)
		return
	}

	maxPrice, err := strconv.Atoi(r.FormValue("MaxPrice"))
	if r.FormValue("MaxPrice") != "" && err != nil {
		log.Printf("無法解析 MaxPrice Error: %s", err)
		return
	}

	minPrice, err := strconv.Atoi(r.FormValue("MinPrice"))
	if r.FormValue("MinPrice") != "" && err != nil {
		log.Printf("無法解析 MinPrice Error: %s", err)
		return
	}

	page, err := strconv.Atoi(r.FormValue("Page"))
	if err != nil {
		log.Printf("無法解析 Page Error: %s", err)
		return
	}

	request := searchpage.VenueFilter{
		ActivityID:     activityID,
		NumberOfPeople: r.FormValue("NumberOfPeople"),
		City:           r.FormValue("City"),
		District:       r.FormValue("District"),
		MaxPrice:       maxPrice,
		MinPrice:       minPrice,
		VenueName:      r.FormValue("VenueName"),
		DayType:        r.FormValue("DayType"),
		RentTime:       r.FormValue("RentTime"),
		Page:           page,
	}

	venueInfos := c.searchService.GetVenueInfos(request)

	response := searchpage.VenuePartialResponse{
		VenueInfos: venueInfos,
		EndOfData:  len(venueInfos) == 0,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
