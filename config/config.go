package config

import (
	"log"
	"os"
	"strconv"

	"github.com/Abdi-Beyond/go-kit/core/auth"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type MAX_UNITS_PER_SERVICE struct {
	CompsUnit  string
	CompsToken string
}

type DynamoDBTables struct {
	Subscriptions  string
	Idempotency    string
	Ledger         string
	Plan           string
	Wallet         string
	Orders         string
	GSIOrders      string
	Notifications  string
	ServiceCatalog string
	Reservations   string
	UserCore       string

	UserPropertyCore   string
	UserSharedProperty string
	Folders            string
	PlanFeatures       string
}

type StripeConfig struct {
	SecretKey     string
	WebhookSecret string
}

type SQSConfig struct {
	PROD_COMPS_SQS string
	DEV_COMPS_SQS  string
}

type CognitoConfig struct {
	UserPoolID   string
	ClientID     string
	Region       string
	JWKSEndpoint string
}

type PlanConfig struct {
	ID         string
	Price      string
	Price_PROD string
	Name       string
}

type Config struct {
	// APP Env
	APP_ENV string
	// AWS Credentials (ONLY for config - not for client creation)
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSRegion          string

	// DynamoDB Tables

	DynamoDBTables DynamoDBTables

	// MAX Units per service
	MAX_UNITS_PER_SERVICE MAX_UNITS_PER_SERVICE

	//
	AWSSQS SQSConfig

	// Stripe
	Stripe StripeConfig

	// Server
	ServerAddress string
	ServerPort    int

	// Plans - organized better
	Plans struct {
		Free    PlanConfig
		Scale   PlanConfig
		Core    PlanConfig
		Team    PlanConfig
		Starter PlanConfig

		// Yearly plans
		CoreYearly  PlanConfig
		ScaleYearly PlanConfig
		TeamYearly  PlanConfig

		// Token packs
		Tokens50  PlanConfig
		Tokens100 PlanConfig
		Tokens250 PlanConfig
		Tokens500 PlanConfig
	}

	// URLs
	AppBaseURL string
	SuccessURL string
	CancelURL  string

	//API KEYS
	APIKey          string
	TokenGatewayKey string
	ServiceKey      string
	// Cognito
	Cognito CognitoConfig

	// Runtime dependencies (set later)
	Authenticator *auth.CognitoAuthenticator
}

func LoadConfig() (*Config, error) {
	LoadEnv()

	cfg := &Config{
		APP_ENV:            getEnv("APP_ENV", ""),
		AWSAccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
		AWSSecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
		AWSRegion:          getEnv("AWS_REGION", ""),

		DynamoDBTables: DynamoDBTables{
			Subscriptions:      getEnv("DYNAMODB_SUBSCRIPTIONS_TABLE", ""),
			Idempotency:        getEnv("DYNAMODB_IDEMPOTENCY_TABLE", ""),
			Ledger:             getEnv("DYNAMODB_LEDGER_TABLE", ""),
			Plan:               getEnv("DYNAMODB_PLAN_TABLE_PROD", ""),
			Wallet:             getEnv("DYNAMODB_WALLET_TABLE", ""),
			Orders:             getEnv("DYNAMODB_ORDERS_TABLE", ""),
			GSIOrders:          getEnv("GSI_ORDERS_TABLE", ""),
			Notifications:      getEnv("DYNAMODB_NOTIFICATIONS_TABLE", ""),
			ServiceCatalog:     getEnv("DYNAMODB_SERVICE_TABLE", ""),
			Reservations:       getEnv("DYNAMODB_WALLET_RESERVATION_TABLE", ""),
			UserCore:           getEnv("DYNAMODB_USER_CORE_TABLE", ""),
			UserPropertyCore:   getEnv("DYNAMODB_USER_PROPERTY_CORE_TABLE", ""),
			UserSharedProperty: getEnv("DYNAMODB_USER_SHARED_PROPERTY_TABLE", ""),
			Folders:            getEnv("DYNAMODB_FOLDERS_TABLE", ""),
			PlanFeatures:       getEnv("DYNAMODB_PLAN_FEATURES_TABLE", ""),
		},
		MAX_UNITS_PER_SERVICE: MAX_UNITS_PER_SERVICE{
			CompsUnit:  getEnv("MAX_UNITS_COMPS", ""),
			CompsToken: getEnv("MAX_TOKENS_COMPS", ""),
		},

		AWSSQS: SQSConfig{
			PROD_COMPS_SQS: getEnv("PROD_SQS_URL", ""),
			DEV_COMPS_SQS:  getEnv("DEV_SQS_URL", ""),
		},

		Stripe: StripeConfig{
			SecretKey:     getEnv("STRIPE_SECRET_KEY", ""),
			WebhookSecret: getEnv("STRIPE_WEBHOOK_SECRET", ""),
		},

		ServerPort:    getEnvAsInt("PORT", 8080),
		ServerAddress: ":" + getEnv("PORT", ""),

		AppBaseURL: getEnv("APP_BASE_URL", ""),
		SuccessURL: getEnv("SUCCESS_URL", ""),
		CancelURL:  getEnv("CANCEL_URL", ""),

		APIKey:          getEnv("API_KEY", ""),
		ServiceKey:      getEnv("SERVICE_KEY", ""),
		TokenGatewayKey: getEnv("TOKEN_GATEWAY_KEY", ""),

		Cognito: CognitoConfig{
			UserPoolID:   getEnv("COGNITO_USER_POOL_ID", ""),
			ClientID:     getEnv("COGNITO_CLIENT_ID", ""),
			Region:       getEnv("COGNITO_REGION", ""),
			JWKSEndpoint: getEnv("COGNITO_JWKS_ENDPOINT", ""),
		},
	}

	// Load plans
	cfg.Plans.Free = PlanConfig{
		ID:    getEnv("PLAN_FREE_ID", ""),
		Price: getEnv("PLAN_FREE_PRICE", ""),
		Name:  getEnv("PLAN_FREE_NAME", ""),
	}

	cfg.Plans.Scale = PlanConfig{
		ID:    getEnv("PLAN_SCALE_ID", ""),
		Price: getEnv("PLAN_SCALE_PRICE", ""),
		Name:  getEnv("PLAN_SCALE_NAME", ""),
	}

	cfg.Plans.Core = PlanConfig{
		ID:    getEnv("PLAN_CORE_ID", ""),
		Price: getEnv("PLAN_CORE_PRICE", ""),
		Name:  getEnv("PLAN_CORE_NAME", ""),
	}

	cfg.Plans.Team = PlanConfig{
		ID:    getEnv("PLAN_TEAM_ID", ""),
		Price: getEnv("PLAN_TEAM_PRICE", ""),
		Name:  getEnv("PLAN_TEAM_NAME", ""),
	}
	cfg.Plans.Starter = PlanConfig{
		ID:    getEnv("PLAN_STARTER_ID", ""),
		Price: getEnv("PLAN_STARTER_PRICE", ""),
		Name:  getEnv("PLAN_STARTER_NAME", ""),
	}

	// Yearly plans
	cfg.Plans.CoreYearly = PlanConfig{
		ID:    getEnv("PLAN_CORE_YEARLY_ID", ""),
		Price: getEnv("PLAN_CORE_YEARLY_PRICE", ""),
		Name:  getEnv("PLAN_CORE_YEARLY_NAME", ""),
	}

	cfg.Plans.ScaleYearly = PlanConfig{
		ID:    getEnv("PLAN_SCALE_YEARLY_ID", ""),
		Price: getEnv("PLAN_SCALE_YEARLY_PRICE", ""),
		Name:  getEnv("PLAN_SCALE_YEARLY_NAME", ""),
	}
	cfg.Plans.TeamYearly = PlanConfig{
		ID:    getEnv("PLAN_TEAM_YEARLY_ID", ""),
		Price: getEnv("PLAN_TEAM_YEARLY_PRICE", ""),
		Name:  getEnv("PLAN_TEAM_YEARLY_NAME", ""),
	}

	// Token packs
	cfg.Plans.Tokens50 = PlanConfig{
		ID:         getEnv("PLAN_TOKENS_50_ID", ""),
		Price:      getEnv("PLAN_TOKENS_50_PRICE", ""),
		Price_PROD: getEnv("PLAN_TOKENS_50_PRICE_PROD", ""),
		Name:       getEnv("PLAN_TOKENS_50_NAME", ""),
	}

	cfg.Plans.Tokens100 = PlanConfig{
		ID:         getEnv("PLAN_TOKENS_100_ID", ""),
		Price:      getEnv("PLAN_TOKENS_100_PRICE", ""),
		Price_PROD: getEnv("PLAN_TOKENS_100_PRICE_PROD", ""),
		Name:       getEnv("PLAN_TOKENS_100_NAME", ""),
	}
	cfg.Plans.Tokens250 = PlanConfig{
		ID:         getEnv("PLAN_TOKENS_250_ID", ""),
		Price:      getEnv("PLAN_TOKENS_250_PRICE", ""),
		Price_PROD: getEnv("PLAN_TOKENS_250_PRICE_PROD", ""),
		Name:       getEnv("PLAN_TOKENS_250_NAME", ""),
	}
	cfg.Plans.Tokens500 = PlanConfig{
		ID:         getEnv("PLAN_TOKENS_500_ID", ""),
		Price:      getEnv("PLAN_TOKENS_500_PRICE", ""),
		Price_PROD: getEnv("PLAN_TOKENS_500_PRICE_PROD", ""),
		Name:       getEnv("PLAN_TOKENS_500_NAME", ""),
	}

	// Validate required fields

	return cfg, nil
}

func (c *Config) InitializeAuthenticator(zap *zap.Logger) error {
	c.Authenticator = auth.NewCognitoAuthenticator(
		c.Cognito.UserPoolID,
		c.Cognito.ClientID,
		c.Cognito.Region,
		zap,
	)
	return nil
}

func LoadEnv() {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(".env"); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt gets an environment variable as an integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
