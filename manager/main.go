package manager

import (
	"time"

	"fmt"

	"math/rand"

	"sync"

	"github.com/LepikovStan/bCrawler2/attributes"
	"github.com/LepikovStan/bCrawler2/crawler"
	"github.com/LepikovStan/bCrawler2/parser"
	"github.com/LepikovStan/bCrawler2/queue"
)

const (
	defaultTimeout = time.Second
	errorTimeout   = time.Second * 2
	maxRetry       = 2
	maxDepth       = 2
)

type Manager struct {
	crawler    *crawler.Instance
	parser     *parser.Instance
	qu         *queue.Qu
	attributes attributes.Storage
	jCount     *JobsCount
}

func (m Manager) handleError(id int, url string) {
	go func() {
		m.jCount.IncErr()
		time.Sleep(errorTimeout)
		m.crawler.In <- composeCrawlerJob(id, url)
	}()
}

func (m Manager) dispatch() {
	for {
		select {
		case job := <-m.crawler.Out:
			fmt.Println("crawled", job.Id, job.Error)
			if job == nil {
				continue
			}
			if job.Error != nil {
				m.jCount.DecErr()
				if m.jCount.EmptyErr() && m.attributes.Retry(job.Id) == 0 {
					m.crawler.Shutdown()
					continue
				}
				if m.attributes.Retry(job.Id) == 0 {
					continue
				}
				m.handleError(job.Id, job.Url)
				m.attributes.Insert(job.Id, &attributes.Attr{
					m.attributes.Depth(job.Id),
					m.attributes.Retry(job.Id) - 1,
				})
			}
			if job.Error == nil {
				go m.runParserJob(job)
			}
		case job := <-m.parser.Out:
			m.jCount.Dec()
			nextDepth := m.attributes.Depth(job.Id) + 1
			fmt.Println("parsed", job.Id, m.jCount.jobscount)
			if job.Error != nil {
				continue
			}
			if m.jCount.EmptyErr() && m.jCount.Empty() && nextDepth > maxDepth {
				m.runCrawlerJobs()
				m.crawler.Shutdown()
				continue
			}
			if nextDepth > maxDepth {
				m.jCount.jobscount++
				m.runCrawlerJobs()
				continue
			}
			m.attributes.Delete(job.Id)
			m.jCount.jobscount++
			m.toQueue(job.List, nextDepth, m.attributes.Retry(job.Id))
			m.incJobsCount(len(job.List))
			m.runCrawlerJobs()
		}
	}
}

func (m Manager) incJobsCount(count int) {
	for i := 0; i < count; i++ {
		m.jCount.Inc()
	}
}

func (m Manager) toQueue(list []string, depth, retry int) {
	for i := 0; i < len(list); i++ {
		job := composeCrawlerJob(-1, list[i])
		m.qu.Push(job)
		m.attributes.Insert(job.Id, &attributes.Attr{depth, retry})
	}
}

func handleSendOnClosedChannel() {
	if r := recover(); r != nil {
		if fmt.Sprintf("%s", r) != "send on closed channel" {
			panic(r)
		}
	}
}

func (m Manager) runParserJob(job *crawler.Result) {
	defer handleSendOnClosedChannel()
	m.parser.In <- composeParserJob(job)
}

func (m Manager) runCrawlerJobs() {
	for i := 0; i < m.crawler.Free(); i++ {
		go func() {
			defer handleSendOnClosedChannel()
			msg := m.qu.Pop()
			if msg == nil {
				return
			}
			m.crawler.In <- msg.(*crawler.Job)
		}()
	}
}

func (m Manager) Work(list []string) {
	m.toQueue(list, 0, maxRetry)
	m.runCrawlerJobs()
	go m.dispatch()
}

func New(cr *crawler.Instance, pr *parser.Instance, qu *queue.Qu) *Manager {
	return &Manager{
		cr,
		pr,
		qu,
		attributes.New(),
		&JobsCount{count: 0, mu: &sync.RWMutex{}},
	}
}

func genId() int {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	return r1.Intn(10000)
}

func composeCrawlerJob(id int, url string) *crawler.Job {
	jobId := id
	if jobId == -1 {
		jobId = genId()
	}
	return &crawler.Job{
		Id:             jobId,
		Url:            url,
		Headers:        make(map[string]string),
		RequestTimeout: defaultTimeout,
	}
}
func composeParserJob(result *crawler.Result) *parser.Job {
	return &parser.Job{
		Id:        result.Id,
		Body:      result.Body,
		ParseTag:  "a",
		ParseAttr: "href",
	}
}
