/*
This package provides sample Go functions to achieve the following functionality:
•	Retrieve GMM API Key
•	Retrieve GMM Gateway Health Status Summary
•	List Gateway Profiles in GMM
•	List Flexible Templates in GMM
•	Download Gateway Profile from GMM
•	Upload Gateway Profile to GMM
•	Download Flexible Template from GMM
•	Upload Flexible Template to GMM
•	Associate Flexible Template with Gateway Profile in GMM
•	Un-claim a Gateway
•	Claim a Gateway
•	Retrieve Gateway GPS Data for Last Hour
•	Name/Rename Gateway within GMM
•	Modify WiFi SSID and PSK
•	Modify WGB SSID and PSK
*/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Function to retrieve the GMM API Key
// Need to supply GMM username and password
func retrieve_gmm_api_key(email string, password string) string {
	fmt.Println("Retrieving GMM API Key")

	// A Response struct to map the entire response
	type Response struct {
		Access_token string		`json: "access_token"`
		Expires_in int			`json: "expires_in"`
		Token_type string		`json: "token_type"`
	}

	jsonData := map[string]string{"email": email, "password": password}
	jsonValue, _ := json.Marshal(jsonData)
	request, _ := http.NewRequest("POST", "https://us.ciscokinetic.io/api/v2/users/access_token", bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	r, err := client.Do(request)

	if err != nil {
		fmt.Printf("API Token retrieve failed with the error %s\n", err)
		os.Exit(1)
	}

	responseData, _ := ioutil.ReadAll(r.Body)
	var responseObject Response
	e := json.Unmarshal(responseData, &responseObject)
	if e != nil {
		fmt.Println("Unmarshall Error: ", e)
	}

	fmt.Println("GMM API Key:", responseObject.Access_token)
	return responseObject.Access_token
}

// Function to retrieve Gateway Health Summary
// Need to supply GMM API Key and GMM Org ID
func retrieve_gmm_gwy_health_summary(gmm_api_key string, org_id int) {

	type gwy_status struct {
		Summary struct {
			Claiming   int `json:"claiming"`
			Inactive   int `json:"inactive"`
			InProgress int `json:"in_progress"`
			Up         int `json:"up"`
			Down       int `json:"down"`
			Failed     int `json:"failed"`
		} `json:"summary"`
	}

	jsonValue, _ := json.Marshal("")
	request, _ := http.NewRequest("GET", "https://us.ciscokinetic.io/api/v2/organizations/" + strconv.Itoa(org_id) + "/gate_ways", bytes.NewBuffer(jsonValue))
	token := "Token " + gmm_api_key
	request.Header.Set("Authorization", token)
	client := &http.Client{}
	r, err := client.Do(request)

	if err != nil {
		fmt.Printf("Retrieve GMM GWY STATUS error %s\n", err)
		os.Exit(1)
	}

	responseData, _ := ioutil.ReadAll(r.Body)
	var responseObject gwy_status
	e := json.Unmarshal(responseData, &responseObject)
	if e != nil {
		fmt.Println("Unmarshall Error: ", e)
	}

	fmt.Println("")
	fmt.Println("GMM Gateway Health Summary")
	fmt.Println(("--------------------------"))
	fmt.Println("Claiming:", responseObject.Summary.Claiming)
	fmt.Println("Inactive:", responseObject.Summary.Inactive)
	fmt.Println("In_Progress:", responseObject.Summary.InProgress)
	fmt.Println("UP:", responseObject.Summary.Up)
	fmt.Println("DOWN:", responseObject.Summary.Down)
	fmt.Println("FAILED:", responseObject.Summary.Failed)
}

// Function to retrieve a list of gateway profiles in GMM
// Need to supply GMM API Key
func retrieve_gmm_gwy_profiles_list(gmm_api_key string, org_id int) {

	type gwy_profiles struct {
		GatewayProfiles []struct {
			ID                               int           `json:"id"`
			Name                             string        `json:"name"`
		} `json:"gateway_profiles"`
		Paging struct {
			Limit  int `json:"limit"`
			Offset int `json:"offset"`
			Pages  int `json:"pages"`
			Count  int `json:"count"`
			Links  struct {
				First string `json:"first"`
				Last  string `json:"last"`
				Next  string `json:"next"`
			} `json:"links"`
		} `json:"paging"`
	}

	jsonValue, _ := json.Marshal("")
	request, _ := http.NewRequest("GET", "https://us.ciscokinetic.io/api/v2/organizations/" + strconv.Itoa(org_id) + "/gateway_profiles?limit=100", bytes.NewBuffer(jsonValue))
	token := "Token " + gmm_api_key
	request.Header.Set("Authorization", token)
	client := &http.Client{}
	r, err := client.Do(request)

	if err != nil {
		fmt.Printf("Retrieve GMM GWY Profiles error %s\n", err)
		os.Exit(1)
	}

	responseData, _ := ioutil.ReadAll(r.Body)

	var responseObject gwy_profiles
	e := json.Unmarshal(responseData, &responseObject)
	if e != nil {
		fmt.Println("Unmarshall Error: ", e)
	}

	fmt.Println("")
	fmt.Println("Total Number of Gateway Profiles in GMM: ", len(responseObject.GatewayProfiles))
	fmt.Println("")
	fmt.Println("Gateway Profiles in GMM")
	fmt.Println("-----------------------")
	for i := 0; i < len(responseObject.GatewayProfiles); i++ {
		fmt.Println("Profile-ID: ", responseObject.GatewayProfiles[i].ID, " Profile Name: ", responseObject.GatewayProfiles[i].Name)
	}
}

// Function to retrieve a List of Flexible Templates in GMM
// Need to supply GMM API Key and GMM Org ID
func retrieve_gmm_flex_template_list(gmm_api_key string, org_id int) {

	type flex_templates_list struct {
		FlexibleTemplates []struct {
			ID          int      `json:"id"`
			Name        string   `json:"name"`
			Description string   `json:"description"`
			Template    string   `json:"template"`
			Variables   []string `json:"variables"`
		} `json:"flexible_templates"`
		Paging struct {
			Limit  int `json:"limit"`
			Offset int `json:"offset"`
			Pages  int `json:"pages"`
			Count  int `json:"count"`
			Links  struct {
				First string `json:"first"`
				Last  string `json:"last"`
				Prev  string `json:"prev"`
				Next  string `json:"next"`
			} `json:"links"`
		} `json:"paging"`
	}

	jsonValue, _ := json.Marshal("")
	request, _ := http.NewRequest("GET", "https://us.ciscokinetic.io/api/v2/organizations/" + strconv.Itoa(org_id) + "/flexible_templates", bytes.NewBuffer(jsonValue))
	token := "Token " + gmm_api_key
	request.Header.Set("Authorization", token)
	client := &http.Client{}
	r, err := client.Do(request)

	if err != nil {
		fmt.Printf("Retrieve GMM GWY Profiles error %s\n", err)
		os.Exit(1)
	}

	responseData, _ := ioutil.ReadAll(r.Body)

	var responseObject flex_templates_list
	e := json.Unmarshal(responseData, &responseObject)
	if e != nil {
		fmt.Println("Unmarshall Error: ", e)
	}

	fmt.Println("")
	fmt.Println("Total Number of Flexible Templates in GMM: ", len(responseObject.FlexibleTemplates))
	fmt.Println("")
	fmt.Println("Flexible Templates in GMM")
	fmt.Println("-------------------------")
	for i := 0; i < len(responseObject.FlexibleTemplates); i++ {
		fmt.Println("Flex-Template-ID: ", responseObject.FlexibleTemplates[i].ID, " Flex Template Name: ", responseObject.FlexibleTemplates[i].Name)
	}
}

// Function to retrieve GMM Gateway ID correponding to a particular Gateway S/N
// Need to supply GMM API Key and Gateway S/N
func retrieve_gmm_gwy_id(gmm_api_key string, org_id int, gwy_sn string) (id int) {

	type gwy_status struct {
		GateWays []struct {
			ID                      int           `json:"id"`
			UUID                    string        `json:"uuid"`
		} `json:"gate_ways"`
		Paging struct {
			Limit  int `json:"limit"`
			Offset int `json:"offset"`
			Pages  int `json:"pages"`
			Count  int `json:"count"`
			Links  struct {
				First string `json:"first"`
				Last  string `json:"last"`
			} `json:"links"`
		} `json:"paging"`
	}

	jsonValue, _ := json.Marshal("")
	request, _ := http.NewRequest("GET", "https://us.ciscokinetic.io/api/v2/organizations/" + strconv.Itoa(org_id) + "/gate_ways", bytes.NewBuffer(jsonValue))
	token := "Token " + gmm_api_key
	request.Header.Set("Authorization", token)
	client := &http.Client{}
	r, err := client.Do(request)

	if err != nil {
		fmt.Printf("Retrieve GMM GWY STATUS error %s\n", err)
		os.Exit(1)
	}

	responseData, _ := ioutil.ReadAll(r.Body)
	var responseObject gwy_status
	e := json.Unmarshal(responseData, &responseObject)
	if e != nil {
		fmt.Println("Unmarshall Error: ", e)
	}

	gwy_id := 0
	for i:= 0; i < len(responseObject.GateWays); i++ {
		if responseObject.GateWays[i].UUID == gwy_sn {
			gwy_id = responseObject.GateWays[i].ID
		}
	}

	return gwy_id
}

// Function to Unclaim a Gateway from GMM
// Need to supply GMM API Key, GMM Org ID and the Gateway S/N
func gmm_unclaim_gwy(gmm_api_key string, org_id int, gwy_sn string) {

	type unclaim struct {
		Id   int
		UUID string
		Name string
	}

	// Retrieving corresponding gateway ID
	gwy_id := retrieve_gmm_gwy_id(gmm_api_key, org_id, gwy_sn)

	if gwy_id == 0 {
		fmt.Println("")
		fmt.Println("Gateway " + gwy_sn + " could not be unclaimed since it's currently not claimed within GMM")
		return
	}

	jsonValue, _ := json.Marshal("")
	url := "https://us.ciscokinetic.io/api/v2/claims/" + strconv.Itoa(gwy_id)
	request, _ := http.NewRequest("DELETE", url, bytes.NewBuffer(jsonValue))
	token := "Token " + gmm_api_key
	request.Header.Set("Authorization", token)
	client := &http.Client{}
	r, err := client.Do(request)

	if err != nil {
		fmt.Printf("Unclaim Gateway error %s\n", err)
		os.Exit(1)
	}

	responseData, _ := ioutil.ReadAll(r.Body)
	var responseObject unclaim
	e := json.Unmarshal(responseData, &responseObject)
	if e != nil {
		fmt.Println("Unmarshall Error: ", e)
	} else {
		fmt.Println("")
		fmt.Println("Gateway", responseObject.UUID, "with ID", responseObject.Id, "Unclaimed")
	}

	time.Sleep(250000000000)
}

// Function to Name/Re-name a Gateway in GMM
// Need to supply GMM API Key, GMM Org ID, Gateway S/N and the Gateway Name to be configured
func gmm_rename_gwy(gmm_api_key string, org_id int, gwy_sn string, gwy_name string) {

	data := []byte(`{ "gate_way": { "name": "` + gwy_name + `" } }`)

	// Retrieving corresponding gateway ID
	gwy_id := retrieve_gmm_gwy_id(gmm_api_key, org_id, gwy_sn)

	url := "https://us.ciscokinetic.io/api/v2/gate_ways/" + strconv.Itoa(gwy_id)
	request, _ := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	token := "Token " + gmm_api_key
	request.Header.Set("Authorization", token)
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	r, err := client.Do(request)

	if err != nil {
		fmt.Printf("Renaming gateway failed with error %s\n", err)
		os.Exit(1)
	}

	responseData, _ := ioutil.ReadAll(r.Body)

	fmt.Println()
	fmt.Println("Rename Gateway Successful: " + string(responseData))
}

// Function to retrieve Gateway GPS Data for last hour
// Need to supply GMM API Key, GMM Org ID and the Gateway S/N
// Returns Gateway GPS Data as a JSON blob
func retrieve_gmm_gwy_gps(gmm_api_key string, org_id int, gwy_sn string) (gps_history string) {

	jsonValue, _ := json.Marshal("")

	// Retrieving corresponding gateway ID
	gwy_id := retrieve_gmm_gwy_id(gmm_api_key, org_id, gwy_sn)

	now := time.Now()
	from_time := (now.Unix() - 3600) * 1000
	to_time := now.Unix() * 1000
	url := "https://us.ciscokinetic.io/api/v2/gate_ways/" + strconv.Itoa(gwy_id) + "/gps_history?from_time=" + strconv.FormatInt(from_time, 10) + "&to_time=" + strconv.FormatInt(to_time, 10)

	request, _ := http.NewRequest("GET", url, bytes.NewBuffer(jsonValue))
	token := "Token " + gmm_api_key
	request.Header.Set("Authorization", token)
	client := &http.Client{}
	r, err := client.Do(request)

	if err != nil {
		fmt.Printf("Retrieve GMM GWY GPS History error %s\n", err)
		os.Exit(1)
	}

	responseData, _ := ioutil.ReadAll(r.Body)
	return string(responseData)
}

// Function to retrieve a particular Gateway Profile
// Need to supply GMM API Key and Profile Name
// Saves Gateway Profile as JSON file in the /tmp directory
func retrieve_gmm_gwy_profile(gmm_api_key string, org_id int, profile_name string) {

	// Retrieve Gateway Profile ID from GMM
	profile_id := retrieve_gmm_profile_id(gmm_api_key, org_id, profile_name)

	jsonValue, _ := json.Marshal("")
	url := "https://us.ciscokinetic.io/api/v2/gateway_profiles/" + strconv.Itoa(profile_id)
	request, _ := http.NewRequest("GET", url, bytes.NewBuffer(jsonValue))
	token := "Token " + gmm_api_key
	request.Header.Set("Authorization", token)
	client := &http.Client{}
	r, err := client.Do(request)

	if err != nil {
		fmt.Printf("Retrieve GMM GWY Profile error %s\n", err)
		os.Exit(1)
	}

	responseData, _ := ioutil.ReadAll(r.Body)
	fmt.Println(string(responseData))

	filename := "/tmp/" + profile_name + ".json"
	ioutil.WriteFile(filename, responseData, 0644)
}

// Function to Upload a Gateway Profile to GMM
// Need to supply GMM API Key, GMM Org ID, Profile as JSON File
// The Gateway Profile JSON file needs to be in the same directory as this script
func gmm_upload_gwy_profile(gmm_api_key string, org_id int, profile_filename string) {

	jsonFile, err := os.Open(profile_filename)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("")
	fmt.Println("Successfully opened " + profile_filename)

	byteValue, _ := ioutil.ReadAll(jsonFile)

	url := "https://us.ciscokinetic.io/api/v2/organizations/" + strconv.Itoa(org_id) + "/gateway_profiles"
	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(byteValue))
	token := "Token " + gmm_api_key
	request.Header.Set("Authorization", token)
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	r, err := client.Do(request)

	if err != nil {
		fmt.Printf("Upload of Gateway Profile failed with error %s\n", err)
		os.Exit(1)
	}

	responseData, _ := ioutil.ReadAll(r.Body)
	fmt.Println("Gateway Profile Uploaded into GMM: " + string(responseData))
}

