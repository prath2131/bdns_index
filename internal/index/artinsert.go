package index

import "bytes"

//import "fmt"

func replace(old *ARTNode, new *ARTNode) {
	if new == nil {
		return
	}

	*old = *new
}

func (node *ARTNode) isFull() bool {
	// Check if the node is full based on its type
	switch node.nodetype {
	case 4:
		if node.ChildCount >= 4 {
			return true
		}
	case 16:
		if node.ChildCount >= 16 {
			return true
		}
	case 48:
		if node.ChildCount >= 48 {
			return true
		}
	default:
		return false
	}
	return false
}

func (node *ARTNode) grow() {
	// Grow the node to the next larger type
	var newNode *ARTNode
	if node.nodetype == 4 { //upgrade to node 16
		n := &Node16{}
		newNode = n.makenode()

		for i := 0; i < node.ChildCount; i++ {
			label := node.data.(*Node4).keys[i]
			j := newNode.ChildCount - 1
			for j >= 0 && node.data.(*Node4).keys[j] > label {
				newNode.data.(*Node16).keys[j+1] = node.data.(*Node4).keys[j]
				newNode.data.(*Node16).childpointers[j+1] = node.data.(*Node4).childpointers[j]
				j = j - 1
			}
			newNode.data.(*Node16).keys[j+1] = label
			newNode.data.(*Node16).childpointers[j+1] = node.data.(*Node4).childpointers[i]
			newNode.ChildCount += 1
		}

	}

	if node.nodetype == 16 {
		n := &Node48{}
		newNode = n.makenode()

		newNode.data.(*Node48).keysIndex[0] = 0xFF
		//newNode.data.(*Node48).keysIndex[1] = 0xFF
		for i := 1; i < 256; i *= 2 { //initialize all keysIndex = 0xFF logarithmically
			copy(newNode.data.(*Node48).keysIndex[i:], newNode.data.(*Node48).keysIndex[:i])
		}

		for i := 0; i < 16; i++ {
			label := node.data.(*Node16).keys[i]
			newNode.data.(*Node48).keysIndex[label] = uint8(i)
			newNode.data.(*Node48).childpointers[uint8(i)] = node.data.(*Node16).childpointers[i]
		}
	}

	if node.nodetype == 48 {
		n := &Node256{}
		newNode = n.makenode()
		//var label byte
		for index, value := range node.data.(*Node48).keysIndex {
			//label := byte(node.data.(*Node48).keysIndex[i])
			if value != 0xFF {
				//fmt.Println(index,value)
				newNode.data.(*Node256).childpointers[index] = node.data.(*Node48).childpointers[value]
			}
		}

	}
	// copy metadata
	newNode.ChildCount = node.ChildCount
	newNode.Edge = node.Edge
	newNode.EdgeLen = node.EdgeLen
	newNode.isEnd = node.isEnd

	replace(node, newNode)
}

func (node *ARTNode) addChild(childnode *ARTNode, key []byte, depth int) {
	// Add child to the node based on its type

	// Guard against out of bounds
	if depth < 0 {
		return
	}

	var label byte

	if depth >= len(key) {
		label = 0x1A
	} else {
		label = key[depth]
	}

	next := node.findChild(label)
	if next != nil { //no duplicate keys
		return
	}

	if depth == len(key) {
		childnode.Edge = nil
		childnode.EdgeLen = 0
	} else {
		childnode.Edge = key[depth+1:]
		childnode.EdgeLen = len(childnode.Edge)
	}

	node.ChildCount += 1
	switch node.nodetype {
	case 4:
		if node.ChildCount <= 4 {
			i := node.ChildCount - 2
			for i >= 0 && node.data.(*Node4).keys[i] > label {
				node.data.(*Node4).keys[i+1] = node.data.(*Node4).keys[i]
				node.data.(*Node4).childpointers[i+1] = node.data.(*Node4).childpointers[i]
				i = i - 1
			}
			node.data.(*Node4).keys[i+1] = label
			node.data.(*Node4).childpointers[i+1] = childnode
		} else {
			node.grow()
			node.addChild(childnode, key, depth)
		}
	case 16:
		if node.ChildCount <= 16 {
			// Insert in sorted order
			i := node.ChildCount - 2
			for i >= 0 && node.data.(*Node16).keys[i] > label {
				node.data.(*Node16).keys[i+1] = node.data.(*Node16).keys[i]
				node.data.(*Node16).childpointers[i+1] = node.data.(*Node16).childpointers[i]
				i = i - 1
			}
			node.data.(*Node16).keys[i+1] = label
			node.data.(*Node16).childpointers[i+1] = childnode
		} else {
			node.grow()
			node.addChild(childnode, key, depth)
		}
	case 48:
		if node.ChildCount <= 48 {
			node.data.(*Node48).keysIndex[label] = uint8(node.ChildCount - 1)
			node.data.(*Node48).childpointers[uint8(node.ChildCount-1)] = childnode
		} else {
			node.grow()
			node.addChild(childnode, key, depth)
		}
	case 256:
		node.data.(*Node256).childpointers[label] = childnode
	}
}

