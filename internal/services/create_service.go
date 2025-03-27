package services

import (
	"rentjoy/internal/dto/create"
	repoInterfaces "rentjoy/internal/interfaces/repositories"
	serviceInterfaces "rentjoy/internal/interfaces/services"
	"rentjoy/internal/models"
	"rentjoy/internal/repositories"
	"rentjoy/pkg/helper"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type CreateService struct {
	spaceTypeRepo     repoInterfaces.SpaceTypeRepository
	activityTypeRepo  repoInterfaces.ActivityTypeRepository
	deviceTypeRepo    repoInterfaces.DeviceTypeRepository
	managementRepo    repoInterfaces.ManagementRepository
	cloudinaryService serviceInterfaces.CloudinaryService
	DB                *gorm.DB
}

func NewCreateService(db *gorm.DB, cloudinaryService serviceInterfaces.CloudinaryService) serviceInterfaces.CreateService {
	return &CreateService{
		spaceTypeRepo:     repositories.NewSpaceTypeRepository(db),
		activityTypeRepo:  repositories.NewActivityTypeRepository(db),
		deviceTypeRepo:    repositories.NewDeviceTypeRepository(db),
		managementRepo:    repositories.NewManagementRepository(db),
		cloudinaryService: cloudinaryService,
		DB:                db,
	}
}

// 取得場地新增頁內的場地類型
func (s *CreateService) GetSpaceTypes() ([]create.SpaceType, error) {
	spaceTypes, err := s.spaceTypeRepo.FindSpaceAndFacility()
	if err != nil {
		return nil, err
	}

	var spaceTypeDto []create.SpaceType
	var dtoInfos []create.SpaceInfo

	for i, spaceType := range spaceTypes {
		spaceTypeDto = append(spaceTypeDto, create.SpaceType{
			ID:       int(spaceType.ID),
			TypeName: spaceType.TypeName,
		})

		for _, venueType := range spaceType.VenueTypes {
			dtoInfos = append(dtoInfos, create.SpaceInfo{
				ID:           int(venueType.ID),
				FacilityName: venueType.TypeName,
			})
		}

		spaceTypeDto[i].SpaceInfos = dtoInfos
		dtoInfos = nil
	}
	return spaceTypeDto, nil
}

// 取得場地新增頁內的活動類型
func (s *CreateService) GetActivities() ([]create.ActivityInfo, error) {
	activities, err := s.activityTypeRepo.FindAll()
	if err != nil {
		return nil, err
	}

	var activityDto []create.ActivityInfo

	for _, activity := range activities {
		activityDto = append(activityDto, create.ActivityInfo{
			ID:   int(activity.ID),
			Name: activity.ActivityName,
		})
	}

	return activityDto, nil
}

// 取得場地新增頁內的設備類型
func (s *CreateService) GetEquipmentTypes() ([]create.EquipmentType, error) {
	deviceTypes, err := s.deviceTypeRepo.FindTypeAndItems()
	if err != nil {
		return nil, err
	}

	var equipmentTypes []create.EquipmentType
	var equipmentInfos []create.EquipmentInfo

	for i, deviceType := range deviceTypes {
		equipmentTypes = append(equipmentTypes, create.EquipmentType{
			ID:       int(deviceType.ID),
			TypeName: deviceType.DeviceTypeName,
		})

		for _, deviceItem := range deviceType.DeviceItems {
			equipmentInfos = append(equipmentInfos, create.EquipmentInfo{
				ID:            int(deviceItem.ID),
				EquipmentName: deviceItem.DeviceName,
			})
		}

		equipmentTypes[i].EquipmentInfos = equipmentInfos
		equipmentInfos = nil
	}

	return equipmentTypes, nil
}

// 取得場地新增頁內的經理人
func (s *CreateService) GetManagers(userId uint) ([]create.ManagerInfo, error) {
	managements, err := s.managementRepo.FindByUserId(userId)
	if err != nil {
		return nil, err
	}

	var managerInfos []create.ManagerInfo

	for _, management := range managements {
		managerInfos = append(managerInfos, create.ManagerInfo{
			ID:                  int(management.ID),
			ManagerName:         management.ManagementName,
			ManagerContact:      management.Contact,
			ManagerDescription:  management.ManagementDescription,
			ManagerPublicPhone:  management.PublicNumber,
			ManagerPrivatePhone: management.PrivateNumber,
			ManagerImgUrl:       management.AvatarImgLinkPath,
		})
	}

	return managerInfos, nil
}

// 新增場地
func (s *CreateService) CreateVenue(userId uint, form *create.CreateForm) error {
	// 建立交易
	tx := s.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 準備資料
	var activities []models.VenueActivity
	var equipments []models.VenueDevice
	var venueImages []models.VenueImg
	var hourPricing []models.BillingRate
	var periodPricing []models.BillingRate

	// 處理活動資料
	for _, activity := range form.SelectedActivities {
		activities = append(activities, models.VenueActivity{
			ActivityID: uint(activity),
		})
	}

	// 處理設備資料
	for _, equipment := range form.Equipments {
		equipments = append(equipments, models.VenueDevice{
			DeviceItemID:      uint(equipment.Id),
			Count:             helper.StringToInt(equipment.Quantity),
			DeviceDescription: equipment.Description,
		})
	}

	// 處理小時定價資料
	for _, pricing := range form.HourPricing.PricingSettings {
		startTime, err := time.Parse(time.RFC3339, pricing.StartTime)
		if err != nil {
			return err
		}
		endTime, err := time.Parse(time.RFC3339, pricing.EndTime)
		if err != nil {
			return err
		}

		hourPricing = append(hourPricing, models.BillingRate{
			RateTypeID:   1,
			Rate:         decimal.NewFromFloat(pricing.Price),
			DayOfWeek:    time.Weekday(pricing.Day),
			StartTime:    startTime,
			EndTime:      endTime,
			MinRentHours: form.HourPricing.LeastRentHours,
		})
	}

	// 處理時段定價資料
	for _, pricing := range form.PeriodPricing.PricingSettings {
		startTime, err := time.Parse(time.RFC3339, pricing.StartTime)
		if err != nil {
			return err
		}
		endTime, err := time.Parse(time.RFC3339, pricing.EndTime)
		if err != nil {
			return err
		}

		periodPricing = append(periodPricing, models.BillingRate{
			RateTypeID: 2,
			Rate:       decimal.NewFromFloat(pricing.Price),
			DayOfWeek:  time.Weekday(pricing.Day),
			StartTime:  startTime,
			EndTime:    endTime,
		})
	}

	// 上傳圖片
	imgs, err := s.cloudinaryService.UploadImages(form.VenueImgs)
	if err != nil {
		tx.Rollback()
		return err
	}

	// 處理圖片資料
	for i, img := range imgs {
		venueImages = append(venueImages, models.VenueImg{
			VenueImgPath: img.ImageURL,
			Sort:         uint(i),
			PublicID:     img.PublicID,
		})
	}

	latitude, longitude := helper.GetCoordinates(form.CitySelect, form.DistrictSelect, form.Address)

	// 創建場地資訊
	venue := models.VenueInformation{
		VenueTypeID:     uint(form.VenueTypeId),
		ManagementID:    uint(form.ManagerID),
		Name:            form.VenueName,
		Rules:           form.VenueRule,
		UnsubscribeRule: form.UnsubscribeRule,
		Address:         form.Address,
		Status:          2,
		City:            form.CitySelect,
		District:        form.DistrictSelect,
		MRTInfo:         form.TransportationMRT,
		BusInfo:         form.TransportationBus,
		ParkInfo:        form.TransportationParking,
		NumOfPeople:     form.NumberOfPeople,
		SpaceSize:       form.SpaceSize,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		Latitude:        latitude,
		Longitude:       longitude,
		OwnerID:         userId,
		Imgs:            venueImages,
		Activities:      activities,
		Devices:         equipments,
		BillingRates:    append(hourPricing, periodPricing...),
	}

	// 創建場地與相關資料
	if err := tx.Create(&venue).Error; err != nil {
		return err
	}

	// 執行交易
	return tx.Commit().Error
}