// Function to Upload a Flexible Template to GMM
// Need to supply GMM API Key, GMM Org ID, Flexible Template as JSON File
// The Flex Template JSON file needs to be in the same directory as this script
func gmm_upload_flex_template(gmm_api_key string, org_id int, flex_template_filename string) {

	jsonFile, err := os.Open(flex_template_filename)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("")
	fmt.Println("Successfully opened " + flex_template_filename)

	byteValue, _ := ioutil.ReadAll(jsonFile)

	url := "https://us.ciscokinetic.io/api/v2/organizations/" + strconv.Itoa(org_id) + "/flexible_templates"
	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(byteValue))
	token := "Token " + gmm_api_key
	request.Header.Set("Authorization", token)
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	r, err := client.Do(request)

	if err != nil {
		fmt.Printf("Upload of Flexible Template failed with error %s\n", err)
		os.Exit(1)
	}

	responseData, _ := ioutil.ReadAll(r.Body)
	fmt.Println("Flexible Template Uploaded into GMM: " + string(responseData))
}

// Function to retrieve Gateway Profile ID corresponding to Profile Name
// Need to supply GMM API Key, GMM Org ID and the Profile Name
// Returns Profile ID
func retrieve_gmm_profile_id(gmm_api_key string, org_id int, profile_name string) (pid int) {

	type gwy_profiles struct {
		GatewayProfiles []struct {
			ID                               int           `json:"id"`
			Name                             string        `json:"name"`
		} `json:"gateway_profiles"`
		Paging struct {
			Limit  int `json:"limit"`
			Offset int `json:"offset"`
			Pages  int `json:"pages"`
			Count  int `json:"count"`
			Links  struct {
				First string `json:"first"`
				Last  string `json:"last"`
				Next  string `json:"next"`
			} `json:"links"`
		} `json:"paging"`
	}

	jsonValue, _ := json.Marshal("")
	request, _ := http.NewRequest("GET", "https://us.ciscokinetic.io/api/v2/organizations/" + strconv.Itoa(org_id) + "/gateway_profiles?limit=100", bytes.NewBuffer(jsonValue))
	token := "Token " + gmm_api_key
	request.Header.Set("Authorization", token)
	client := &http.Client{}
	r, err := client.Do(request)

	if err != nil {
		fmt.Printf("Retrieve GMM GWY Profiles error %s\n", err)
		os.Exit(1)
	}

	responseData, _ := ioutil.ReadAll(r.Body)
	var responseObject gwy_profiles
	e := json.Unmarshal(responseData, &responseObject)
	if e != nil {
		fmt.Println("Unmarshall Error: ", e)
		os.Exit(1)
	}

	profile_id := 0
	for i := 0; i < len(responseObject.GatewayProfiles); i++ {
		if responseObject.GatewayProfiles[i].Name == profile_name {
			profile_id = responseObject.GatewayProfiles[i].ID
		}
	}

	fmt.Println("")
	fmt.Println("Profile ID for Gateway Profile " + profile_name + " is: " + strconv.Itoa(profile_id))
	return profile_id
}

