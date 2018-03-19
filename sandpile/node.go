package sandpile

import (
	"sync"
)

type sandNode struct {
	sync.Mutex

	val int
	loc point

	Top    *sandNode
	Right  *sandNode
	Bottom *sandNode
	Left   *sandNode

	sandpile *Sandpile
}

func (n *sandNode) topple() (nodes []*sandNode) {
	n.Lock()

	n.checkNodes()

	nodes = make([]*sandNode, 0, 5)

	if n.val > n.sandpile.config.MaxGrains {
		n.val -= 4
		n.Unlock()

		if n.val > n.sandpile.config.MaxGrains {
			nodes = append(nodes, n)
		}

		n.Top.Lock()
		n.Top.val++
		if n.Top.val > n.sandpile.config.MaxGrains {
			nodes = append(nodes, n.Top)
		}
		n.Top.Unlock()

		n.Right.Lock()
		n.Right.val++
		if n.Right.val > n.sandpile.config.MaxGrains {
			nodes = append(nodes, n.Right)
		}
		n.Right.Unlock()

		n.Bottom.Lock()
		n.Bottom.val++
		if n.Bottom.val > n.sandpile.config.MaxGrains {
			nodes = append(nodes, n.Bottom)
		}
		n.Bottom.Unlock()

		n.Left.Lock()
		n.Left.val++
		if n.Left.val > n.sandpile.config.MaxGrains {
			nodes = append(nodes, n.Left)
		}
		n.Left.Unlock()
	} else {
		n.Unlock()
	}
	return
}

func (n *sandNode) checkNodes() {
	if n.Top == nil {
		n.Top = n.sandpile.findOrCreateNode(point{n.loc.x, n.loc.y + 1})
		n.Top.Lock()
		n.Top.Bottom = n
		n.Top.Unlock()
	}

	if n.Right == nil {
		n.Right = n.sandpile.findOrCreateNode(point{n.loc.x + 1, n.loc.y})
		n.Right.Lock()
		n.Right.Left = n
		n.Right.Unlock()
	}

	if n.Bottom == nil {
		n.Bottom = n.sandpile.findOrCreateNode(point{n.loc.x, n.loc.y - 1})
		n.Bottom.Lock()
		n.Bottom.Top = n
		n.Bottom.Unlock()
	}

	if n.Left == nil {
		n.Left = n.sandpile.findOrCreateNode(point{n.loc.x - 1, n.loc.y})
		n.Left.Lock()
		n.Left.Right = n
		n.Left.Unlock()
	}

}
