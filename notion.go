package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
)

// NotionClient handles communication with the Notion API.
type NotionClient struct {
	APIURL           string
	IntegrationToken string
	PageID           string
}

// NewNotionClient initializes a new NotionClient.
func NewNotionClient(token, pageID string) *NotionClient {
	return &NotionClient{
		APIURL:           "https://api.notion.com/v1",
		IntegrationToken: token,
		PageID:           pageID,
	}
}

// UpdatePage updates the specified Notion page with the provided content.
func (c *NotionClient) UpdatePage(content string) error {
	// Prepare the request body according to Notion's API.
	// Here, we'll update the page's title property. Adjust according to your page structure.

	type Text struct {
		Content string `json:"content"`
	}

	type Title struct {
		Text Text `json:"text"`
	}

	type Properties struct {
		Title []Title `json:"Title"`
	}

	type Payload struct {
		Properties Properties `json:"properties"`
	}

	payload := Payload{
		Properties: Properties{
			Title: []Title{
				{
					Text: Text{
						Content: fmt.Sprintf("Current IP: %s", content),
					},
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Create the HTTP request.
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/pages/%s", c.APIURL, c.PageID), bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	// Set necessary headers.
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.IntegrationToken))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Notion-Version", "2022-06-28")

	// Send the request.
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check for successful response.
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		log.Info().Msg("Successfully updated Notion page with the new IP address.")
		return nil
	}

	// Read response body for error details.
	var respBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return fmt.Errorf("failed to decode Notion response: %v", err)
	}

	return fmt.Errorf("failed to update Notion page: %v", respBody)
}
