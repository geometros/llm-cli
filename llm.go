package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	apiURL = "https://api.anthropic.com/v1/messages"
)

type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	MaxTokens int      `json:"max_tokens"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Response struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a string argument")
		os.Exit(1)
	}

	userInput := os.Args[1]

	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		fmt.Println("ANTHROPIC_API_KEY environment variable is not set")
		os.Exit(1)
	}

	request := Request{
		Model: "claude-3-opus-20240229",
		Messages: []Message{
			{Role: "user", Content: userInput},
		},
		MaxTokens: 2048,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		os.Exit(1)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	client := &http.Client{}
	//fmt.Println("Sending request to Anthropic API...")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// fmt.Printf("Received response with status code: %d\n", resp.StatusCode)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		os.Exit(1)
	}

	//fmt.Printf("Response body: %s\n", string(body))

	if resp.StatusCode != http.StatusOK {
		// fmt.Printf("API request failed with status code %d\n", resp.StatusCode)
		// fmt.Printf("Response body: %s\n", string(body))
		os.Exit(1)
	}

	var apiResponse Response
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		fmt.Println("Error parsing JSON response:", err)
		fmt.Printf("Response body: %s\n", string(body))
		os.Exit(1)
	}

	if len(apiResponse.Content) > 0 {
		// fmt.Println("API Response:")
		fmt.Println(apiResponse.Content[0].Text)
	} else {
		fmt.Println("No content in the response")
		fmt.Printf("Full response: %+v\n", apiResponse)
	}
}
