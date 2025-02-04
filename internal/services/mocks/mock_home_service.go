package mocks

import (
	"log"
	"rentjoy/internal/dto/homepage"
)

type MockHomeService struct{}

func (s *MockHomeService) GetHomePage() homepage.HomePage {
	log.Println("取得HomePage假資料")
	return homepage.HomePage{
		ActivityList:           s.GetActivities(),
		PeopleCountList:        s.GetPeopleCounts(),
		GalleryList:            s.GetGalleries(),
		ExhibitList:            s.GetExhibits(),
		ExhibitDescriptionList: s.GetExhibitDescriptions(),
	}
}

func (s *MockHomeService) GetActivities() []homepage.Activity {
	log.Println("取得Activity假資料")
	return []homepage.Activity{
		{ID: 1, ActivityName: "課程講座"},
		{ID: 2, ActivityName: "派對"},
		{ID: 3, ActivityName: "會議"},
		{ID: 4, ActivityName: "聚會餐敘"},
		{ID: 5, ActivityName: "拍攝攝影"},
		{ID: 6, ActivityName: "發表會"},
		{ID: 7, ActivityName: "音樂表演"},
		{ID: 8, ActivityName: "靜態展覽"},
		{ID: 9, ActivityName: "教育訓練"},
		{ID: 10, ActivityName: "聚會"},
		{ID: 11, ActivityName: "運動"},
		{ID: 12, ActivityName: "工作"},
		{ID: 13, ActivityName: "私人談話"},
		{ID: 14, ActivityName: "親子活動"},
	}
}

func (s *MockHomeService) GetPeopleCounts() []homepage.PeopleCount {
	log.Println("取得PeopleCount假資料")
	return []homepage.PeopleCount{
		{PeopleCount: "1 - 10"},
		{PeopleCount: "11 - 20"},
		{PeopleCount: "21 - 40"},
		{PeopleCount: "41 - 60"},
		{PeopleCount: "61 - 80"},
		{PeopleCount: "81 - 100"},
		{PeopleCount: "101 - 200"},
		{PeopleCount: "201 - 300"},
		{PeopleCount: "301 - 400"},
		{PeopleCount: "401 - 500"},
		{PeopleCount: "500+"},
	}

}

func (s *MockHomeService) GetGalleries() []homepage.Gallery {
	log.Println("取得Gallery假資料")
	return []homepage.Gallery{
		{GalleryTitle: "課程講座", GalleryImgUrl: "/static/images/event_icons-01.jpeg", GalleryLinkUrl: "1"},
		{GalleryTitle: "派對", GalleryImgUrl: "/static/images/event_icons-02.jpeg", GalleryLinkUrl: "2"},
		{GalleryTitle: "會議", GalleryImgUrl: "/static/images/event_icons-03.jpeg", GalleryLinkUrl: "3"},
		{GalleryTitle: "聚餐餐敘", GalleryImgUrl: "/static/images/event_icons-04.jpeg", GalleryLinkUrl: "4"},
		{GalleryTitle: "拍攝攝影", GalleryImgUrl: "/static/images/event_icons-05.jpeg", GalleryLinkUrl: "5"},
		{GalleryTitle: "發表會", GalleryImgUrl: "/static/images/event_icons-06.jpeg", GalleryLinkUrl: "6"},
		{GalleryTitle: "音樂表演", GalleryImgUrl: "/static/images/event_icons-07.jpeg", GalleryLinkUrl: "7"},
		{GalleryTitle: "靜態展覽", GalleryImgUrl: "/static/images/event_icons-08.jpeg", GalleryLinkUrl: "8"},
		{GalleryTitle: "教育訓練", GalleryImgUrl: "/static/images/event_icons-09.jpeg", GalleryLinkUrl: "9"},
	}
}

func (s *MockHomeService) GetExhibits() []homepage.Exhibit {
	log.Println("取得Exhibit假資料")
	return []homepage.Exhibit{
		{
			ExhibitName: "南港展覽館捷運站走路2分鐘，明亮寬敞大會議室",
			SiteImgUrl:  "/static/images/siteimg1.jpg",
			SiteLinkUrl: "#",
		},
		{
			ExhibitName: "大包廂-包場計時收費",
			SiteImgUrl:  "/static/images/siteimg2.jpg",
			SiteLinkUrl: "#",
		},
		{
			ExhibitName: "Swiming Taiwan 藝文空間!!美術園區+台中場地出租[Swiming Taiwan藝文空間]!",
			SiteImgUrl:  "/static/images/siteimg3.jpg",
			SiteLinkUrl: "#",
		},
		{
			ExhibitName: "小草屋|椰子草101",
			SiteImgUrl:  "/static/images/siteimg4.jpg",
			SiteLinkUrl: "#",
		},
		{
			ExhibitName: "九十九yucci|複合式商務空間_北車京站館301",
			SiteImgUrl:  "/static/images/siteimg5.jpg",
			SiteLinkUrl: "#",
		},
		{
			ExhibitName: "HOWSHUAI台北新生C空間",
			SiteImgUrl:  "/static/images/siteimg6.jpg",
			SiteLinkUrl: "#",
		},
	}
}

func (s *MockHomeService) GetExhibitDescriptions() []homepage.ExhibitDESC {
	log.Println("取得ExhibitDESC假資料")
	return []homepage.ExhibitDESC{
		{
			ExhibitTitle: "課程講座",
			Description:  "類型場地可舉辦課程、演講、說明會、發表會等",
			Exhibits: []homepage.Exhibit{
				{SiteImgUrl: "/static/images/siteimg1.jpg", SiteLinkUrl: "#", SiteLocation: "台北市", CapaCity: "30", Price: "$5200", PerTime: "4hr"},
				{SiteImgUrl: "/static/images/siteimg2.jpg", SiteLinkUrl: "#", SiteLocation: "台北市", CapaCity: "40", Price: "$1200", PerTime: "14hr"},
				{SiteImgUrl: "/static/images/siteimg3.jpg", SiteLinkUrl: "#", SiteLocation: "台南市", CapaCity: "4", Price: "$220", PerTime: "hr"},
			}},
		{
			ExhibitTitle: "會議",
			Description:  "類型場地可舉辦會議、研討會、讀書會等",
			Exhibits: []homepage.Exhibit{
				{SiteImgUrl: "/static/images/siteimg4.jpg", SiteLinkUrl: "#", SiteLocation: "台北市", CapaCity: "400", Price: "$70000", PerTime: "5hr"},
				{SiteImgUrl: "/static/images/siteimg5.jpg", SiteLinkUrl: "#", SiteLocation: "台北市", CapaCity: "60", Price: "$5775", PerTime: "4hr"},
				{SiteImgUrl: "/static/images/siteimg6.jpg", SiteLinkUrl: "#", SiteLocation: "高雄市", CapaCity: "20", Price: "$400", PerTime: "hr"},
			}},
		{
			ExhibitTitle: "聚會派對",
			Description:  "類型場地可舉辦校友會、慶生會、同樂會等",
			Exhibits: []homepage.Exhibit{
				{SiteImgUrl: "/static/images/siteimg7.jpg", SiteLinkUrl: "#", SiteLocation: "台北市", CapaCity: "10", Price: "$360", PerTime: "hr"},
				{SiteImgUrl: "/static/images/siteimg8.jpg", SiteLinkUrl: "#", SiteLocation: "台中市", CapaCity: "7", Price: "$460", PerTime: "hr"},
			}},
	}
}
