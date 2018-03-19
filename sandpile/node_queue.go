package sandpile

import (
	"sync"
)

type dequeueFunc func() (*sandNode, bool)
type enqueueFunc func(...*sandNode)

type nodeQueue struct {
	sync.Mutex
	queue   []*sandNode
	inQueue map[*sandNode]bool

	dequeue dequeueFunc
	enqueue enqueueFunc
	done    func()
}

func newNodeQueue() *nodeQueue {
	nq := &nodeQueue{
		queue:   make([]*sandNode, 0, 0),
		inQueue: make(map[*sandNode]bool, 0),
	}

	nq.init()

	return nq
}

func (n *nodeQueue) init() {
	enqueueChan := make(chan []*sandNode)
	enqueueFinished := make(chan struct{})
	n.enqueue = func(nodes ...*sandNode) {
		enqueueChan <- nodes
		<-enqueueFinished
	}

	dequeueChan := make(chan *sandNode)
	n.dequeue = func() (*sandNode, bool) {
		node, ok := <-dequeueChan
		return node, ok
	}

	doneChan := make(chan struct{})
	n.done = func() {
		doneChan <- struct{}{}
	}

	inflight := struct {
		sync.Mutex
		i int
	}{i: 0}
	go func() {
		for {
			<-doneChan
			inflight.Lock()
			inflight.i--
			n.Lock()
			if inflight.i == 0 && len(n.queue) == 0 {
				close(dequeueChan)
			}
			n.Unlock()
			inflight.Unlock()
		}
	}()

	go func() {
		for {
			select {
			case nodes := <-enqueueChan:
				for _, node := range nodes {
					n.Lock()
					if !n.inQueue[node] {
						n.queue = append(n.queue, node)
						n.inQueue[node] = true
					}
					n.Unlock()
				}
				enqueueFinished <- struct{}{}
			default:
			}
			if len(n.queue) > 0 {
				inflight.Lock()
				inflight.i++
				inflight.Unlock()

				n.Lock()
				var node *sandNode
				node, n.queue = n.queue[0], n.queue[1:]
				n.inQueue[node] = false
				n.Unlock()

				dequeueChan <- node
			}
			select {
			case nodes := <-enqueueChan:
				for _, node := range nodes {
					n.Lock()
					if !n.inQueue[node] {
						n.queue = append(n.queue, node)
						n.inQueue[node] = true
					}
					n.Unlock()
				}
				enqueueFinished <- struct{}{}
			}

		}
	}()
}
