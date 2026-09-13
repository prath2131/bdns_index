package index

import "bytes"

//import ("fmt","sort")

func Delete(node *ARTNode, key []byte, depth int) (*ARTNode, bool) {

	if node == nil {
		return node, false
	}
	if len(key) == 0 {
		return node, false
	}

	if node.nodetype == 1 { //leafnode

		//key2 := []byte(node.data.(*Leaf).wholekey) //leafnode key

		// compute match length safely
		/*matchLen := 0
		for j := depth; j < len(key) && j < len(key2) && key[j] == key2[j]; j++ {
			matchLen++
		}
		depth = depth + matchLen
		//fmt.Println("inside delete function depth=",depth, "matchlen",len(node.data.(*Leaf).wholekey))
		if depth >= len(node.data.(*Leaf).wholekey) {
			node = nil
		}*/
		//fmt.Println("inside delete function node",node)
		if bytes.Equal(node.data.(*Leaf).wholekey, key) {
			return nil, true
		}
		return node, false
	}

	p := len(LongestCommonPrefix(node.Edge, key[depth:]))
	//p := checkPrefix(node, key, depth)
	//fmt.Println("prefixlen", p)
	if p != node.EdgeLen { //  prefix mismatch ie: if the key of the new leaf differs from a compressed path
		return node, false
	}

	depth = depth + node.EdgeLen
	var next *ARTNode
	if depth >= len(key) {
		next = node.findChild(0x1A)
	} else {
		next = node.findChild(key[depth])
	}
	// if next ==nil return false
	var result bool
	if next == nil {
		return node, false
	}

	next, result = Delete(next, key, depth+1)

	if !result {
		return node, false
	} else {
		if next == nil { // if child pointer points to nil, then remove this key and child
			replace(node, removeChild(node, key, depth))
		}
		// path compression case: if only one child and not root merge parent and child
		if node.ChildCount == 1 && node.EdgeLen != depth {
			node = node.MergeChild()
		} else if node.canShrink() { // not only child, shrink or adapt node
			node.shrink()
		}
		return node, true
	}
}

func (node *ARTNode) MergeChild() *ARTNode {
	onlyChild := node.getSingleChild()
	newEdge := node.Edge
	if onlyChild.nodetype == 1 {
		newEdge = append(newEdge, node.data.(*Node4).keys[0])
		newEdge = append(newEdge, onlyChild.Edge...)
		replace(node, onlyChild)
		node.Edge = newEdge
		node.EdgeLen = len(newEdge)
		//fmt.Println(string(node.Edge))
	}
	return node
}

func (node *ARTNode) getSingleChild() *ARTNode {
	if node.nodetype == 4 {
		return node.data.(*Node4).childpointers[0]
	}
	return nil
}

/*func isEmpty(child *Node) bool {
	if child == nil{
		return true
	}
	return false
}*/

func removeChild(node *ARTNode, key []byte, depth int) *ARTNode {

	var label byte
	if depth >= len(key) {
		label = 0x1A
	} else {
		label = key[depth]
	}

	if label == 0 {
		return nil
	}

	var index int
	switch node.nodetype {
	case 4:
		for i := 0; i < node.ChildCount+1; i = i + 1 {
			if node.data.(*Node4).keys[i] == label {
				node.data.(*Node4).keys[i] = 0
				node.data.(*Node4).childpointers[i] = nil
				index = i
				break
			}
		}
		for i := index; i < node.ChildCount-1; i = i + 1 { //shifting remaining left
			node.data.(*Node4).keys[i] = node.data.(*Node4).keys[i+1]
			node.data.(*Node4).childpointers[i] = node.data.(*Node4).childpointers[i+1]
		}
		node.data.(*Node4).keys[node.ChildCount-1] = 0
		node.data.(*Node4).childpointers[node.ChildCount-1] = nil
		node.ChildCount = node.ChildCount - 1

	case 16:
		for i := 0; i < node.ChildCount+1; i = i + 1 {
			if node.data.(*Node16).keys[i] == label {
				node.data.(*Node16).keys[i] = 0
				node.data.(*Node16).childpointers[i] = nil
				index = i
				break
			}
		}
		for i := index; i < node.ChildCount-1; i = i + 1 { //shifting remaining left
			node.data.(*Node16).keys[i] = node.data.(*Node16).keys[i+1]
			node.data.(*Node16).childpointers[i] = node.data.(*Node16).childpointers[i+1]

		}
		node.data.(*Node16).keys[node.ChildCount-1] = 0
		node.data.(*Node16).childpointers[node.ChildCount-1] = nil
		node.ChildCount = node.ChildCount - 1

	case 48:
		childindex := node.data.(*Node48).keysIndex[label]
		node.data.(*Node48).childpointers[childindex] = nil
		node.data.(*Node48).keysIndex[label] = 0xFF
		node.ChildCount = node.ChildCount - 1
	case 256:
		node.data.(*Node256).childpointers[label] = nil
		node.ChildCount = node.ChildCount - 1
	}

	return node
}

func (node *ARTNode) canShrink() bool {
	switch node.nodetype {
	case 16:
		if node.ChildCount <= 4 {
			return true
		}
	case 48:
		if node.ChildCount <= 16 {
			//fmt.Println("child",node.ChildCount)
			return true
		}
	case 256:
		if node.ChildCount <= 48 {
			return true
		}
	default:
		return false
	}
	return false
}

func (node *ARTNode) shrink() {
	var newNode *ARTNode
	if node.nodetype == 16 { // downgrade to node 4
		n := &Node4{}
		newNode = n.makenode()
		for i := 0; i < node.ChildCount; i++ {
			newNode.data.(*Node4).keys[i] = node.data.(*Node16).keys[i]
			newNode.data.(*Node4).childpointers[i] = node.data.(*Node16).childpointers[i]
			newNode.ChildCount += 1
		}
	}

	if node.nodetype == 48 {
		n := &Node16{}
		newNode = n.makenode()
		//var keys []byte
		var n16index int
		for index, value := range node.data.(*Node48).keysIndex {
			if value != 0xFF {
				//keys = append(keys, byte(index))
				newNode.data.(*Node16).keys[n16index] = byte(index)
				newNode.data.(*Node16).childpointers[n16index] = node.data.(*Node48).childpointers[value]
				n16index = n16index + 1
			}
		}
	}

	if node.nodetype == 256 {
		n := &Node48{}
		newNode = n.makenode()

		newNode.data.(*Node48).keysIndex[0] = 0xFF
		for i := 1; i < 256; i *= 2 { //initialize all keysIndex = 0xFF logarithmically
			copy(newNode.data.(*Node48).keysIndex[i:], newNode.data.(*Node48).keysIndex[:i])
		}
		var ChildCount uint8
		for idx, ptrvalue := range node.data.(*Node256).childpointers {
			if ptrvalue != nil {
				newNode.data.(*Node48).keysIndex[idx] = ChildCount
				newNode.data.(*Node48).childpointers[ChildCount] = ptrvalue
				ChildCount++
			}

		}

	}

	newNode.ChildCount = node.ChildCount
	newNode.Edge = node.Edge
	newNode.EdgeLen = node.EdgeLen
	newNode.isEnd = node.isEnd
	replace(node, newNode)

	//fmt.Println(node)
}
