package ksqldb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// ServerInfo gets the info for your server
func (cl *Client) ServerInfo() (*ServerInfo, error) {
	info := ServerInfoResponse{}
	res, err := http.Get(cl.url + INFO_ENDPOINT)
	if err != nil {
		return nil, fmt.Errorf("can't get server info: %v", err)
	}
	defer res.Body.Close()

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, fmt.Errorf("could not read response body: %v", readErr)
	}

	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("could not parse the response as JSON:\n%w\n%v", err, string(body))
	}

	return &info.KSQLServerInfo, nil
}