// Function to retrieve Flexible Template ID
// Need to supply GMM API Key, GMM Org ID and the Flexible Template Name
// Returns the GMM Flexible Template ID
func retrieve_gmm_flex_template_id(gmm_api_key string, org_id int, flex_template_name string) (ftid int) {

	type flex_template_list struct {
		FlexibleTemplates []struct {
			ID          int      `json:"id"`
			Name        string   `json:"name"`
			Description string   `json:"description"`
			Template    string   `json:"template"`
			Variables   []string `json:"variables"`
		} `json:"flexible_templates"`
		Paging struct {
			Limit  int `json:"limit"`
			Offset int `json:"offset"`
			Pages  int `json:"pages"`
			Count  int `json:"count"`
			Links  struct {
				First string `json:"first"`
				Last  string `json:"last"`
				Prev  string `json:"prev"`
				Next  string `json:"next"`
			} `json:"links"`
		} `json:"paging"`
	}

	jsonValue, _ := json.Marshal("")
	request, _ := http.NewRequest("GET", "https://us.ciscokinetic.io/api/v2/organizations/" + strconv.Itoa(org_id) + "/flexible_templates", bytes.NewBuffer(jsonValue))
	token := "Token " + gmm_api_key
	request.Header.Set("Authorization", token)
	client := &http.Client{}
	r, err := client.Do(request)

	if err != nil {
		fmt.Printf("Retrieve GMM GWY Profiles error %s\n", err)
		os.Exit(1)
	}

	responseData, _ := ioutil.ReadAll(r.Body)
	var responseObject flex_template_list
	e := json.Unmarshal(responseData, &responseObject)
	if e != nil {
		fmt.Println("Unmarshall Error: ", e)
		os.Exit(1)
	}

	flex_template_id := 0
	for i := 0; i < len(responseObject.FlexibleTemplates); i++ {
		if responseObject.FlexibleTemplates[i].Name == flex_template_name {
			flex_template_id = responseObject.FlexibleTemplates[i].ID
		}
	}

	fmt.Println("")
	fmt.Println("Flexible Template ID for Flexible Template " + flex_template_name + " is: " + strconv.Itoa(flex_template_id))
	return flex_template_id
}

