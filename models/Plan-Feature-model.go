package models

// ─────────────────────────────────────────────
// ENUMS
// ─────────────────────────────────────────────

type FeatureStatus string

const (
	FeatureStatusDraft      FeatureStatus = "draft"
	FeatureStatusBeta       FeatureStatus = "beta"
	FeatureStatusActive     FeatureStatus = "active"
	FeatureStatusDeprecated FeatureStatus = "deprecated"
	FeatureStatusDisabled   FeatureStatus = "disabled"
)

type TokenCostModel string

const (
	TokenCostModelPerUnit TokenCostModel = "per_unit"
	TokenCostModelFlat    TokenCostModel = "flat"
	TokenCostModelTiered  TokenCostModel = "tiered"
	TokenCostModelFree    TokenCostModel = "free"
)

type LimitType string

const (
	LimitTypeConcurrent LimitType = "concurrent" // max at the same time
	LimitTypeMonthly    LimitType = "monthly"    // rolling counter, resets each cycle
	LimitTypeHardCap    LimitType = "hard_cap"   // absolute max, never resets
)

type LimitReset string

const (
	LimitResetMonthly LimitReset = "monthly"
	LimitResetNever   LimitReset = "never"
)

type ExceedBehavior string

const (
	ExceedBehaviorBlock   ExceedBehavior = "block"
	ExceedBehaviorQueue   ExceedBehavior = "queue"
	ExceedBehaviorNotify  ExceedBehavior = "notify"
	ExceedBehaviorOverage ExceedBehavior = "overage"
)

// ─────────────────────────────────────────────
// FEATURE KEYS
// single source of truth — import this, never use raw strings
// ─────────────────────────────────────────────

type FeatureKey string

const (
	FeatureBasicAPI           FeatureKey = "basic_api"
	FeatureBatchProcessing    FeatureKey = "batch_processing"
	FeatureExport             FeatureKey = "export"
	FeatureWebhooks           FeatureKey = "webhooks"
	FeatureCustomDomains      FeatureKey = "custom_domains"
	FeaturePriorityProcessing FeatureKey = "priority_processing"
	FeatureAnalytics          FeatureKey = "analytics"
	FeatureAPIAccess          FeatureKey = "api_access"
)

// ─────────────────────────────────────────────
// PLAN FEATURE
// ─────────────────────────────────────────────

type PlanFeature struct {
	// ── Identity ─────────────────────────────
	PK string `dynamodbav:"PK"             json:"-"`
	SK string `dynamodbav:"SK"             json:"-"`

	PlanID         string     `dynamodbav:"plan_id"         json:"plan_id" binding:"required"`
	PlanVersion    int        `dynamodbav:"plan_version"    json:"plan_version"`
	FeatureKey     FeatureKey `dynamodbav:"feature_key"     json:"feature_key" binding:"required"`
	FeatureVersion int        `dynamodbav:"feature_version" json:"feature_version"` // increments per config change on this feature

	// ── State ────────────────────────────────
	Status  FeatureStatus `dynamodbav:"status"  json:"status"`
	Enabled bool          `dynamodbav:"enabled" json:"enabled"`

	// ── Limits ───────────────────────────────
	Limit            *int64         `dynamodbav:"limit"              json:"limit"` // nil = unlimited
	LimitType        *LimitType     `dynamodbav:"limit_type"         json:"limit_type,omitempty"`
	LimitUnit        string         `dynamodbav:"limit_unit"         json:"limit_unit,omitempty"` // jobs | records | requests | domains
	Reset            *LimitReset    `dynamodbav:"reset"              json:"reset,omitempty"`
	BehaviorOnExceed ExceedBehavior `dynamodbav:"behavior_on_exceed" json:"behavior_on_exceed"`

	// ── Overage ──────────────────────────────
	OverageAllowed bool `dynamodbav:"overage_allowed"    json:"overage_allowed"`
	//OverageTokenCost int64 `dynamodbav:"overage_token_cost" json:"overage_token_cost"` // higher cost per unit when over limit

	// // ── Visibility ───────────────────────────
	// IsPublic        bool   `dynamodbav:"is_public"          json:"is_public"` // show on pricing page
	// IsAddon         bool   `dynamodbav:"is_addon"           json:"is_addon"`  // sold separately
	//AddonTokenGrant *int64 `dynamodbav:"addon_token_grant"  json:"addon_token_grant,omitempty"`

	// ── Rollout ──────────────────────────────
	EffectiveAt     int64 `dynamodbav:"effective_at"     json:"effective_at"`               // schedule future changes
	ExpiresAt       int64 `dynamodbav:"expires_at"       json:"expires_at,omitempty"`       // for temporary promotions
	IsDraft         bool  `dynamodbav:"is_draft"         json:"is_draft"`                   // true = staged, not live
	PreviousVersion *int  `dynamodbav:"previous_version" json:"previous_version,omitempty"` // pointer for rollback

	// ── Audit ────────────────────────────────
	CreatedAt int64  `dynamodbav:"created_at"    json:"created_at"`
	UpdatedAt int64  `dynamodbav:"updated_at"    json:"updated_at"`
	UpdatedBy string `dynamodbav:"updated_by"    json:"updated_by"`
}
