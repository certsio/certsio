package resolver

import "sync"

const (
	Alive ResultType = iota + 1
	Error
)

// Pool is a pool of resolvers to resolve certificate names.
type Pool struct {
	*Resolver
	Tasks   chan HostEntry
	Results chan Result
	wg      *sync.WaitGroup
}

// NewPool creates a pool of resolvers for resolving certificates.
func (r *Resolver) NewPool(workers int) *Pool {
	resolutionPool := &Pool{
		Resolver: r,
		Tasks:    make(chan HostEntry),
		Results:  make(chan Result),
		wg:       &sync.WaitGroup{},
	}

	go func() {
		for i := 0; i < workers; i++ {
			resolutionPool.wg.Add(1)
			go resolutionPool.resolve()
		}
		resolutionPool.wg.Wait()
		close(resolutionPool.Results)
	}()

	return resolutionPool
}

// resolve resolves a hostname to IP addresses.
func (p *Pool) resolve() {
	for task := range p.Tasks {
		hosts, err := p.client.Lookup(task.Host)
		if err != nil {
			p.Results <- Result{Type: Error, Task: task, Error: err}
			continue
		}

		if len(hosts) == 0 {
			continue
		}

		p.Results <- Result{Type: Alive, Task: task, IPs: hosts}
	}
	p.wg.Done()
}