// Function to retrieve a particular Flexible Template
// Need to supply GMM API Key, GMM Org ID and the Flexible Template Name
// Saves Flexible Template as JSON file in the /tmp directory
func retrieve_gmm_flex_template(gmm_api_key string, org_id int, ft_name string) {

	// Retrieve Flexible Template ID from GMM
	ft_id := retrieve_gmm_flex_template_id(gmm_api_key, org_id, ft_name)

	jsonValue, _ := json.Marshal("")
	url := "https://us.ciscokinetic.io/api/v2/flexible_templates/" + strconv.Itoa(ft_id)
	request, _ := http.NewRequest("GET", url, bytes.NewBuffer(jsonValue))
	token := "Token " + gmm_api_key
	request.Header.Set("Authorization", token)
	client := &http.Client{}
	r, err := client.Do(request)

	if err != nil {
		fmt.Printf("Retrieve GMM Flexible Template error %s\n", err)
		os.Exit(1)
	}

	responseData, _ := ioutil.ReadAll(r.Body)
	fmt.Println(string(responseData))

	filename := "/tmp/" + ft_name + ".json"
	ioutil.WriteFile(filename, responseData, 0644)
}

// Function to Modify WiFi SSID and PSK
// Need to supply GMM API Key, GMM Org ID, Gateway Profile Name, New WiFi SSID and/or New WiFi PSK
func gmm_modify_gwy_wifi(gmm_api_key string, org_id int, profile_name string, wifi_ssid string, wifi_psk string) {

	data := []byte(`{ "wifi_ssid": "` + wifi_ssid + `", "wifi_pre_shared_key": "` + wifi_psk + `" }`)

	profile_id := retrieve_gmm_profile_id(gmm_api_key, org_id, profile_name)

	url := "https://us.ciscokinetic.io/api/v2/gateway_profiles/" + strconv.Itoa(profile_id)
	request, _ := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	token := "Token " + gmm_api_key
	request.Header.Set("Authorization", token)
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	r, err := client.Do(request)

	if err != nil {
		fmt.Printf("Modifying gateway WiFi settings failed with error %s\n", err)
		os.Exit(1)
	}

	responseData, _ := ioutil.ReadAll(r.Body)
	fmt.Println("Modifying Gateway WiFi settings : " + string(responseData))
}

