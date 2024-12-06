package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

const (
	apiURL = "https://api.anthropic.com/v1/messages"
)

type Request struct {
	Model        string    `json:"model"`
	Messages     []Message `json:"messages"`
	MaxTokens    int       `json:"max_tokens"`
	SystemPrompt string    `json:"system"`
	Stream       bool      `json:"stream"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type contentBlockDelta struct {
	Type  string `json:"type"`
	Delta struct {
		Text string `json:"text"`
	} `json:"delta"`
}

//Need to get responses that look like this
// 	event: content_block_delta
// data: {"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"Hello"}          }

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
		MaxTokens:    2048,
		SystemPrompt: "You are a CLI assistant program. Please be brief and format your responses so they can be easily read by the user and handled by other CLI programs",
		Stream:       true,
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

	//our new streaming code will need a streaming handler instead of the below

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data:") {
			jsonData := strings.TrimPrefix(line, "data:")
			//fmt.Println(jsonData)
			var delta contentBlockDelta
			err = json.Unmarshal([]byte(jsonData), &delta)
			if err != nil {
				fmt.Println("Error parsing JSON data:", err)
				fmt.Printf("Data: %s\n", string(jsonData))
				os.Exit(1)
			}
			if delta.Type == "content_block_delta" {
				fmt.Printf(delta.Delta.Text)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	/*
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			os.Exit(1)
		}

		   fmt.Printf("Response body: %s\n", string(body))

		   	if resp.StatusCode != http.StatusOK {
		   		fmt.Printf("API request failed with status code %d\n", resp.StatusCode)
		   		fmt.Printf("Response body: %s\n", string(body))
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
	*/
}
