package api_client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func GetTargets() (*Targets, error) {
	req, err := http.NewRequest("GET", "https://itarmy.pp.ua/api/?type=all", nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	targets := &Targets{}
	err = json.Unmarshal(body, targets)
	if err != nil {
		return nil, err
	}
	return targets, nil
}
