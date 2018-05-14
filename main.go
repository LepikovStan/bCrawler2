package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"

	"time"

	"github.com/LepikovStan/bCrawler2/crawler"
	"github.com/LepikovStan/bCrawler2/manager"
	"github.com/LepikovStan/bCrawler2/parser"
	"github.com/LepikovStan/bCrawler2/queue"
)

const (
	crawlers = 2
	parsers  = 2
)

func Startlist() []string {
	return []string{
		"http://ya.ru",
		"ya.ru",
	}
}

func interceptInterrupt(cr *crawler.Instance, pr *parser.Instance) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c
	fmt.Println()
	fmt.Println("Force Shutdown...")
	cr.Shutdown()
}

func main() {
	start := time.Now()

	crWg := &sync.WaitGroup{}
	prWg := &sync.WaitGroup{}

	crWg.Add(crawlers)
	cr := crawler.New(crawlers, crWg)
	prWg.Add(parsers)
	pr := parser.New(parsers, prWg)
	q := queue.New()
	man := manager.New(cr, pr, q)

	cr.Run()
	pr.Run()
	man.Work(Startlist())
	go interceptInterrupt(cr, pr)

	crWg.Wait()
	fmt.Println("crawler stopped")
	pr.Shutdown()
	prWg.Wait()
	fmt.Println("parser stopped")
	fmt.Println(time.Now().Sub(start))

}
