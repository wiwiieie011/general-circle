package api_http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"ticket-service/internal/dto"
	dto_api "ticket-service/internal/dto/api"
)

type EventClient struct {
	baseURL string
	client  *http.Client
}

func NewEventClient(baseURL string) *EventClient {
	return &EventClient{
		baseURL: baseURL,
		client: &http.Client{},
	}
}

func (c *EventClient) GetEvent(ctx context.Context, eventId uint64) (*dto_api.EventResponse, error) {
	url := fmt.Sprintf("%s/event/%d", c.baseURL, eventId)

	req, err := http.NewRequestWithContext(
		ctx, http.MethodGet, url, nil,
	)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, dto.ErrEventNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected status: %d, got: %d", http.StatusOK, resp.StatusCode)
	}

	var respDto dto_api.EventResponse
	if err := json.NewDecoder(resp.Body).Decode(&respDto); err != nil {
		return nil, err
	}

	return &respDto, nil
}
