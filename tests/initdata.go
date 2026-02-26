package main

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Mock user data similar to the Python version
var mockUserData = map[string]any{
	"id":                 687627953,
	"first_name":         "linuxfight",
	"last_name":          "",
	"username":           "linuxfight",
	"language_code":      "ru",
	"is_premium":         false,
	"allows_write_to_pm": false,
	"photo_url":          "https://example.com/cat.jpeg",
}

// loadBotToken reads BOT_TOKEN from dev.env (same format as the Python script)
func loadBotToken() (string, error) {
	file, err := os.Open("dev.env")
	if err != nil {
		return "", fmt.Errorf("dev.env not found: %w", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if after, ok := strings.CutPrefix(line, "BOT_TOKEN="); ok {
			token := after
			token = strings.TrimSpace(token)
			if token == "" {
				return "", fmt.Errorf("BOT_TOKEN is empty in dev.env")
			}
			return token, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("BOT_TOKEN not found in dev.env")
}

// createUserJSON removes empty string values (but keeps false booleans) and returns a compact JSON string.
func createUserJSON(userData map[string]any) (string, error) {
	filtered := make(map[string]any)
	for k, v := range userData {
		// Keep false booleans, remove empty strings, nil values are not expected here.
		if str, ok := v.(string); ok && str == "" {
			continue
		}
		filtered[k] = v
	}
	// Marshal without indentation and with compact separators.
	b, err := json.Marshal(filtered)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// generateHash replicates the algorithm from the Python script.
func generateHash(data map[string]string, botToken string) string {
	// Step 1: create key=value pairs excluding "hash"
	pairs := make([]string, 0, len(data))
	for k, v := range data {
		if k == "hash" {
			continue
		}
		pairs = append(pairs, fmt.Sprintf("%s=%s", k, v))
	}
	// Step 2: sort alphabetically
	sort.Strings(pairs)
	// Step 3: HMAC‑SHA256 with key "WebAppData" and bot token
	secretKey := hmac.New(sha256.New, []byte(botToken))
	secretKey.Write([]byte("WebAppData"))
	secret := secretKey.Sum(nil)
	// Step 4: HMAC‑SHA256 with secret key and joined pairs
	dataString := strings.Join(pairs, "\n")
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(dataString))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// generateMockInitData builds a mock init data string similar to the Python version.
func generateMockInitData(botToken string, userData map[string]any, authDate int64) (string, error) {
	if userData == nil {
		userData = mockUserData
	}
	if authDate == 0 {
		authDate = time.Now().Unix()
	}
	userJSON, err := createUserJSON(userData)
	if err != nil {
		return "", err
	}
	// Prepare init data without hash
	initData := map[string]string{
		"user":      userJSON,
		"chat_type": "private",
		"auth_date": strconv.FormatInt(authDate, 10),
	}
	// Generate hash
	hash := generateHash(initData, botToken)
	initData["hash"] = hash
	// Build URL‑encoded query string
	values := url.Values{}
	for k, v := range initData {
		values.Set(k, v)
	}
	return values.Encode(), nil
}

// parseAndDisplay prints the init data in a human‑readable form.
func parseAndDisplay(initDataString string) {
	parsed, _ := url.ParseQuery(initDataString)
	fmt.Println("=== Parsed Init Data ===")
	for key, vals := range parsed {
		value := vals[0]
		if key == "user" {
			var user map[string]any
			if err := json.Unmarshal([]byte(value), &user); err == nil {
				fmt.Printf("%s:\n", key)
				for uk, uv := range user {
					fmt.Printf("  %s: %v\n", uk, uv)
				}
				continue
			}
		}
		fmt.Printf("%s: %s\n", key, value)
	}
}

func main() {
	botToken, err := loadBotToken()
	if err != nil {
		fmt.Printf("❌ Error loading bot token: %v\n", err)
		fmt.Println("\n💡 Create a dev.env file with: BOT_TOKEN=your_telegram_bot_token_here")
		return
	}
	fmt.Println("✅ Bot token loaded from dev.env")
	mockInit, err := generateMockInitData(botToken, nil, 0)
	if err != nil {
		fmt.Printf("❌ Error generating mock init data: %v\n", err)
		return
	}
	fmt.Println("\n🎯 Generated Init Data:")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println(mockInit)
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("\n📊 Parsed Data:")
	fmt.Println(strings.Repeat("-", 50))
	parseAndDisplay(mockInit)
}
