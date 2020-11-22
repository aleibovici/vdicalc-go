package was

/* https://api.was.cloud.com/swagger/ui/index# */

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"vdicalc/citrixcloud/trust"
	"vdicalc/functions"
)

// UserExperienceTrendData exported
type UserExperienceTrendData struct {
	UserExperienceReport userExperienceAggregateData `json:"userExperienceReport"`
}

type userExperienceAggregateData struct {
	Aggregated userExperienceSummary          `json:"aggregated"`
	Items      []userExperienceTimeSeriesData `json:"items"`
}

type userExperienceSummary struct {
	TotalUsers int                              `json:"totalUsers"`
	Items      []userExperienceAggregateDetails `json:"items"`
}

type userExperienceTimeSeriesData struct {
	DateTime  int64 `json:"dateTime"`
	Excellent int   `json:"excellent"`
	Fair      int   `json:"fair"`
	Poor      int   `json:"poor"`
}

type userExperienceAggregateDetails struct {
	ChangeCount int    `json:"changeCount"`
	Category    string `json:"category"`
	Value       int    `json:"value"`
}

// RequestUserExperienceTrend export
func RequestUserExperienceTrend(clients trust.Clients, customerID string, interval int, startTime int64, endTime int64, timerange string, siteID string, isUpdate bool, isDev bool) userExperienceSummary {

	bearer := "CwsAuth bearer=" + clients.Token
	baseURL, err := url.Parse("https://api-b.was.cloud.com/wsanalytics/api/v1/" + customerID + "/userexperiencetrend")
	if isDev == true {
		baseURL, err = url.Parse("https://api-b.was.cloudnacho.com/wsanalytics/api/v1/" + customerID + "/userexperiencetrend")
	}

	params := url.Values{}
	params.Add("interval", functions.InttoStr(interval))
	params.Add("endTime", functions.Int64toStr(endTime))
	params.Add("startTime", functions.Int64toStr(startTime))
	params.Add("range", timerange)
	params.Add("siteId", siteID)
	params.Add("interval", strconv.FormatBool(isUpdate))

	// Add Query Parameters to the URL
	baseURL.RawQuery = params.Encode() // Escape Query Parameters

	req, err := http.NewRequest(http.MethodGet, baseURL.String(), nil)
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	result := &UserExperienceTrendData{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	return result.UserExperienceReport.Aggregated
}

// CalculateScore export
/* This public fundtion is reponsible for calculating the score indicator
Each category is baselined, dividing by the total number of users/machines/vdas
Excelent = 10 / Fair = 5 / Poor = 1
Multiply by 100 for bigger number and effect */
func CalculateScore(totalUsers int, excelent int, fair int, poor int) int {

	if totalUsers != 0 {
		var excelentBaseline float64 = (float64(excelent) / float64(totalUsers)) * 10
		var fairBaseline float64 = (float64(fair) / float64(totalUsers)) * 5
		var poorBaseline float64 = (float64(poor) / float64(totalUsers)) * 1
		return int((excelentBaseline + fairBaseline + poorBaseline) * 100)
	}
	return 0
}
