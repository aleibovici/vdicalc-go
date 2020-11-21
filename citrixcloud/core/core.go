package core

/* https://api.was.cloud.com/swagger/ui/index#/ */

import (
	"citrixcloud/trust"
	"encoding/json"
	"log"
	"net/http"
)

// ServiceStates exported
type ServiceStates struct {
	Items []items `json:"items"`
}

type items struct {
	ServiceName                string `json:"serviceName"`
	State                      string `json:"state"`
	Type                       string `json:"type"`
	Quantity                   int    `json:"quantity"`
	DaysToExpiration           int    `json:"daysToExpiration"`
	FutureEntitlementStartDate string `json:"futureEntitlementStartDate"`
}

// InventoryStatus exported
type InventoryStatus struct {
	Status            string `json:"status"`
	Activity          string `json:"activity"`
	Device_ip_address string `json:"device_ip_address"`
	Message           string `json:"message"`
	Act_id            string `json:"act_id"`
	Is_last           string `json:"is_last"`
	Starttime         string `json:"starttime"`
	Type              string `json:"type"`
	Id                string `json:"id"`
}

// RequestServiceEntitlement exported
func RequestServiceEntitlement(clients trust.Clients, customerID string) *ServiceStates {

	bearer := "CwsAuth Bearer=" + clients.Token
	url := "https://core.citrixworkspacesapi.net/" + customerID + "/serviceStates"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("Authorization", bearer)
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	result := &ServiceStates{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	return result
}

// RequestInventoryStatus exported
func RequestInventoryStatus(clients trust.Clients, customerID string) *InventoryStatus {

	bearer := "CwsAuth Bearer=" + clients.Token
	url := "https://core.citrixworkspacesapi.net/" + customerID + "/serviceStates"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Add("Authorization", bearer)
	req.Header.Add("isCloud", "true")
	req.Header.Add("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	result := &InventoryStatus{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	return result

}
