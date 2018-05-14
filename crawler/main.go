package crawler

import (
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type Job struct {
	Id             int
	Url            string
	Headers        map[string]string
	RequestTimeout time.Duration
}

type Result struct {
	Id    int
	Url   string
	Body  []byte
	Error error
}

type Instance struct {
	concurrency int
	In          chan *Job
	Out         chan *Result
	wg          *sync.WaitGroup
}

func New(concurrency int, wg *sync.WaitGroup) *Instance {
	return &Instance{
		In:          make(chan *Job, concurrency),
		Out:         make(chan *Result, concurrency),
		concurrency: concurrency,
		wg:          wg,
	}
}

func (ins *Instance) Free() int {
	return ins.concurrency - len(ins.In)
}

func (ins *Instance) Shutdown() {
	close(ins.In)
}

func (ins *Instance) Run() {
	for i := 0; i < ins.concurrency; i++ {
		go func() {
			client := createHttpClient()
			for job := range ins.In {
				result, err := crawl(client, job)
				ins.Out <- &Result{job.Id, job.Url, result, err}
			}
			ins.wg.Done()
		}()
	}
}

func createHttpClient() *http.Client {
	return &http.Client{}
}

func createRequest(job *Job) *http.Request {
	req, err := http.NewRequest("GET", job.Url, nil)
	if err != nil {
		log.Fatal(err)
	}
	for headerName, headerValue := range job.Headers {
		req.Header.Add(headerName, headerValue)
	}
	return req
}

func crawl(client *http.Client, job *Job) ([]byte, error) {
	req := createRequest(job)
	client.Timeout = job.RequestTimeout
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return result, nil
}
