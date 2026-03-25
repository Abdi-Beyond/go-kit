package models

type SubscriptionList struct {
	PK string `json:"PK,omitempty" dynamodbav:"PK"`
	SK string `json:"SK,omitempty" dynamodbav:"SK"`

	UserID         string `json:"UserID" dynamodbav:"user_id"`
	SubscriptionID string `json:"SubscriptionID,omitempty" dynamodbav:"subscription_id"`

	PlanID           string `json:"PlanID,omitempty"  dynamodbav:"plan_id"`
	PriceID          string `json:"PriceID,omitempty" dynamodbav:"price_id"`
	StripePriceID    string `json:"StripePriceID,omitempty" dynamodbav:"stripe_price_id"`
	StripeCustomerID string `json:"StripeCustomerID,omitempty" dynamodbav:"stripe_customer_id"`

	Status       string `json:"Status,omitempty" dynamodbav:"status"`
	StripeStatus string `json:"StripeStatus,omitempty" dynamodbav:"stripe_status"`

	Amount   int64  `json:"TokenAmount,omitempty" dynamodbav:"amount"`
	Currency string `json:"Currency,omitempty" dynamodbav:"currency"`

	CurrentPeriodStart string `json:"CurrentPeriodStart,omitempty" dynamodbav:"current_period_start"`
	CurrentPeriodEnd   string `json:"CurrentPeriodEnd,omitempty" dynamodbav:"current_period_end"`
	CancelAtPeriodEnd  bool   `json:"CancelAtPeriodEnd,omitempty" dynamodbav:"cancel_at_period_end"`

	Upgrade_plan   string `json:"Upgrade_plan" dynamodbav:"upgrade_plan_id"`
	Downgrade_plan string `json:"Downgrade_plan" dynamodbav:"downgrade_plan_id"`
	Current_plan   string `json:"Current_plan" dynamodbav:"current_plan_id"`

	Downgrade_effective_date string `json:"Downgrade_effective_date" dynamodbav:"downgrade_effective_date"`
	TrialEnd                 string `json:"TrialEnd,omitempty" dynamodbav:"trial_end"`
	CanceledAt               string `json:"CanceledAt,omitempty" dynamodbav:"canceled_at"`
	LatestInvoiceID          string `json:"LatestInvoiceID,omitempty" dynamodbav:"latest_invoice_id"`
	LastPaymentStatus        string `json:"LastPaymentStatus,omitempty" dynamodbav:"last_payment_status"`
	PlanName                 string `json:"PlanName,omitempty" dynamodbav:"plan_name"`
	PlanInterval             string `json:"PlanInterval,omitempty" dynamodbav:"plan_interval"`
	Metadata                 string `json:"Metadata,omitempty" dynamodbav:"metadata"`
	TotalTokenAmount         int64  `json:"TotalTokenAmount"`

	CreatedAt int64 `json:"CreatedAt,omitempty" dynamodbav:"created_at"`
	UpdatedAt int64 `json:"UpdatedAt,omitempty" dynamodbav:"updated_at"`
}
