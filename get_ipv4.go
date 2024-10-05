package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// IPService represents an external service that returns the IPv4 address.
type IPService struct {
	URL string
}

// GetIPv4Address attempts to retrieve the external IPv4 address using the provided services.
// It tries each service in order and returns the first successful result.
func GetIPv4Address(services []IPService, timeout time.Duration) (string, error) {
	client := &http.Client{
		Timeout: timeout,
	}

	for _, service := range services {
		log.Info().Str("service", service.URL).Msg("Attempting to fetch IPv4 address")

		resp, err := client.Get(service.URL)
		if err != nil {
			log.Error().Str("service", service.URL).Err(err).Msg("Failed to make request")
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Error().
				Str("service", service.URL).
				Int("status_code", resp.StatusCode).
				Msg("Non-OK HTTP status")
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error().Str("service", service.URL).Err(err).Msg("Failed to read response body")
			continue
		}

		ip := strings.TrimSpace(string(body))
		if isValidIPv4(ip) {
			log.Info().Str("service", service.URL).Str("ip", ip).Msg("Successfully retrieved IPv4 address")
			return ip, nil
		}

		log.Error().
			Str("service", service.URL).
			Str("ip", ip).
			Msg("Invalid IPv4 address format")
	}

	return "", errors.New("unable to retrieve external IPv4 address from all services")
}

// isValidIPv4 performs a basic validation of the IPv4 address format.
// For more robust validation, consider using regex or the net package.
func isValidIPv4(ip string) bool {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return false
	}
	for _, p := range parts {
		if len(p) == 0 {
			return false
		}
		// Further validation can be added here (e.g., numerical range)
	}
	return true
}

