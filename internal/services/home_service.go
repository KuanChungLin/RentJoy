package services

import (
	"fmt"
	"log"
	"math"
	"rentjoy/internal/dto/homepage"
	repoInterfaces "rentjoy/internal/interfaces/repositories"
	serviceInterfaces "rentjoy/internal/interfaces/services"
	"rentjoy/internal/repositories"
	"strconv"

	"gorm.io/gorm"
)

type HomeService struct {
	participantRangeRepo repoInterfaces.ParticipantRangeRepository
	activityTypeRepo     repoInterfaces.ActivityTypeRepository
	venueInformationRepo repoInterfaces.VenueInformationRepository
}

func NewHomeService(db *gorm.DB) serviceInterfaces.HomeService {
	return &HomeService{
		participantRangeRepo: repositories.NewParticipantRangeRepository(db),
		activityTypeRepo:     repositories.NewActivityTypeRepository(db),
		venueInformationRepo: repositories.NewVenueInformationRepository(db),
	}
}

// 取得整個首頁所需的資料
func (s *HomeService) GetHomePage() *homepage.HomePage {
	return &homepage.HomePage{
		ActivityList:           s.GetActivities(),
		PeopleCountList:        s.GetPeopleCounts(),
		GalleryList:            s.GetGalleries(),
		ExhibitList:            s.GetExhibits(),
		ExhibitDescriptionList: s.GetExhibitDESC(),
	}
}

// 取得活動類型資料
func (s *HomeService) GetActivities() []homepage.Activity {
	ranges, err := s.activityTypeRepo.FindAll()
	if err != nil {
		log.Printf("ActivityType Get Error: %s", err)
		return []homepage.Activity{}
	}

	activityTypeList := make([]homepage.Activity, len(ranges))
	for i, r := range ranges {
		activityTypeList[i] = homepage.Activity{
			ID:           r.ID,
			ActivityName: r.ActivityName,
		}
	}
	return activityTypeList
}

// 取得活動人數區間資料
func (s *HomeService) GetPeopleCounts() []homepage.PeopleCount {
	ranges, err := s.participantRangeRepo.FindAll()
	if err != nil {
		log.Printf("PeopleOfRange Get Error: %s", err)
		return []homepage.PeopleCount{}
	}

	peopleCountList := make([]homepage.PeopleCount, len(ranges))
	for i, r := range ranges {
		peopleCountList[i] = homepage.PeopleCount{
			PeopleCount: r.PeopleOfRange,
		}
	}

	return peopleCountList
}

// 取得活動性質 sweeper 內容物
func (s *HomeService) GetGalleries() []homepage.Gallery {
	ranges, err := s.activityTypeRepo.FindAll()
	if err != nil {
		log.Printf("ActivityType Get Error: %s", err)
		return []homepage.Gallery{}
	}

	galleryList := make([]homepage.Gallery, len(ranges))
	for i, r := range ranges {
		galleryList[i] = homepage.Gallery{
			GalleryTitle:   r.ActivityName,
			GalleryImgUrl:  r.ActivityIcon,
			GalleryLinkUrl: "/SearchPage?ActivityId=" + strconv.Itoa(int(r.ID)),
		}
	}
	return galleryList
}

func (s *HomeService) GetExhibits() []homepage.Exhibit {
	ranges, err := s.venueInformationRepo.FindExhibits()
	if err != nil {
		log.Printf("Exhibits Get Error: %s", err)
	}

	exhibitList := make([]homepage.Exhibit, len(ranges))
	for i, r := range ranges {
		imgUrl := r.Imgs[0].VenueImgPath

		if imgUrl == "" {
			log.Printf("場地 %d 圖片取得錯誤", r.ID)
		}

		exhibitList[i] = homepage.Exhibit{
			ExhibitName: r.Name,
			SiteImgUrl:  imgUrl,
			SiteLinkUrl: "/Venue/VenuePage?venueId=" + strconv.Itoa(int(r.ID)),
		}
	}
	return exhibitList
}

func (s *HomeService) GetExhibitDESC() []homepage.ExhibitDESC {
	// 從 Repository 獲取資料
	activityTypes, venueMap, err := s.venueInformationRepo.FindExhibitDESC()
	if err != nil {
		log.Printf("Error getting exhibit descriptions: %v", err)
		return []homepage.ExhibitDESC{}
	}

	var result []homepage.ExhibitDESC

	// 轉換資料
	for _, actType := range activityTypes {
		venues := venueMap[actType.ID]
		if len(venues) == 0 {
			continue
		}

		var exhibits []homepage.Exhibit
		for _, venue := range venues {
			// 處理圖片
			var imgURL string
			if len(venue.Imgs) > 0 {
				imgURL = venue.Imgs[0].VenueImgPath
			}

			// 處理計費資訊
			var price, perTime string
			if len(venue.BillingRates) > 0 {
				rate := venue.BillingRates[0]
				price = (rate.Rate).String()

				duration := rate.EndTime.Sub(rate.StartTime)
				hours := int(math.Round(duration.Hours()))
				perTime = fmt.Sprintf("%d hr", hours)
			}

			exhibits = append(exhibits, homepage.Exhibit{
				SiteImgUrl:   imgURL,
				SiteLinkUrl:  fmt.Sprintf("/Venue/VenuePage?venueId=%d", venue.ID),
				SiteLocation: venue.City,
				CapaCity:     strconv.Itoa(int(venue.NumOfPeople)),
				Price:        price,
				PerTime:      perTime,
			})
		}

		result = append(result, homepage.ExhibitDESC{
			ExhibitTitle: actType.ActivityName,
			Description:  actType.ActivityDescription,
			Exhibits:     exhibits,
		})
	}

	return result
}
