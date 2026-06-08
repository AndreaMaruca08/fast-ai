package cli

import (
	"encoding/json"
	"fast_ai_client/core"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func UsagePage(page *Page) *Page {
	ClearTerminal()

	_ = godotenv.Load()

	key := os.Getenv("KEY")
	if key == "" {
		log.Fatal("KEY di openrouter mancante nel .env")
	}

	resp, err := getWithHeader("https://openrouter.ai/api/v1/key", key)
	if err != nil {
		log.Fatal(err)
	}

	var response core.Response
	err = json.Unmarshal(resp, &response)
	if err != nil {
		log.Fatal(err)
	}
	data := response.Data
	builder := strings.Builder{}

	builder.WriteString(fmt.Sprintf("BYOKUsage: %.4f\n", data.BYOKUsage))
	builder.WriteString(fmt.Sprintf("BYOKUsageDaily: %.4f\n", data.BYOKUsageDaily))
	builder.WriteString(fmt.Sprintf("BYOKUsageMonthly: %.f\n", data.BYOKUsageMonthly))
	builder.WriteString(fmt.Sprintf("BYOKUsageWeekly: %.4f\n", data.BYOKUsageWeekly))
	builder.WriteString(fmt.Sprintf("CreatorUserID: %s\n", data.CreatorUserID))
	builder.WriteString(fmt.Sprintf("IncludeBYOKInLimit: %t\n", data.IncludeBYOKInLimit))
	builder.WriteString(fmt.Sprintf("IsFreeTier: %t\n", data.IsFreeTier))
	builder.WriteString(fmt.Sprintf("IsManagementKey: %t\n", data.IsManagementKey))
	builder.WriteString(fmt.Sprintf("Label: %s\n", data.Label))
	builder.WriteString(fmt.Sprintf("Limit: %.4f\n", data.Limit))
	builder.WriteString(fmt.Sprintf("LimitRemaining: %.4f\n", data.LimitRemaining))
	builder.WriteString(fmt.Sprintf("LimitReset: %s\n", data.LimitReset))
	builder.WriteString(fmt.Sprintf("Usage: %.4f\n", data.Usage))
	builder.WriteString(fmt.Sprintf("UsageDaily: %.4f\n", data.UsageDaily))
	builder.WriteString(fmt.Sprintf("UsageMonthly: %.4f\n", data.UsageMonthly))
	builder.WriteString(fmt.Sprintf("UsageWeekly: %.4f\n", data.UsageWeekly))
	builder.WriteString(fmt.Sprintf("IsProvisioningKey: %t\n", data.IsProvisioningKey))
	builder.WriteString(fmt.Sprintf("ExpiresAt: %s\n", data.ExpiresAt.Format(time.RFC3339)))

	page = NewPage("Usage", builder.String(), false)
	page.Update()

	return page
}

func getWithHeader(url string, token string) ([]byte, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, fmt.Errorf("GET failed: status %d: %s", res.StatusCode, string(body))
	}

	return body, nil
}
