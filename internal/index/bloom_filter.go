package index

import (
	"github.com/bits-and-blooms/bloom/v3"
)

type BloomFilterManager struct {
	revocationFilter *bloom.BloomFilter
	validFilter      *bloom.BloomFilter
}

func InitFilter() *BloomFilterManager {
	return &BloomFilterManager{
		revocationFilter: bloom.NewWithEstimates(1000, 0.01),
		validFilter:      bloom.NewWithEstimates(1000, 0.01),
	}
}

func (bfm *BloomFilterManager) AddToValidList(validDomain string) {
	bfm.validFilter.AddString(validDomain)
}

func (bfm *BloomFilterManager) AddToRevocationList(revokedDomain string) {
	bfm.revocationFilter.AddString(revokedDomain)
}

func (bfm *BloomFilterManager) IsValid(domain string) bool {
	return bfm.validFilter.TestString(domain)
}

func (bfm *BloomFilterManager) IsRevoked(domain string) bool {
	return bfm.revocationFilter.TestString(domain)
}
