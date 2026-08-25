package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type InventoryClient struct {
	baseURL    string
	httpClient *http.Client
}

type InventoryResponse struct {
	ID               int64 `json:"id"`
	ProductID        int64 `json:"product_id"`
	Quantity         int   `json:"quantity"`
	ReservedQuantity int   `json:"reserved_quantity"`
	ReorderLevel     int   `json:"reorder_level"`
}

type InventoryAPIResponse struct {
	Status    string            `json:"status"`
	Inventory InventoryResponse `json:"inventory"`
}

func NewInventoryClient(baseURL string) *InventoryClient {
	return &InventoryClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *InventoryClient) GetInventory(
	ctx context.Context,
	productID int64,
) (*InventoryResponse, error) {

	url := fmt.Sprintf(
		"%s/api/v1/inventory/%d",
		c.baseURL,
		productID,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to create inventory request: %w",
			err,
		)
	}

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf(
			"inventory service unavailable: %w",
			err,
		)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"inventory service returned status %d",
			resp.StatusCode,
		)
	}

	var result InventoryAPIResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf(
			"failed to decode inventory response: %w",
			err,
		)
	}

	return &result.Inventory, nil
}

func (c *InventoryClient) HasEnoughStock(
	ctx context.Context,
	productID int64,
	requiredQuantity int,
) (bool, error) {

	inventory, err := c.GetInventory(
		ctx,
		productID,
	)

	if err != nil {
		return false, err
	}

	available := inventory.Quantity -
		inventory.ReservedQuantity

	return available >= requiredQuantity, nil
}
