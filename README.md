# B-DNS Blockchain  

B-DNS is a blockchain-based decentralized DNS system using a **Proof-of-Stake (PoS)** consensus mechanism.  It stores domain records as transactions in an immutable ledger,  maintaining security and decentralization.
## Instructions

### Dependencies
Before running the project, ensure all dependencies are installed by executing:
   ```sh
   go mod tidy
   ```

### Run Program
   ```sh
   cd benchmark
   go test -bench=.
   ```

### Run Simulation
   ```sh
   make
   ```
## My contribution

I implemented the Adaptive Radix Tree (ART) index layer used for DNS record lookups (`art.go`, `artinsert.go`, `artsearch.go`, `artdelete.go`, `prefix.go`), and bug fixes to the blockchain and index manager components.

[See my full diff against the original here](https://github.com/bleasey/bdns/compare/master...prath2131:bdns_index:master)
