package api_client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func GetTargetsItArmyPpUa() (*TargetsItArmy, error) {
	req, err := http.NewRequest("GET", "https://itarmy.pp.ua/api/?type=all", nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t, fErr := targetsFromFile()
		if fErr != nil {
			return nil, err
		}
		fmt.Println("Could not get targets from api, using default targets from file...")
		return t, nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	targets := &TargetsItArmy{}
	err = json.Unmarshal(body, targets)
	if err != nil {
		return nil, err
	}
	return targets, nil
}

func GetTargets() (*TargetsItArmy, error) {
	req, err := http.NewRequest("GET", "http://164.92.247.88:9300/victims", nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t, fErr := targetsFromFile()
		if fErr != nil {
			return nil, err
		}
		fmt.Println("Could not get targets from api, using default targets from file...")
		return t, nil
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
	return sortTargets(targets), nil
}

func sortTargets(t *Targets) *TargetsItArmy {
	// TODO sort on more criteria
	onlineTargets := []*Statuses{}
	targets := []string{}
	for i := 0; i < len(t.Statuses); i++ {
		if strings.Contains(t.Statuses[i].Status, "UP") {
			onlineTargets = append(onlineTargets, &t.Statuses[i])
		}
	}
	if len(onlineTargets) >= 8 {
		for i := 0; i < 8; i++ {
			if onlineTargets[i].URL != "" {
				targets = append(targets, onlineTargets[i].URL)
			}
		}
	} else {
		for i := 0; i < len(onlineTargets); i++ {
			if onlineTargets[i].URL != "" {
				targets = append(targets, onlineTargets[i].URL)
			}
		}
	}
	return &TargetsItArmy{
		Online:  targets,
		Offline: nil,
	}
}

func targetsFromFile() (*TargetsItArmy, error) {
	f, err := os.Open("default_targets.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	byteValue, _ := ioutil.ReadAll(f)
	t := TargetsItArmy{}
	err = json.Unmarshal(byteValue, &t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
