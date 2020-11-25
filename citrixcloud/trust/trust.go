package trust

/* https://api.was.cloud.com/swagger/ui/index#/ */

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// Clients exported
type Clients struct {
	Principal   string `json:"principal"`
	Locale      string `json:"locale"`
	Subject     string `json:"subject"`
	Token       string `json:"token"`
	OpenIDToken string `json:"openIdToken"`
	ExpiresIn   int    `jason:"expiresIn"`
}

// Reset method
/* This method reset a struct objet */
func (re *Clients) Reset() {
	var zeroA = &Clients{}
	*re = *zeroA
}

// RequestToken exported
func RequestToken(customerID string, clientID string, clientSecret string) (Clients, error) {

	url := "https://trust.citrixworkspacesapi.net/" + customerID + "/tokens/clients"

	postBody, _ := json.Marshal(map[string]string{
		"ClientId":     clientID,
		"ClientSecret": clientSecret,
	})
	responseBody := bytes.NewBuffer(postBody)

	resp, err := http.Post(url, "application/json", responseBody)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	clients1 := Clients{}
	err = json.Unmarshal(body, &clients1)

	return clients1, err
}
