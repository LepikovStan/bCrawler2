package parser

import (
	"bytes"
	"sync"

	"golang.org/x/net/html"
)

type Job struct {
	Id        int
	Body      []byte
	ParseTag  string
	ParseAttr string
}
type Result struct {
	Id    int
	List  []string
	Error error
}

type Instance struct {
	In          chan *Job
	Out         chan *Result
	concurrency int
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

func parse(n *html.Node, parseTag, parseAttr string) []string {
	result := make([]string, 0)
	if n.Type == html.ElementNode && n.Data == parseTag {
		for i := 0; i < len(n.Attr); i++ {
			if n.Attr[i].Key == parseAttr {
				result = append(result, n.Attr[i].Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result = append(result, parse(c, parseTag, parseAttr)...)
	}
	return result
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
			for job := range ins.In {
				doc, err := html.Parse(bytes.NewReader(job.Body))
				ins.Out <- &Result{
					Id:    job.Id,
					List:  parse(doc, job.ParseTag, job.ParseAttr),
					Error: err,
				}
			}
			ins.wg.Done()
		}()
	}
}
