package homepage

type HomePage struct {
	ActivityList []Activity `json:"activityList"`
	// Cities []City `json:"cities"`
	GalleryList            []Gallery     `json:"galleryList"`
	ExhibitList            []Exhibit     `json:"exhibitList"`
	ExhibitDescriptionList []ExhibitDESC `json:"exhibitDescriptionList"`
	PeopleCountList        []PeopleCount `json:"peopleCountList"`
}

type Activity struct {
	ID           uint   `json:"id"`
	ActivityName string `json:"activityName"`
}

type City struct {
	CityName string `json:"cityName"`
}

type Gallery struct {
	GalleryTitle   string `json:"galleryTitle"`
	GalleryImgUrl  string `json:"galleryImgUrl"`
	GalleryLinkUrl string `json:"galleryLinkUrl"`
}

type Exhibit struct {
	ExhibitName  string `json:"exhibitName"`
	SiteImgUrl   string `json:"siteImgUrl"`
	SiteLocation string `json:"siteLocation"`
	CapaCity     string `json:"capaCity"`
	Price        string `json:"price"`
	PerTime      string `json:"perTime"`
	SiteLinkUrl  string `json:"siteLinkUrl"`
}

type ExhibitDESC struct {
	ExhibitTitle string    `json:"exhibitTitle"`
	Description  string    `json:"description"`
	Exhibits     []Exhibit `json:"exhibits"`
}

type PeopleCount struct {
	PeopleCount string `json:"peopleCount"`
}
