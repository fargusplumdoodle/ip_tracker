package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// NotionClient handles communication with the Notion API.
type NotionClient struct {
	APIURL           string
	IntegrationToken string
	PageID           string
	HTTPClient       *http.Client
}

// NewNotionClient initializes a new NotionClient.
func NewNotionClient(token, pageID string) *NotionClient {
	return &NotionClient{
		APIURL:           "https://api.notion.com/v1",
		IntegrationToken: token,
		PageID:           pageID,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetLatestIP retrieves the most recent IP address from the Notion page.
func (nc *NotionClient) GetLatestIP() (string, error) {
	endpoint := nc.APIURL + "/blocks/" + nc.PageID + "/children?page_size=100"

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create request to fetch latest IP")
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+nc.IntegrationToken)
	req.Header.Set("Notion-Version", "2022-06-28")
	req.Header.Set("Content-Type", "application/json")

	resp, err := nc.HTTPClient.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute request to fetch latest IP")
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Error().
			Int("status_code", resp.StatusCode).
			Msg("Non-OK HTTP status when fetching latest IP")
		return "", errors.New("failed to fetch latest IP from Notion")
	}

	var result NotionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Error().Err(err).Msg("Failed to decode Notion response for latest IP")
		return "", err
	}

	// Iterate through the blocks in reverse to find the latest IP entry
	for i := len(result.Results) - 1; i >= 0; i-- {
		block := result.Results[i]
		if block.Type == "paragraph" && block.Paragraph != nil {
			text := getTextFromBlock(block.Paragraph)
			if isValidIPv4(text.IP) {
				return text.IP, nil
			}
		}
	}

	return "", errors.New("no valid IP address found in Notion page")
}

// AppendIP appends a new IP address block with the current timestamp to the Notion page.
func (nc *NotionClient) AppendIP(ip string, timestamp time.Time) error {
	endpoint := nc.APIURL + "/blocks/" + nc.PageID + "/children"

	newBlock := NotionAppendBlock{
		Children: []NotionChild{
			{
				Object: "block",
				Type:   "paragraph",
				Paragraph: &NotionParagraph{
					RichText: []NotionTextObject{
						{
							Type: "text",
							Text: NotionTextContent{
								Content: ip + " (Updated: " + timestamp.Format(time.RFC3339) + ")",
							},
						},
					},
				},
			},
		},
	}

	payload, err := json.Marshal(newBlock)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal new IP block payload")
		return err
	}

	req, err := http.NewRequest("PATCH", endpoint, bytes.NewBuffer(payload))
	if err != nil {
		log.Error().Err(err).Msg("Failed to create request to append new IP")
		return err
	}

	req.Header.Set("Authorization", "Bearer "+nc.IntegrationToken)
	req.Header.Set("Notion-Version", "2022-06-28")
	req.Header.Set("Content-Type", "application/json")

	resp, err := nc.HTTPClient.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to execute request to append new IP")
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Read response body for more detailed error message
		var errorResponse NotionErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err == nil {
			log.Error().
				Int("status_code", resp.StatusCode).
				Str("error", errorResponse.Message).
				Msg("Non-OK HTTP status when appending new IP")
		} else {
			log.Error().
				Int("status_code", resp.StatusCode).
				Msg("Non-OK HTTP status when appending new IP")
		}
		return errors.New("failed to append new IP to Notion page")
	}

	log.Info().Str("ip", ip).Msg("Successfully appended new IP to Notion page")
	return nil
}

// NotionResponse represents the structure of Notion API responses when fetching blocks.
type NotionResponse struct {
	Object  string               `json:"object"`
	Results []NotionFetchedBlock `json:"results"`
}

// NotionFetchedBlock represents a single block fetched from Notion.
type NotionFetchedBlock struct {
	Object    string           `json:"object"`
	ID        string           `json:"id"`
	Type      string           `json:"type"`
	Paragraph *NotionParagraph `json:"paragraph,omitempty"`
}

// NotionParagraph represents a paragraph block in Notion.
type NotionParagraph struct {
	RichText []NotionTextObject `json:"rich_text"`
}

// NotionTextObject represents a text object within a paragraph.
type NotionTextObject struct {
	Type string            `json:"type"`
	Text NotionTextContent `json:"text"`
}

// NotionTextContent represents the content of a text object.
type NotionTextContent struct {
	Content string `json:"content"`
}

// NotionAppendBlock represents the payload structure for appending blocks.
type NotionAppendBlock struct {
	Children []NotionChild `json:"children"`
}

// NotionChild represents the children structure for appending blocks.
type NotionChild struct {
	Object    string           `json:"object"`
	Type      string           `json:"type"`
	Paragraph *NotionParagraph `json:"paragraph,omitempty"`
}

// NotionErrorResponse represents the structure of Notion API error responses.
type NotionErrorResponse struct {
	Object  string `json:"object"`
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// IPText represents the extracted IP and timestamp from a Notion paragraph.
type IPText struct {
	IP string
}

// getTextFromBlock extracts the IP and timestamp from a paragraph block.
func getTextFromBlock(paragraph *NotionParagraph) IPText {
	if len(paragraph.RichText) == 0 {
		return IPText{}
	}
	content := paragraph.RichText[0].Text.Content
	// Assuming the format "IP (Updated: Timestamp)"
	parts := strings.Split(content, " (Updated: ")
	if len(parts) < 1 {
		return IPText{}
	}
	ip := parts[0]
	return IPText{IP: ip}
}
