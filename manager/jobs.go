package manager

import "sync"

type JobsCount struct {
	mu                         *sync.RWMutex
	count, errcount, jobscount int
}

func (jl *JobsCount) Inc() {
	jl.mu.Lock()
	defer jl.mu.Unlock()
	jl.count++
}
func (jl *JobsCount) Dec() {
	jl.mu.Lock()
	defer jl.mu.Unlock()
	jl.count--
	if jl.count < 0 {
		jl.count = 0
	}
}
func (jl JobsCount) Empty() bool {
	jl.mu.RLock()
	defer jl.mu.RUnlock()

	return jl.count == 0
}
func (jl *JobsCount) IncErr() {
	jl.mu.Lock()
	defer jl.mu.Unlock()
	jl.errcount++
}
func (jl *JobsCount) DecErr() {
	jl.mu.Lock()
	defer jl.mu.Unlock()
	jl.errcount--
	if jl.errcount < 0 {
		jl.errcount = 0
	}
}
func (jl JobsCount) EmptyErr() bool {
	jl.mu.RLock()
	defer jl.mu.RUnlock()

	return jl.errcount == 0
}