// Function to Modify WGB SSID and PSK
// Need to supply GMM API Key, GMM Org ID, Gateway Profile Name, New WGB SSID and/or New WGB PSK
func gmm_modify_gwy_wgb(gmm_api_key string, org_id int, profile_name string, wgb_ssid string, wgb_psk string) {

	data := []byte(`{ "wgb_ssid": "` + wgb_ssid + `", "wgb_pre_shared_key": "` + wgb_psk + `" }`)

	profile_id := retrieve_gmm_profile_id(gmm_api_key, org_id, profile_name)

	url := "https://us.ciscokinetic.io/api/v2/gateway_profiles/" + strconv.Itoa(profile_id)
	request, _ := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	token := "Token " + gmm_api_key
	request.Header.Set("Authorization", token)
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	r, err := client.Do(request)

	if err != nil {
		fmt.Printf("Modifying gateway WGB settings failed with error %s\n", err)
		os.Exit(1)
	}

	responseData, _ := ioutil.ReadAll(r.Body)
	fmt.Println("Modifying Gateway WGB settings : " + string(responseData))
}

// Function to claim a gateway in GMM
// Need to supply GMM API Key, GMM Org ID, Gateway S/N, Gateway Model and Gateway Profile
func gmm_claim_gwy(gmm_api_key string, org_id int, gwy_sn string, model string, profile_name string) {

	profile_id := retrieve_gmm_profile_id(gmm_api_key, org_id, profile_name)

	payload := `{ "claim_ids": ["` + gwy_sn + `"], "gateway_profile_id": ` + strconv.Itoa(profile_id) + `, "model": "` + model + `" }`

	data := []byte(payload)

	url := "https://us.ciscokinetic.io/api/v2/organizations/" + strconv.Itoa(org_id) + "/claims"
	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
	token := "Token " + gmm_api_key
	request.Header.Set("Authorization", token)
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	r, err := client.Do(request)

	if err != nil {
		fmt.Printf("Claiming Gateway failed with error %s\n", err)
		os.Exit(1)
	}

	responseData, _ := ioutil.ReadAll(r.Body)
	fmt.Println("")
	fmt.Println("Claiming Gateway : " + string(responseData))
}

