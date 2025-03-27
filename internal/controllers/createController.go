package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"rentjoy/internal/dto/create"
	serviceInterfaces "rentjoy/internal/interfaces/services"
	"rentjoy/internal/middleware"
	"rentjoy/pkg/helper"
)

type CreateController struct {
	BaseController
	createService serviceInterfaces.CreateService
}

func NewCreateController(createService serviceInterfaces.CreateService, templates map[string]*template.Template) *CreateController {
	return &CreateController{
		BaseController: NewBaseController(templates),
		createService:  createService,
	}
}

// 場地新增頁面
func (c *CreateController) CreateVenue(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		c.RenderTemplate(w, r, "create_venue", create.CreateForm{})
	} else if r.Method == http.MethodPost {
		// 解析 multipart/form-data 表單資料
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			http.Error(w, "Failed to parse form", http.StatusBadRequest)
			return
		}

		form := &create.CreateForm{
			// 基本信息
			VenueTypeId:        helper.StringToInt(r.FormValue("facilityId")),
			SelectedActivities: helper.StringSliceToIntSlice(r.Form["selectedActivities"]),
			VenueName:          r.FormValue("VenueName"),
			VenueRule:          r.FormValue("VenueRule"),
			UnsubscribeRule:    r.FormValue("UnsubscribeRule"),

			// 位置信息
			CitySelect:            r.FormValue("citySelect"),
			DistrictSelect:        r.FormValue("districtSelect"),
			Address:               r.FormValue("areaInfoStreetAddress"),
			TransportationMRT:     r.FormValue("transportInfoMRT"),
			TransportationBus:     r.FormValue("transportInfoBus"),
			TransportationParking: r.FormValue("transportInfoParking"),

			// 空间配置
			SpaceSize:      helper.StringToInt(r.FormValue("spaceSize")),
			NumberOfPeople: helper.StringToInt(r.FormValue("numberOfSpace")),

			// 设备信息
			Equipments: parseEquipments(r),

			// 价格设置
			HourPricing:   parseHourPricing(r),
			PeriodPricing: parsePeriodPricing(r),

			// 管理员信息
			ManagerID: helper.StringToInt(r.FormValue("manager")),

			// 图片信息
			VenueImgs: r.MultipartForm.File["VenueImgs"],
		}

		userId, ok := middleware.GetUserIDFromContext(r.Context())
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if err := c.createService.CreateVenue(userId, form); err != nil {
			http.Error(w, "Failed to create venue", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/Manage/VenueManagement", http.StatusSeeOther)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

}

// 取得空間類型資料
func (c *CreateController) GetSpaceTypesData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	spaceTypes, err := c.createService.GetSpaceTypes()
	if err != nil {
		http.Error(w, "Failed to get space types", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	if err := json.NewEncoder(w).Encode(spaceTypes); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// 取得活動資料
func (c *CreateController) GetActivitiesData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	activities, err := c.createService.GetActivities()
	if err != nil {
		http.Error(w, "Failed to get activities", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	if err := json.NewEncoder(w).Encode(activities); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// 取得設備類型資料
func (c *CreateController) GetEquipmentTypesData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	equipmentTypes, err := c.createService.GetEquipmentTypes()
	if err != nil {
		http.Error(w, "Failed to get equipment types", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	if err := json.NewEncoder(w).Encode(equipmentTypes); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// 取得管理員資料
func (c *CreateController) GetManagersInfoData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userId, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	managers, err := c.createService.GetManagers(userId)
	if err != nil {
		http.Error(w, "Failed to get managers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	if err := json.NewEncoder(w).Encode(managers); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// 解析設備資料
func parseEquipments(r *http.Request) []create.Equipment {
	var equipments []create.Equipment
	i := 0
	for {
		id := r.FormValue(fmt.Sprintf("EquipmentInfos[%d].Id", i))
		if id == "" {
			break
		}
		equipment := create.Equipment{
			Id:          helper.StringToInt(id),
			Quantity:    r.FormValue(fmt.Sprintf("EquipmentInfos[%d].Quantity", i)),
			Description: r.FormValue(fmt.Sprintf("EquipmentInfos[%d].Description", i)),
		}
		equipments = append(equipments, equipment)
		i++
	}
	return equipments
}

// 解析小時定價
func parseHourPricing(r *http.Request) create.HourPricingConfig {
	// 從表單中取得 JSON 字串
	hourPricingJSON := r.FormValue("HourPricing")
	if hourPricingJSON == "" {
		return create.HourPricingConfig{}
	}

	// 解析 JSON
	var hourPricing create.HourPricingConfig
	if err := json.Unmarshal([]byte(hourPricingJSON), &hourPricing); err != nil {
		log.Printf("Failed to parse hour pricing JSON: %v", err)
		return create.HourPricingConfig{}
	}

	return hourPricing
}

// 解析時段定價
func parsePeriodPricing(r *http.Request) create.PeriodPricingConfig {
	// 從表單中取得 JSON 字串
	periodPricingJSON := r.FormValue("PeriodPricing")
	if periodPricingJSON == "" {
		return create.PeriodPricingConfig{}
	}

	// 解析 JSON
	var periodPricing create.PeriodPricingConfig
	if err := json.Unmarshal([]byte(periodPricingJSON), &periodPricing); err != nil {
		log.Printf("Failed to parse period pricing JSON: %v", err)
		return create.PeriodPricingConfig{}
	}

	return periodPricing
}
