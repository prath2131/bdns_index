package index

import (
	//"crypto/sha256"
	//"fmt"
	//"strconv"
	"github.com/bleasey/bdns/internal/blockchain"
)

type ARTTree struct {
	root *ARTNode
}

func (t *ARTTree) Add(key []byte, value *blockchain.Transaction) {
	l := &Leaf{wholekey: key, value: value}
	lnode := l.makenode()
	t.root.Insert(key, 0, lnode)

}

func (t *ARTTree) Remove(key []byte) {
	t.root, _ = Delete(t.root, key, 0)

}

func (t *ARTTree) Update(key []byte, newValue *blockchain.Transaction) {
	t.root, _ = Delete(t.root, key, 0)
	l := &Leaf{wholekey: key, value: newValue}
	lnode := l.makenode()
	t.root.Insert(key, 0, lnode)
}

func (t *ARTTree) Search(key []byte) (txn *blockchain.Transaction, found bool) {
	return t.root.search(key, 0)

}

/*func (t *ARTTree) DisplayInOrder() {
	t.root.displayNodesInOrder()
}*/

type ARTNode struct {
	nodetype   int                              // Node type: 4-Node4, 16-Node16, 48-Node48, 256-Node256, 1-leaf
	ChildCount int                              // Number of children
	isEnd      bool                             // Is this node a leaf
	Edge       []byte                           // edge pointing this node and whole key if leaf node
	EdgeLen    int                              // Length of the common prefix
	data       interface{ makenode() *ARTNode } // could point to one of Node4, Node16, Node48, Node256, leaf

}

type Node4 struct {
	keys          [4]byte
	childpointers [4]*ARTNode
}

type Node16 struct {
	keys          [16]byte
	childpointers [16]*ARTNode
}

type Node48 struct {
	keysIndex     [256]uint8
	childpointers [48]*ARTNode
}

type Node256 struct {
	childpointers []*ARTNode
}

type Leaf struct {
	wholekey []byte
	value    *blockchain.Transaction
	leafhash [32]byte
}

func (n *Node4) makenode() *ARTNode {
	return &ARTNode{nodetype: 4,
		ChildCount: 0,
		isEnd:      false,
		//Edge: make([]byte, 0, 50),
		data: &Node4{keys: n.keys, childpointers: n.childpointers},
	}

}

func (n *Node16) makenode() *ARTNode {
	return &ARTNode{nodetype: 16,
		ChildCount: 0,
		isEnd:      false,
		data:       &Node16{keys: n.keys, childpointers: n.childpointers},
	}
}

func (n *Node48) makenode() *ARTNode {
	return &ARTNode{nodetype: 48,
		ChildCount: 0,
		isEnd:      false,
		data:       &Node48{keysIndex: n.keysIndex, childpointers: n.childpointers},
	}
}

func (n *Node256) makenode() *ARTNode {
	n.childpointers = make([]*ARTNode, 256)
	return &ARTNode{nodetype: 256,
		ChildCount: 0,
		isEnd:      false,
		data:       &Node256{childpointers: n.childpointers},
	}
}

func (l *Leaf) makenode() *ARTNode {
	//l.leafhash = ComputeLeafHash(l)
	return &ARTNode{nodetype: 1,
		ChildCount: 0,
		isEnd:      true,
		data:       &Leaf{wholekey: l.wholekey, value: l.value},
	}
}