// Function to Modify WGB SSID and PSK
// Need to supply GMM API Key, GMM Org ID, Gateway Profile Name, New WGB SSID and/or New WGB PSK
func gmm_associate_flex_template(gmm_api_key string, org_id int, profile_name string, flex_template_name string) {

	pid := retrieve_gmm_profile_id(gmm_api_key, org_id, profile_name)
	ftid := retrieve_gmm_flex_template_id(gmm_api_key, org_id, flex_template_name)

	data := []byte(`{"flexible_template_id": ` + strconv.Itoa(ftid) + `, "flexible_template_enable": true, "flexible_template_advanced": false, "flexible_template_variables": [{ "name": "", "value": "none"}]}`)

	url := "https://us.ciscokinetic.io/api/v2/gateway_profiles/" + strconv.Itoa(pid)
	request, _ := http.NewRequest("PUT", url, bytes.NewBuffer(data))
	token := "Token " + gmm_api_key
	request.Header.Set("Authorization", token)
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	r, err := client.Do(request)

	if err != nil {
		fmt.Printf("Associating Flex Template with Base Template failed with error %s\n", err)
		os.Exit(1)
	}

	responseData, _ := ioutil.ReadAll(r.Body)
	fmt.Println("Associated Flexible Template : " + string(responseData))
}

