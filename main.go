package main

import (
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Load environment variables from .env file if present.
	err := godotenv.Load()
	if err != nil {
		log.Warn().Msg("No .env file found. Relying on environment variables.")
	}

	// Configure zerolog to write to the console in a human-friendly format.
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})

	// Retrieve environment variables.
	notionToken := os.Getenv("NOTION_TOKEN")
	notionPageID := os.Getenv("NOTION_PAGE_ID")
	ipServicesEnv := os.Getenv("IP_SERVICES")

	if notionToken == "" || notionPageID == "" {
		log.Fatal().Msg("NOTION_TOKEN and NOTION_PAGE_ID must be set in environment variables.")
	}

	// Parse IP services from environment variable or use defaults.
	var ipServices []IPService
	if ipServicesEnv != "" {
		services := strings.Split(ipServicesEnv, ",")
		for _, s := range services {
			ipServices = append(ipServices, IPService{URL: strings.TrimSpace(s)})
		}
	} else {
		ipServices = []IPService{
			{URL: "https://api.ipify.org"},
			{URL: "https://ipv4.icanhazip.com"},
		}
	}

	// Define a timeout for HTTP requests.
	timeout := 5 * time.Second

	// Fetch the IPv4 address.
	ip, err := GetIPv4Address(ipServices, timeout)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to retrieve external IPv4 address")
	}

	// Log the retrieved IPv4 address.
	log.Info().Str("IPv4", ip).Msg("Retrieved external IPv4 address")

	// Initialize Notion client.
	notionClient := NewNotionClient(notionToken, notionPageID)

	// Update Notion page with the IP address.
	err = notionClient.UpdatePage(ip)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to update Notion page with IP address")
	}

	log.Info().Msg("IP address successfully updated in Notion.")
}
