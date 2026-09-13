package index

import (
	"bytes"

	"github.com/bleasey/bdns/internal/blockchain"
)

func BinarySearch(arr [16]byte, target byte, childCount int) int {
	left, right := 0, childCount-1

	for left <= right {
		mid := left + (right-left)/2

		if arr[mid] == target {
			return mid
		} else if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return -1 // not found
}

func (node *ARTNode) findChild(label byte) *ARTNode {
	// Implementation of findChild based on node type

	if node.nodetype == 4 {
		for i := 0; i < node.ChildCount; i = i + 1 {
			if node.data.(*Node4).keys[i] == label {
				return node.data.(*Node4).childpointers[i]
			}
		}
		return nil
	}
	if node.nodetype == 16 {
		keys := node.data.(*Node16).keys
		index := BinarySearch(keys, label, node.ChildCount)
		if index != -1 {
			return node.data.(*Node16).childpointers[index]
		}
		return nil
	}
	if node.nodetype == 48 {
		if node.data.(*Node48).keysIndex[label] != 0xFF {
			return node.data.(*Node48).childpointers[node.data.(*Node48).keysIndex[label]]
		} else {
			return nil
		}
	}
	if node.nodetype == 256 {
		return node.data.(*Node256).childpointers[label]
	}
	return nil
}

func (node *ARTNode) search(key []byte, depth int) (*blockchain.Transaction, bool) {
	if len(key) == 0 {
		return nil, false
	}
	if node == nil {
		return nil, false
	}
	if node.nodetype == 1 { // Leaf node
		if bytes.Equal(node.data.(*Leaf).wholekey, key) {
			return node.data.(*Leaf).value, true
		}
		return nil, false
	}
	// Traverse internal nodes
	p := len(LongestCommonPrefix((node.Edge), key[depth:]))
	if p != node.EdgeLen {
		return nil, false
	}
	depth = depth + int(node.EdgeLen)
	var next *ARTNode
	if depth >= len(key) {
		next = node.findChild(0x1A)
	} else {
		next = node.findChild(key[depth])
	}
	return next.search(key, depth+1)
}
