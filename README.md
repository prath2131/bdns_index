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
