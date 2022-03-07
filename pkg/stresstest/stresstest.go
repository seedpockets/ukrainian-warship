package stresstest

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Warship struct {
	Target         string
	Workers        int
	HttpClient     *http.Client
	Debug          bool
	Stop           chan bool
	SuccessRequest int64
	AmountRequests int64
	sync.Mutex
}

func New(target string, debug bool, workers int) (*Warship, error) {
	u, err := url.Parse(target)
	if err != nil || len(u.Host) == 0 {
		return nil, fmt.Errorf("Undefined host or error = %v", err)
	}
	w := &Warship{
		Target:  u.String(),
		Workers: workers,
		HttpClient: &http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       time.Millisecond * 500,
		},
		Debug:          debug,
		SuccessRequest: 0,
		AmountRequests: 0,
	}
	return w, nil
}

func (w Warship) String() string {
	if w.Workers == 0 {
		return fmt.Sprintf("Total request: %d\t\tSuccessful request: %d\t\tTarget: %s", w.AmountRequests, w.SuccessRequest, w.Target)
	}
	return fmt.Sprintf("%d\t\t%d\t\t%s", w.AmountRequests, w.SuccessRequest, w.Target)
}

func (w *Warship) Fire() {
	w.Stop = make(chan bool)
	for i := 0; i < w.Workers; i++ {
		go func() {
			for {
				select {
				case stop := <-w.Stop:
					if stop {
						fmt.Println("Stopping worker for target: ", w.Target)
						return
					}
				default:
					w.Lock()
					err := w.send()
					if err == nil {
						w.AmountRequests += 1
						w.SuccessRequest += 1
					}
					w.AmountRequests += 1
					w.Unlock()
				}
			}
		}()
	}
}

func (w *Warship) FocusFire() {
	for {
		err := w.send()
		if err == nil {
			w.AmountRequests++
		}
		w.AmountRequests++
		w.SuccessRequest++
		if !w.Debug {
			fmt.Println("\x1b[2J")
			fmt.Println(w)
		}
	}
}

func (w *Warship) send() error {
	req, err := http.NewRequest("GET", w.Target, nil)
	if err != nil {
		if w.Debug {
			fmt.Println(err.Error())
		}
		return err
	}
	w.addHeaders(req)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if w.Debug {
			fmt.Println(err.Error())
		}
		return err
	}
	if w.Debug {
		fmt.Println("Response status: ", resp.Status)
		io.Copy(ioutil.Discard, resp.Body)
		defer resp.Body.Close()
	} else {
		io.Copy(ioutil.Discard, resp.Body)
		defer resp.Body.Close()
	}
	return nil
}

func (w *Warship) addHeaders(r *http.Request) {
	// r.Header.Add("User-Agent", getUserAgent())
	r.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.2) AppleWebKit/535.7 (KHTML, like Gecko) Comodo_Dragon/16.1.1.0 Chrome/16.0.912.63 Safari/535.7")
	r.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	r.Header.Add("Accept-Language", "en-us,en;q=0.5")
	r.Header.Add("Accept-Encoding", "gzip,deflate")
	r.Header.Add("Accept-Charset", "ISO-8859-1,utf-8;q=0.7,*;q=0.7")
	r.Header.Add("Keep-Alive", "115")
	r.Header.Add("Connection", "keep-alive")
}

func getUserAgent() string {
	rand.Seed(time.Now().UnixNano())
	userAgents := []string{
		"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.0) Opera 12.14",
		"Mozilla/5.0 (X11; Ubuntu; Linux i686; rv:26.0) Gecko/20100101 Firefox/26.0",
		"Mozilla/5.0 (X11; U; Linux x86_64; en-US; rv:1.9.1.3) Gecko/20090913 Firefox/3.5.3",
		"Mozilla/5.0 (Windows; U; Windows NT 6.1; en; rv:1.9.1.3) Gecko/20090824 Firefox/3.5.3 (.NET CLR 3.5.30729)",
		"Mozilla/5.0 (Windows NT 6.2) AppleWebKit/535.7 (KHTML, like Gecko) Comodo_Dragon/16.1.1.0 Chrome/16.0.912.63 Safari/535.7",
		"Mozilla/5.0 (Windows; U; Windows NT 5.2; en-US; rv:1.9.1.3) Gecko/20090824 Firefox/3.5.3 (.NET CLR 3.5.30729)",
		"Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US; rv:1.9.1.1) Gecko/20090718 Firefox/3.5.1",
	}
	return userAgents[rand.Intn(len(userAgents))]
}
