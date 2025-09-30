package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"
)

const (
	BASE_API_URL      = "https://api.tyrads.com"
	API_KEY           = "YOUR_API_KEY"    // Replace with your actual API key
	API_SECRET        = "YOUR_API_SECRET" // Replace with your actual API secret
	SDK_VERSION       = "3.0"
	SDK_PLATFORM      = "Web"
	LANGUAGE          = "en"
	AGE               = 18                  // Replace with actual age
	GENDER            = 1                   // 1 for male, 2 for female
	PUBLISHER_USER_ID = "PUBLISHER_USER_ID" // Replace with actual publisher user ID
	PORT              = 8080
)

// request payload
type AuthRequest struct {
	PublisherUserID string `json:"publisherUserId"`
	Age             int    `json:"age"`
	Gender          int    `json:"gender"`
}

// response payload from API
type AuthResponse struct {
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s/v%s/auth?lang=%s", BASE_API_URL, SDK_VERSION, LANGUAGE)

	payload := AuthRequest{
		PublisherUserID: PUBLISHER_USER_ID,
		Age:             AGE,
		Gender:          GENDER,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header.Set("X-Api-Key", API_KEY)
	req.Header.Set("X-Api-Secret", API_SECRET)
	req.Header.Set("X-SDK-Version", SDK_VERSION)
	req.Header.Set("X-SDK-Platform", SDK_PLATFORM)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var apiResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	success := apiResp.Data.Token != ""
	token := ""
	if success {
		token = apiResp.Data.Token
	}

	iframeURL := fmt.Sprintf("https://sdk.tyrads.com?token=%s", token)
	iframePremiumURL := fmt.Sprintf("https://sdk.tyrads.com/widget?token=%s", token)

	result := map[string]interface{}{
		"success":          success,
		"token":            token,
		"iframeUrl":        iframeURL,
		"iframePremiumUrl": iframePremiumURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	handlerWithCORS := cors.Default().Handler(mux)

	addr := fmt.Sprintf(":%d", PORT)
	fmt.Printf("Server running at http://localhost%s\n", addr)

	if err := http.ListenAndServe(addr, handlerWithCORS); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