/*
func loadkey(node *Node) []byte {
	// Load the key from the leaf node
	return []byte{}
}
func makeNode4(keyByte byte, child *Node ) *Node {
	n:= &Node4{}
		n.keys[0]=keyByte
		n.childpointers[0]= child
	return n.makenode()
}*/

func preserveNode(node *ARTNode) *ARTNode {
	oldNode := &ARTNode{
		nodetype:   node.nodetype,
		ChildCount: node.ChildCount,
		isEnd:      node.isEnd,
		Edge:       make([]byte, len(node.Edge)),
		EdgeLen:    node.EdgeLen,
	}
	copy(oldNode.Edge, node.Edge)

	// Copy the node's data based on its type
	switch node.nodetype {
	case 1:
		oldData := &Leaf{wholekey: node.data.(*Leaf).wholekey, value: node.data.(*Leaf).value} //, leafhash: node.data.(*Leaf).leafhash}
		oldNode.data = oldData
	case 4:
		oldData := &Node4{}
		*oldData = *node.data.(*Node4) // Copy all fields
		oldNode.data = oldData
	case 16:
		oldData := &Node16{}
		*oldData = *node.data.(*Node16)
		oldNode.data = oldData
	case 48:
		oldData := &Node48{}
		*oldData = *node.data.(*Node48)
		oldNode.data = oldData
	case 256:
		oldData := &Node256{}
		*oldData = *node.data.(*Node256)
		oldNode.data = oldData
	}
	return oldNode
}

/*func preserveNodeShallow(n *Node) *Node {
    oldNode := *n  // struct copy (NO allocations yet)
    return &oldNode
}*/

func loadleafkey(node *ARTNode) []byte {
	return node.data.(*Leaf).wholekey
}

func (node *ARTNode) lazyexpand(matchLen int, key []byte, depth int, newleaf *ARTNode) {
	// Expand a leaf node lazily into a Node4
	n := &Node4{}
	newNode := n.makenode()
	oldNode := preserveNode(node)
	newNode.EdgeLen = matchLen
	newNode.Edge = make([]byte, matchLen)
	if node.nodetype == 1 {
		copy(newNode.Edge[:matchLen], key[depth:depth+matchLen])
	} else {
		copy(newNode.Edge[:matchLen], node.Edge[:matchLen])
		//oldNode.Edge = oldNode.Edge[matchLen:]
		//oldNode.EdgeLen = oldNode.EdgeLen - (matchLen)

	}

	if node.nodetype == 1 {
		newNode.addChild(oldNode, loadleafkey(oldNode), depth+matchLen)
	} else {
		newNode.addChild(oldNode, oldNode.Edge, matchLen)
	}
	newNode.addChild(newleaf, key, depth+matchLen)
	replace(node, newNode)

}

func (node *ARTNode) Insert(key []byte, depth int, newleaf *ARTNode) {
	if node == nil {
		return
	}
	if len(key) == 0 {
		return
	}
	if node.nodetype == 0 { // empty node
		newleaf.Edge = key
		newleaf.EdgeLen = len(key)
		replace(node, newleaf)
		//fmt.Println(node)
		return
	}

	if node.nodetype == 1 { //case: lazy expansion= if encountered existing leaf, replace with new node and existing leaf
		key2 := loadleafkey(node) //leafnode key

		if bytes.Equal(key2, key) { //no duplicate key
			return
		} else {
			matchLen := len(LongestCommonPrefix(key[depth:], key2[depth:]))
			node.lazyexpand(matchLen, key, depth, newleaf)
		}
		//fmt.Println(node)
		return
	}

	p := len(LongestCommonPrefix(node.Edge, key[depth:]))
	if p < node.EdgeLen { //  prefix mismatch ie: if the key of the new leaf differs from a compressed path
		node.lazyexpand(p, key, depth, newleaf)
		//fmt.Println(node)
		return
	}

	depth = depth + node.EdgeLen
	if depth >= len(key) {
		node.addChild(newleaf, key, depth)
		return
	}
	next := node.findChild(key[depth])
	if next != nil {
		next.Insert(key, depth+1, newleaf)
	} else {
		// Add to inner node if no matching child is found
		if node.isFull() {
			node.grow()
		}
		node.addChild(newleaf, key, depth)
		//fmt.Println(node)
		return
	}
}
