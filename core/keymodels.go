package core

import "time"

type KeyUsageData struct {
	BYOKUsage          float64   `json:"byok_usage"`
	BYOKUsageDaily     float64   `json:"byok_usage_daily"`
	BYOKUsageMonthly   float64   `json:"byok_usage_monthly"`
	BYOKUsageWeekly    float64   `json:"byok_usage_weekly"`
	CreatorUserID      string    `json:"creator_user_id"`
	IncludeBYOKInLimit bool      `json:"include_byok_in_limit"`
	IsFreeTier         bool      `json:"is_free_tier"`
	IsManagementKey    bool      `json:"is_management_key"`
	Label              string    `json:"label"`
	Limit              float64   `json:"limit"`
	LimitRemaining     float64   `json:"limit_remaining"`
	LimitReset         string    `json:"limit_reset"`
	Usage              float64   `json:"usage"`
	UsageDaily         float64   `json:"usage_daily"`
	UsageMonthly       float64   `json:"usage_monthly"`
	UsageWeekly        float64   `json:"usage_weekly"`
	IsProvisioningKey  bool      `json:"is_provisioning_key"`
	ExpiresAt          time.Time `json:"expires_at"`
}

type Response struct {
	Data KeyUsageData `json:"data"`
}
