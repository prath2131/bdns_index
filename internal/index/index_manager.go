package index

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/bleasey/bdns/internal/blockchain"
)

// revive:disable-next-line
type IndexManager struct {
	tree *ARTTree
	//tree         *AVLTree
	filter       *BloomFilterManager
	currentIndex int64
}

func NewIndexManager() *IndexManager {
	return &IndexManager{
		tree: &ARTTree{&ARTNode{}},
		//tree:         &AVLTree{nil},
		filter:       InitFilter(),
		currentIndex: 0, // Corresponds to genesis block
	}
}

func (im *IndexManager) GetIP(domain string) *blockchain.Transaction {
	// Check if the domain is valid
	if !im.filter.IsValid(domain) {
		//fmt.Printf("Domain %s is not valid.\n", domain)
		return nil
	}

	targetNode, found := im.tree.Search(reverseDomain([]byte(domain)))
	if !found {
		//targetNode := im.tree.Search(HashDomain(domain))
		//if targetNode == nil {
		//fmt.Printf("Domain %s not found in ART tree.\n", domain)
		return nil
	}

	return targetNode
}

//func (im *IndexManager) GetIndexHash() []byte {
//return ComputeIndexNodeHash(im.tree.root)
//}

func (im *IndexManager) Add(domain string, tx *blockchain.Transaction) {
	im.tree.Add(reverseDomain([]byte(domain)), tx)
	//im.tree.Add(HashDomain(domain), tx)
	im.filter.AddToValidList(domain)
}

func (im *IndexManager) Update(domain string, tx *blockchain.Transaction) {
	im.tree.Update(reverseDomain([]byte(domain)), tx)
	//im.tree.Update(HashDomain(domain), tx)
}

func (im *IndexManager) Remove(domain string) {
	im.tree.Remove(reverseDomain([]byte(domain)))
	//im.tree.Remove(HashDomain(domain))
	im.filter.AddToRevocationList(domain)
}

func HashDomain(domain string) string {
	hash := sha256.Sum256([]byte(domain))
	return hex.EncodeToString(hash[:])
}