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

	builder.WriteString(fmt.Sprintf(core.WrapIn("BYOKUsage:", core.Blue)+" %.7f\n", data.BYOKUsage))
	builder.WriteString(fmt.Sprintf(core.WrapIn("BYOKUsageDaily:", core.Blue)+" %.7f\n", data.BYOKUsageDaily))
	builder.WriteString(fmt.Sprintf(core.WrapIn("BYOKUsageMonthly:", core.Blue)+" %.7f\n", data.BYOKUsageMonthly))
	builder.WriteString(fmt.Sprintf(core.WrapIn("BYOKUsageWeekly:", core.Blue)+" %.7f\n", data.BYOKUsageWeekly))
	builder.WriteString(fmt.Sprintf(core.WrapIn("CreatorUserID:", core.Blue)+" %s\n", data.CreatorUserID))
	builder.WriteString(fmt.Sprintf(core.WrapIn("IncludeBYOKInLimit:", core.Blue)+" %t\n", data.IncludeBYOKInLimit))
	builder.WriteString(fmt.Sprintf(core.WrapIn("IsFreeTier:", core.Blue)+" %t\n", data.IsFreeTier))
	builder.WriteString(fmt.Sprintf(core.WrapIn("IsManagementKey:", core.Blue)+" %t\n", data.IsManagementKey))
	builder.WriteString(fmt.Sprintf(core.WrapIn("Label:", core.Blue)+" %s\n", data.Label))
	builder.WriteString(fmt.Sprintf(core.WrapIn("Limit:", core.Blue)+" %.7f\n", data.Limit))
	builder.WriteString(fmt.Sprintf(core.WrapIn("LimitRemaining:", core.Blue)+" %.7f\n", data.LimitRemaining))
	builder.WriteString(fmt.Sprintf(core.WrapIn("LimitReset:", core.Blue)+" %s\n", data.LimitReset))
	builder.WriteString(fmt.Sprintf(core.WrapIn("Usage:", core.Blue)+" %.7f\n", data.Usage))
	builder.WriteString(fmt.Sprintf(core.WrapIn("UsageDaily:", core.Blue)+" %.7f\n", data.UsageDaily))
	builder.WriteString(fmt.Sprintf(core.WrapIn("UsageMonthly:", core.Blue)+" %.7f\n", data.UsageMonthly))
	builder.WriteString(fmt.Sprintf(core.WrapIn("UsageWeekly:", core.Blue)+" %.7f\n", data.UsageWeekly))
	builder.WriteString(fmt.Sprintf(core.WrapIn("IsProvisioningKey:", core.Blue)+" %t\n", data.IsProvisioningKey))
	builder.WriteString(fmt.Sprintf(core.WrapIn("ExpiresAt:", core.Blue)+" %s\n", data.ExpiresAt.Format(time.RFC3339)))

	page = NewPage(core.WrapIn(
		"▗▖ ▗▖ ▗▄▄▖ ▗▄▖  ▗▄▄▖▗▄▄▄▖\n"+
			"▐▌ ▐▌▐▌   ▐▌ ▐▌▐▌   ▐▌   \n"+
			"▐▌ ▐▌ ▝▀▚▖▐▛▀▜▌▐▌▝▜▌▐▛▀▀▘\n"+
			"▝▚▄▞▘▗▄▄▞▘▐▌ ▐▌▝▚▄▞▘▐▙▄▄▖\n", core.Blue), builder.String(), false)
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
