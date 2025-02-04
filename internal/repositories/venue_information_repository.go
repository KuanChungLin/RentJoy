package repositories

import (
	"rentjoy/internal/dto/searchpage"
	interfaces "rentjoy/internal/interfaces/repositories"
	"rentjoy/internal/models"
	"rentjoy/pkg/helper"

	"gorm.io/gorm"
)

type VenueInformationRepository struct {
	*GenericRepository[models.VenueInformation]
}

func NewVenueInformationRepository(db *gorm.DB) interfaces.VenueInformationRepository {
	return &VenueInformationRepository{
		GenericRepository: NewGenericRepository[models.VenueInformation](db),
	}
}

func (r *VenueInformationRepository) FindExhibits() ([]models.VenueInformation, error) {
	var venues []models.VenueInformation

	err := r.GenericRepository.DB.Where("Status = ?", 1).
		Preload("Imgs", "Sort = ?", 0).
		Limit(6).
		Find(&venues).Error

	return venues, err
}

func (r *VenueInformationRepository) FindExhibitDESC() ([]models.ActivityType, map[uint][]models.VenueInformation, error) {
	var activityTypes []models.ActivityType
	err := r.GenericRepository.DB.Where("ActivityDescription IS NOT NULL").
		Limit(3).
		Find(&activityTypes).Error

	if err != nil {
		return nil, nil, err
	}

	venueMap := make(map[uint][]models.VenueInformation)

	for _, actType := range activityTypes {
		var venues []models.VenueInformation
		err := r.GenericRepository.DB.Where("Status = ?", 1).
			Joins("JOIN VenueActivities va ON va.VenueId = VenueInformations.Id").
			Where("va.ActivityId = ?", actType.ID).
			Preload("Imgs", "Sort = ?", 0).
			Preload("BillingRates").
			Limit(3).
			Find(&venues).Error

		if err != nil {
			return nil, nil, err
		}

		venueMap[actType.ID] = venues
	}

	return activityTypes, venueMap, nil
}

func (r *VenueInformationRepository) FindSearchPageInfos(filter searchpage.VenueFilter) ([]models.VenueInformation, error) {
	var venues []models.VenueInformation

	// 基本查詢 with preloads
	query := r.DB.Distinct("VenueInformations.*").Where("status = ?", 1).
		Preload("BillingRates.RateType").
		Preload("Management").
		Preload("Imgs", "Sort = ?", 0).
		Preload("Activities.Activity")

	// 動態添加篩選條件
	if filter.ActivityID != 0 {
		query = query.Joins("JOIN VenueActivities va ON va.VenueId = VenueInformations.Id").
			Joins("JOIN ActivityType a ON a.Id = va.ActivityId").
			Where("va.ActivityId = ?", filter.ActivityID)
	}

	if filter.VenueName != "" {
		query = query.Where("VenueName = ?", filter.VenueName)
	}

	if filter.City != "" {
		query = query.Where("City = ?", filter.City)
	}

	if filter.District != "" {
		query = query.Where("District = ?", filter.District)
	}

	if filter.NumberOfPeople != "" {
		max, min := helper.GetNumberOfPeopleFilter(filter.NumberOfPeople)
		query = query.Where("NumberOfPeople BETWEEN ? AND ?", min, max)
	}

	if filter.DayType != "" || filter.RentTime != "" || filter.MaxPrice != 0 || filter.MinPrice != 0 {
		query = query.Joins("JOIN BillingRate br ON br.VenueId = VenueInformations.Id")

		if filter.DayType != "" {
			dayList := helper.GetDayTypeCondition(filter.DayType)
			query = query.Where("DATEPART(WEEKDAY, br.DayOfWeek) IN ?", dayList)
		}

		if filter.RentTime != "" {
			startHour, endHour := helper.GetTimeSlotCondition(filter.RentTime)
			query = query.Where("DATEPART(HOUR, br.StartTime) BETWEEN DATEPART(HOUR, ?) AND DATEPART(HOUR, ?)", startHour, endHour)
		}

		if filter.MaxPrice != 0 {
			query = query.Where("br.Rate <= ?", filter.MaxPrice)
		}

		if filter.MinPrice != 0 {
			query = query.Where("br.Rate >= ?", filter.MinPrice)
		}
	}

	offset := (filter.Page - 1) * 4
	query = query.Offset(offset).Limit(4)

	err := query.Find(&venues).Error
	if err != nil {
		return nil, err
	}

	query.Debug()

	return venues, nil
}
