package ksqldb

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

func (api *KsqldbClient) CloseQuery(ctx context.Context, queryID string) error {
	// Try to close the query
	payload := strings.NewReader(`{"queryId":"` + queryID + `"}`)
	req, err := newCloseQueryRequest(api.http, ctx, payload)

	if err != nil {
		return fmt.Errorf("failed to construct http request to cancel query\n%w", err)
	}

	res, err := api.http.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute http request to cancel query\n%w", err)
	}
	defer res.Body.Close()

	body, err := api.readBody(res.Body)
	if err != nil {
		return fmt.Errorf("can't read response body:\n%w", err)
	}
	// handleError
	if res.StatusCode != http.StatusOK {
		return handleRequestError(res.StatusCode, body)
	}
	return nil
}
