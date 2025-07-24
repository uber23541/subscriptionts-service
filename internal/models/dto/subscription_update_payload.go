package dto

type SubscriptionUpdatePayload struct {
	Price     int     `json:"price" binding:"required,gte=0" example:"400"`
	StartDate string  `json:"start_date" binding:"required" example:"07-2025"`
	EndDate   *string `json:"end_date" example:"12-2025"`
}
