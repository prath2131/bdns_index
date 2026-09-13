# B-DNS Blockchain  

BloQDNS is a blockchain-based decentralized DNS system using a **Proof-of-Stake (PoS)** consensus mechanism.  It stores domain records as transactions in an immutable ledger,  maintaining security and decentralization. 

Here an on-chain index structure based on adaptive radix tree is used to optimize query resolution.
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

### Run Linting
```sh
golangci-lint run  # to identify all the issues
golangci-lint run --fix # to automatically fix the fixable issues
```