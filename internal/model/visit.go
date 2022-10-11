package model

// Visit represent visit data model
type Visit struct {
	VisitId    uint32 `json:"visit_id" validate:"required"`
	LocationId uint32 `json:"location_id" validate:"required"`
	UserId     uint32 `json:"user_id" validate:"required"`
	VisitedAt  string `json:"visited_at" validate:"required,datetime=2006-01-02"`
	Mark       uint8  `json:"mark" validate:"required,min=0,max=5"`
}

// UserVisit represent user visit data model
type UserVisit struct {
	Place     string `json:"place"`
	Country   string `json:"country"`
	VisitedAt string `json:"visited_at"`
	Mark      uint8  `json:"mark"`
}

// UserVisits represent user visits list data model
type UserVisits struct {
	Visits []UserVisit `json:"visits"`
}
