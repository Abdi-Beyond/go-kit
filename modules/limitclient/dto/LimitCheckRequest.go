package dto

type LimitCheckRequest struct {
	ServiceType string `json:"serviceType" binding:"required"`
	PropertyID  string `json:"property_id,omitempty"`
}
