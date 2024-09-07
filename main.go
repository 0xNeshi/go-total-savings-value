package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type ExchangeRate struct {
	Result             string             `json:"result"`
	Provider           string             `json:"provider"`
	Documentation      string             `json:"documentation"`
	TermsOfUse         string             `json:"terms_of_use"`
	TimeLastUpdateUnix int64              `json:"time_last_update_unix"`
	TimeLastUpdateUTC  string             `json:"time_last_update_utc"`
	TimeNextUpdateUnix int64              `json:"time_next_update_unix"`
	TimeNextUpdateUTC  string             `json:"time_next_update_utc"`
	TimeEolUnix        int64              `json:"time_eol_unix"`
	BaseCode           string             `json:"base_code"`
	Rates              map[string]float64 `json:"rates"`
}

func main() {
	csv_input := flag.String("input", "", "saving input in CSV format: 100,USD,3449,RSD,1010,CAD...")
	flag.Parse()

	inputArr := strings.Split((*csv_input), ",")
	if len(inputArr)%2 != 0 {
		log.Fatalf("Missing currency for the last input amount\n")
	}

	savings := make(map[string]int)
	for i := 0; i < len(inputArr); i += 2 {
		amt, err := strconv.Atoi(inputArr[i])
		if err != nil {
			log.Fatalf("Item at index %d ('%s') is not a valid number\n", i, inputArr[i])
		}

		curCode := inputArr[i+1]
		if len(curCode) != 3 {
			log.Fatalf("Currency code at %d should be 3 characters in length: '%s'\n", i+1, curCode)
		}
		savings[curCode] += amt
	}

	// Send an HTTP GET request
	resp, err := http.Get("https://open.er-api.com/v6/latest/USD")
	if err != nil {
		log.Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}

	// Parse the JSON into the struct
	var data ExchangeRate
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	totalUsd := 0.0
	for curCode, amount := range savings {
		rate, exists := data.Rates[curCode]
		if !exists {
			log.Fatalf("Non-existent currency code: %v", curCode)
		}
		totalUsd += float64(amount) / rate
	}

	rsdRate, exists := data.Rates["RSD"]
	if !exists {
		log.Fatalf("Unexpected: RSD is not listed as a currency")
	}
	total := totalUsd * rsdRate

	fmt.Printf("Total: %0.2f\n", total)
}
