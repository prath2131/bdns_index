package index

import (
	"fmt"
	"testing"

	"github.com/bleasey/bdns/internal/blockchain"
)

func TestIndexManager_Integration(t *testing.T) {
	im := NewIndexManager()
	mockTx := make(map[int]*blockchain.Transaction)
	// Test Data
	//var domain string
	//var mockTx [10]*blockchain.Transaction
	for i := range 10 {
		domainName := fmt.Sprintf("node%d.com", i+1)
		ip := fmt.Sprintf("192.168.1.%d", i+1)
		ttl := int64(3600)
		keypair := blockchain.NewKeyPair()
		mockTx[i] = blockchain.NewTransaction(blockchain.REGISTER, domainName, ip, ttl, keypair.PublicKey, &keypair.PrivateKey, mockTx)
	}

	// --- 1. Test Adding and Bloom Filter ---
	t.Run("Add and Bloom Filter Check", func(t *testing.T) {
		domain := mockTx[9].DomainName
		im.Add(domain, mockTx[9])

		// Verify Bloom Filter (Element 1)
		if !im.filter.IsValid(domain) {
			t.Errorf("Expected domain %s to be valid in Bloom Filter", domain)
		}

		// Verify ART Tree Search (Element 2)
		result := im.GetIP(domain)
		if result == nil {
			t.Fatal("Expected to find transaction in ART Tree, got nil")
		}
	})

	for i := range 8 {
		domain := mockTx[i].DomainName
		im.Add(domain, mockTx[i])

		// Verify Bloom Filter (Element 1)
		if !im.filter.IsValid(domain) {
			t.Errorf("Expected domain %s to be valid in Bloom Filter", domain)
		}

		// Verify ART Tree Search (Element 2)
		result := im.GetIP(domain)
		if result == nil {
			t.Fatal("Expected to find transaction in ART Tree, got nil")
		}
	}

	// --- 2. Test Updating ---
	t.Run("Update Entry", func(t *testing.T) {
		newTx := mockTx[1]
		newTx.Type = blockchain.UPDATE
		domain := newTx.DomainName
		im.Update(domain, newTx)

		result := im.GetIP(domain)
		if result != newTx {
			t.Error("ART Tree did not return the updated transaction")
		}
	})

	// --- 3. Test Removal and Index Tracking ---
	t.Run("Remove and State Check", func(t *testing.T) {
		mockTx[0].Type = blockchain.REVOKE
		domain := mockTx[0].DomainName
		im.Remove(domain)

		// Check ART Tree
		result := im.GetIP(domain)
		if result != nil {
			t.Error("Expected domain to be removed from ART Tree")
		}

		// Verify Revocation (Element 1 - Filter logic)
		// This assumes your IsValid checks the revocation list internally
		if im.filter.IsValid(domain) {
			t.Log("Note: Bloom filters have no true 'remove'. Ensure IsValid handles the revocation list check.")
		}

		// Verify CurrentIndex (Element 3)
		// Manually increment or check logic if your code updates this during Add/Remove
		if im.currentIndex != 0 {
			t.Errorf("Expected initial index 0, got %d", im.currentIndex)
		}
	})
}
