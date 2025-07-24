package dto

type SubscriptionResponse struct {
	ServiceName string  `json:"service_name" example:"Yandex Plus"`
	Price       int     `json:"price" example:"400"`
	UserID      string  `json:"user_id"  example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   string  `json:"start_date"  example:"07-2025"`
	EndDate     *string `json:"end_date,omitempty"  example:"12-2025"`
}
