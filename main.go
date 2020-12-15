package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

func main() {
	app := NewApp()
	app.initFlag()
	app.Parse()

	if err := app.Validate(); err != nil {
		log.Fatal(err)
	}
	app.Run()
	app.ShowResult()
}

var errEmptyDomain = errors.New("the domain can not be empty")

type DNSPressureOptions struct {
	// concurrent users, default is 10
	Concurrent int `json:"concurrent"`
	// RESP, number of times to run the test in every user, default is 500
	REPS int `json:"resp"`
	// The total time it takes to execute the test, default is 1*time.Hour
	Time time.Duration `json:"time"`
	// The target domain name that needs to be resolved
	Domain string `json:"domain"`
	//
	OverTime time.Duration `json:"overTime"`
}

type App interface {
	initFlag()
	Parse()
	Validate() error

	Run()
	ShowResult()
}

type app struct {
	flagSets *flag.FlagSet
	options  DNSPressureOptions

	resolver *net.Resolver

	result appResult
}

type appResult struct {
	resultLock         sync.Mutex
	totalCount         int64
	totalOverTimeCount int64
	totalTime          time.Duration
	MaxTime            time.Duration
	MinTime            time.Duration
}

func NewApp() App {
	a := &app{
		resolver: &net.Resolver{},
	}
	flagSets := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	a.flagSets = flagSets

	return a
}

func (a *app) initFlag() {
	a.flagSets.IntVar(&a.options.Concurrent, "concurrent", 10, "concurrent users")
	a.flagSets.IntVar(&a.options.REPS, "resp", 500, "RESP, number of times to run test in every user")
	a.flagSets.DurationVar(&a.options.Time, "time", time.Hour, "The total time it takes to execute the test")
	a.flagSets.DurationVar(&a.options.OverTime, "over-time", 5*time.Second, "over time")
	a.flagSets.StringVar(&a.options.Domain, "domain", "", "The target domain name that needs to be resolved")
}

func (a *app) Parse() {
	// Ignore errors; a.flagSets is set for ExitOnError.
	a.flagSets.Parse(os.Args[1:])
}

func (a *app) Validate() error {
	if a.options.Domain == "" {
		return errEmptyDomain
	}
	return nil
}

func (a *app) Run() {
	ctx, cancel := context.WithTimeout(context.Background(), a.options.Time)
	defer cancel()

	wg := &sync.WaitGroup{}
	for i := 0; i < a.options.Concurrent; i++ {
		wg.Add(1)
		go func() {
			a.Worker(wg, ctx)
		}()
	}
	wg.Wait()
}

func (a *app) Worker(wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	var count int
	for {
		select {
		case <-ctx.Done():
			log.Println("the context timeout expired ")
			return
		default:
			if count >= a.options.REPS {
				return
			}
			startAt := time.Now()
			_, err := a.resolver.LookupIPAddr(context.Background(), a.options.Domain)
			if err != nil {
				return
			}
			endAt := time.Now()
			costTime := endAt.Sub(startAt)

			a.result.resultLock.Lock()
			a.result.totalCount++
			log.Println("number:", a.result.totalCount, fmt.Sprintf("cost time is %v", costTime))
			if costTime > a.result.MaxTime {
				a.result.MaxTime = costTime
			}
			if costTime < a.result.MinTime {
				a.result.MinTime = costTime
			}
			if costTime > a.options.OverTime {
				a.result.totalOverTimeCount++
			}
			a.result.totalTime += costTime
			a.result.resultLock.Unlock()
			count++
		}
	}
}

func (a *app) ShowResult() {
	log.Println("----------------- Result Analyse-----------------")
	log.Println("total count		:\t\t", a.result.totalCount)
	log.Println(fmt.Sprintf("over time(%s)  	:\t\t %d", a.options.OverTime, a.result.totalOverTimeCount))
	log.Println("max time   		:\t\t", a.result.MaxTime)
	log.Println("min time   		:\t\t", a.result.MinTime)
	log.Println("avg time   		:\t\t", fmt.Sprintf("%.2fms", float64(a.result.totalTime.Milliseconds())/float64(a.result.totalCount)))
	log.Println("qps        		:\t\t", fmt.Sprintf("%.2fq/s", float64(a.result.totalCount)/a.result.totalTime.Seconds()))
}