func main() {

	// Retrieve GMM API Key
	gmm_api_key := retrieve_gmm_api_key("uid@enterprise.com", "password")

	// Retrieve Gateway Health Status Summary
	retrieve_gmm_gwy_health_summary(gmm_api_key, 1234)

	// Retrieve the list of Gateway Profiles in GMM
	retrieve_gmm_gwy_profiles_list(gmm_api_key, 1234)

	// Retrieve the list of Flexible Templates in GMM
	retrieve_gmm_flex_template_list(gmm_api_key, 1234)

	// Upload a Gateway Profile to GMM
	gmm_upload_gwy_profile(gmm_api_key,1234,"profile.json")

	// Upload a Flexible Template to GMM
	gmm_upload_flex_template(gmm_api_key,1234,"flex.json")

	// Retrieve a specific Flexible Template
	retrieve_gmm_flex_template(gmm_api_key, 1234, "Add User")

	// Associate a Flexible Template with a Base Template in GMM
	gmm_associate_flex_template(gmm_api_key,1234,"Sample-Template","Add User")

	// Retrieve a specific Gateway Profile from GMM
	retrieve_gmm_gwy_profile(gmm_api_key, 1234, "Sample-Template")

	base_templates := []string { "Basic Mobile Asset - IR829",
		"Basic Remote Asset - IR1101", "Fleet - IR829-2LTE", "Remote Site - IR1101 - Custom Subnet",
	}

	flex_templates := []string { "Adv-QoS 1", "Adv-Firewall 1", "Adv - Static Addressing",
		"Adv - Dynamic Addressing",
	}

	for i := 0; i < len(base_templates); i++ {
		retrieve_gmm_gwy_profile(gmm_api_key, 1234, base_templates[i])
	}

	for j := 0; j < len(flex_templates); j++ {
		retrieve_gmm_flex_template(gmm_api_key, 1234, flex_templates[j])
	}

	for k := 0; k < len (flex_templates); k++ {
		gmm_upload_flex_template(gmm_api_key, 2174, "/tmp/" + flex_templates[k] + ".json")
	}

	for k := 0; k < len (base_templates); k++ {
		gmm_upload_gwy_profile(gmm_api_key, 2174, "/tmp/" + base_templates[k] + ".json")
	}

	// Unclaim a Gateway
	gmm_unclaim_gwy(gmm_api_key, 1234, "ABCDE123456")

	// Claim a Gateway
	gmm_claim_gwy(gmm_api_key, 1234, "ABCDE123456" + "", "IR829", "Sample-Template")

	// Retrieve Gateway GPS Data for the Last Hour
	gps_data := retrieve_gmm_gwy_gps(gmm_api_key, 1234,"ABCDE123456")
	fmt.Println("Gateway FCW22450097 GPS History for the Last Hour: " + gps_data)

	// Name or Rename the Gateway within GMM
	gmm_rename_gwy(gmm_api_key, 1234,"ABCDE123456" , "Cab-123")

	// Modify Gateway WiFi SSID and PSK
	gmm_modify_gwy_wifi(gmm_api_key, 1234, "Sample-Template", "MOD-WiFi-SSID", "MOD-WiFi-PSK")

	// Modify Gateway WGB SSID and PSK
	gmm_modify_gwy_wgb(gmm_api_key, 1234, "Sample-Template", "MOD-WGB-SSID", "MOD-WGB-PSK")
}