package services

import (
	"fmt"
	"log"
	"rentjoy/internal/dto/homepage"
	"rentjoy/internal/dto/searchpage"
	repoInterfaces "rentjoy/internal/interfaces/repositories"
	serviceInterfaces "rentjoy/internal/interfaces/services"
	"rentjoy/internal/repositories"
	"rentjoy/pkg/helper"
	"sort"

	"gorm.io/gorm"
)

type SearchService struct {
	participantRangeRepo repoInterfaces.ParticipantRangeRepository
	activityTypeRepo     repoInterfaces.ActivityTypeRepository
	venueInformationRepo repoInterfaces.VenueInformationRepository
	budgetRepo           repoInterfaces.BudgetRepository
}

func NewSearchService(db *gorm.DB) serviceInterfaces.SearchPageService {
	return &SearchService{
		participantRangeRepo: repositories.NewParticipantRangeRepository(db),
		activityTypeRepo:     repositories.NewActivityTypeRepository(db),
		venueInformationRepo: repositories.NewVenueInformationRepository(db),
		budgetRepo:           repositories.NewBudgetRepository(db),
	}
}

// 取得搜尋頁資料
func (s *SearchService) GetSearchPage(filter searchpage.VenueFilter) *searchpage.SearchPage {
	return &searchpage.SearchPage{
		ActivityList:    s.GetActivities(),
		PeopleCountList: s.GetPeopleCounts(),
		MaxPriceList:    s.GetSearchPrice(),
		MinPriceList:    s.GetSearchPrice(),
		VenueInfos:      s.GetVenueInfos(filter),
		VenueFilter:     filter,
	}
}

// 取得場地活動搜尋欄位
func (s *SearchService) GetActivities() []homepage.Activity {
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

// 取得場地人數搜尋欄位
func (s *SearchService) GetPeopleCounts() []homepage.PeopleCount {
	ranges, err := s.participantRangeRepo.FindAll()
	if err != nil {
		log.Printf("PeopleCounts Get Error: %s", err)
		return []homepage.PeopleCount{}
	}

	PeopleCountList := make([]homepage.PeopleCount, len(ranges))
	for i, r := range ranges {
		PeopleCountList[i] = homepage.PeopleCount{
			PeopleCount: r.PeopleOfRange,
		}
	}

	return PeopleCountList
}

// 取得價格搜尋欄位
func (s *SearchService) GetSearchPrice() []int {
	ranges, err := s.budgetRepo.FindAll()
	if err != nil {
		log.Printf("SearchPrice Get Error: %s", err)
		return nil
	}

	priceList := make([]int, len(ranges))
	for i, r := range ranges {
		priceList[i] = helper.DecimalToIntRounded(r.Price)
	}

	return priceList
}

// 取得場地資訊
func (s *SearchService) GetVenueInfos(filter searchpage.VenueFilter) []searchpage.VenueInfo {
	ranges, err := s.venueInformationRepo.FindSearchPageInfos(filter)
	if err != nil {
		log.Printf("SearchPage Venue Get Error: %s", err)
		return []searchpage.VenueInfo{}
	}

	var venueList []searchpage.VenueInfo
	for _, venue := range ranges {
		var imgURL string
		if len(venue.Imgs) > 0 {
			imgURL = venue.Imgs[0].VenueImgPath
		}

		// 找出最低價格的計費方式
		var priceStr string
		if len(venue.BillingRates) > 0 {
			// 使用 sort.Slice 對 BillingRates 進行排序
			// func(i, j int) bool 是比較函數，返回 true 表示 i 應該排在 j 前面
			// LessThan 用於比較兩個 decimal 值的大小
			sort.Slice(venue.BillingRates, func(i, j int) bool {
				return venue.BillingRates[i].Rate.LessThan(venue.BillingRates[j].Rate)
			})

			lowestRate := venue.BillingRates[0]
			// 組合價格顯示字串："{收費方式說明} ${金額} 起"
			// StringFixed(0) 將 decimal 轉為字串，不保留小數位
			priceStr = fmt.Sprintf("%s $%s 起",
				lowestRate.RateType.RateDescription,
				lowestRate.Rate.StringFixed(0))
		}

		// 活動標籤
		var activityTags []homepage.Activity
		for _, va := range venue.Activities {
			activityTags = append(activityTags, homepage.Activity{
				ID:           va.Activity.ID,
				ActivityName: va.Activity.ActivityName,
			})
		}

		// 建立 ViewModel
		venueList = append(venueList, searchpage.VenueInfo{
			VenueImgUrl:    imgURL,
			VenueID:        venue.ID,
			VenueName:      venue.Name,
			VenueOwner:     venue.Management.ManagementName,
			VenueCity:      venue.City,
			VenueDistrict:  venue.District,
			NumberOfPeople: venue.NumOfPeople,
			VenuePrice:     priceStr,
			ActivityTags:   activityTags,
		})
	}

	return venueList
}
