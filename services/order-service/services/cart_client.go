package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"shopsphere/order-service/models"
)

type CartClient struct {
	baseURL    string
	httpClient *http.Client
}

type CartResponse struct {
	Status string      `json:"status"`
	Cart   models.Cart `json:"cart"`
}

func NewCartClient(baseURL string) *CartClient {
	return &CartClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *CartClient) GetCart(
	ctx context.Context,
	token string,
) (*models.Cart, error) {

	url := fmt.Sprintf(
		"%s/api/v1/cart",
		c.baseURL,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to create cart request: %w",
			err,
		)
	}

	req.Header.Set(
		"Authorization",
		"Bearer "+token,
	)

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf(
			"cart service unavailable: %w",
			err,
		)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"cart service returned status %d",
			resp.StatusCode,
		)
	}

	var result CartResponse

	if err := json.NewDecoder(
		resp.Body,
	).Decode(&result); err != nil {
		return nil, fmt.Errorf(
			"failed to decode cart response: %w",
			err,
		)
	}

	return &result.Cart, nil
}
