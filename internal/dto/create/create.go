package create

type SpaceType struct {
	ID         int         `json:"id"`
	TypeName   string      `json:"typeName"`
	SpaceInfos []SpaceInfo `json:"spaceInfos"`
}

type SpaceInfo struct {
	ID           int    `json:"id"`
	FacilityName string `json:"facilityName"`
}

type ActivityInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type EquipmentType struct {
	ID             int             `json:"id"`
	TypeName       string          `json:"typeName"`
	EquipmentInfos []EquipmentInfo `json:"equipmentInfos"`
}

type EquipmentInfo struct {
	ID            int    `json:"id"`
	EquipmentName string `json:"equipmentName"`
	Quantity      int    `json:"quantity"`
	Description   string `json:"description"`
}

type ManagerInfo struct {
	ID                  int    `json:"id"`
	ManagerName         string `json:"managerName"`
	ManagerContact      string `json:"managerContact"`
	ManagerDescription  string `json:"managerDescription"`
	ManagerPublicPhone  string `json:"managerPublicPhone"`
	ManagerPrivatePhone string `json:"managerPrivatePhone"`
	ManagerImgUrl       string `json:"managerImgUrl"`
}
