package sandpile

import (
	"image"
	"image/png"
	"os"
	"sync"
)

type Sandpile struct {
	imageBounds bounds
	config      *Config

	sand             map[point]*sandNode
	queue            *nodeQueue
	findOrCreateNode findOrCreateNodeFunc
}

func NewSandpile(c *Config) *Sandpile {
	s := &Sandpile{
		config: c,
		sand:   make(map[point]*sandNode, 0),
		queue:  newNodeQueue(),
	}
	s.init()
	return s
}

func (s *Sandpile) init() {
	s.findOrCreateNode = s.startLookupWorker()
}

func (s *Sandpile) startLookupWorker() findOrCreateNodeFunc {
	lookupChan := make(chan point)
	outputChan := make(chan *sandNode)
	go func() {
		for {
			select {
			case p := <-lookupChan:
				if node, ok := s.sand[p]; ok {
					outputChan <- node
				} else {
					node = &sandNode{
						sandpile: s,
						loc:      p,
					}
					s.sand[p] = node
					outputChan <- node
				}

			}
		}
	}()

	return func(p point) *sandNode {
		lookupChan <- p
		return <-outputChan
	}
}

func (s *Sandpile) Process() {
	wg := &sync.WaitGroup{}
	node := s.findOrCreateNode(point{0, 0})

	node.Lock()
	node.val = s.config.NumGrains
	node.Unlock()

	s.queue.enqueue(node)

	for i := 0; i < s.config.NumWorkers; i++ {
		wg.Add(1)
		go s.worker(wg)
	}

	wg.Wait()
	s.SaveImage()
}

func (s *Sandpile) worker(wg *sync.WaitGroup) {
	for {
		node, ok := s.queue.dequeue()
		if ok {
			s.queue.enqueue(node.topple()...)
			s.queue.done()
		} else {
			wg.Done()
			return
		}
	}
}

func (s *Sandpile) SaveImage() {

	for p := range s.sand {
		if p.x > s.imageBounds.maxX {
			s.imageBounds.maxX = p.x
		} else if p.x < s.imageBounds.minX {
			s.imageBounds.minX = p.x
		}

		if p.y > s.imageBounds.maxY {
			s.imageBounds.maxY = p.y
		} else if p.y < s.imageBounds.minY {
			s.imageBounds.minY = p.y
		}
	}

	img := image.NewRGBA(
		image.Rect(0, 0, s.imageBounds.maxX*2+1, s.imageBounds.maxY*2+1),
	)

	for p, n := range s.sand {
		img.Set(p.x+s.imageBounds.maxX, p.y+s.imageBounds.maxY, s.config.Colors[n.val])
	}

	f, _ := os.OpenFile(s.config.FileName, os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	png.Encode(f, img)
}
