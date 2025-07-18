# AGENTS Instructions

These guidelines apply to the entire repository. Follow them when adding or modifying code.

## Implementation Plan

USDWZ will be built as a dedicated Cosmos SDK blockchain with modules for stablecoin issuance and validator arbitration. Key modules:

- **stablecoin**: manage minting/burning of USDWZ backed by M^0 deposits.
- **escrow**: handle milestone-based escrows and validator approvals.
- **validator**: track votes and apply slashing for malicious behavior.
- **yield**: distribute yield from M^0 reserves.

M^0 acts as the underlying collateral. Every USDWZ minted requires an equivalent M^0 deposit. Burning USDWZ releases M^0 back to the user. Governance can pause minting or blacklist addresses for compliance.

### Technical Implementation Plan

1. **Environment & Repository Setup**
   - Scaffold a new Cosmos SDK chain named `usdWZ`.
   - Configure GitHub Actions with Go 1.21+ to run `go vet`, `gofmt`, and `go test ./...`.

2. **Core Modules**

   | Module | Responsibilities | Key Functions |
   |--------|-----------------|---------------|
   |`stablecoin`|Handles USDWZ minting, burning, and M^0-backed collateral management.|`MsgMint`, `MsgBurn`, `DepositCollateral`, `RedeemCollateral`|
   |`escrow`|Manages milestone-based escrows between requesters and labelers.|`CreateEscrow`, `SubmitMilestone`, `VoteMilestone`, `FinalizeEscrow`|
|`validator`|Tracks validator voting for milestone approval and exposes slashing logic.|`SubmitVote`, `TallyVotes`, `Slash`|
|`yield`|Distributes M^0 yield to USDWZ holders.|`AccrueYield`, `DistributeYield`|
|`audit`|Publishes on-chain attestations for reserve verification.|`PublishAttestation`|
|`monitoring`|Exposes Prometheus metrics for collateral and validator performance.|`StartServer`, `SetCollateral`|

Label validators participate directly in consensus. Votes on each milestone are
recorded on-chain and determine whether escrowed funds are released or
refunded. Disputes are resolved automatically through block finality.

### Validator Arbitration Architecture

Milestone approval rules are encoded as simple scripts interpreted by the
`escrow` module. Each script references validator sets and uses opcodes to
define how many approvals are needed. The goal is to avoid complex
smart-contract logic while enabling deterministic execution within consensus.

- **OP_SET `<id>`** – load a predefined validator set (`A`, `B`, etc.).
- **OP_ANY** – approval succeeds if any validator in the current set votes `yes`.
- **OP_QUORUM `<n>`** – succeed once `n` validators in the set vote `yes`.
- **OP_ALL** – require unanimous approval from the current set.
- **OP_THEN** – chain multiple conditions sequentially.

Example scripts:

1. `OP_SET A OP_ANY` – any validator in set A can approve.
2. `OP_SET A OP_QUORUM 2` – at least two of three validators in set A approve.
3. `OP_SET A OP_ALL` – unanimous approval from set A.
4. `OP_SET A OP_QUORUM 2 OP_THEN OP_SET B OP_ALL` – quorum from set A followed by unanimous approval from set B.

Validators stake tokens and earn a portion of escrow fees for timely votes. The
`validator` module slashes stake for failing to vote or for submitting invalid
votes. This incentive model aligns economic rewards with accurate milestone
validation.

3. **M^0 Integration**
   - Maintain a module account mirroring M^0 deposits.
   - Burn USDWZ to release M^0 via IBC transfer.
   - Mint additional USDWZ to distribute yield.
   - Verify each deposit transaction before minting to maintain 1:1 collateralization.
   - Track yield accrual from M^0 and automatically mint distribution via the yield module.

4. **Governance & Security**
   - Parameter changes and emergency controls via on-chain proposals.
   - Implement slashing and blacklisting for malicious actors.
   - Manage chain treasury via multisig custody to secure collateral deposits.
5. **APIs & Wallet Support**
   - Expose REST and gRPC endpoints.
   - Keep ICS20 transfers enabled for interoperability.

6. **Test Plan**
   - Unit tests for each module.
   - Integration tests launching a local chain and mocking M^0 deposits.
   - Tests ensuring the M^0 collateral account equals or exceeds supply.
   - GitHub Actions jobs for formatting and tests.
   - Include testcases for the `monitoring` and `audit` packages to validate metrics output and attestation storage.

7. **Deployment, Monitoring & Operational Housekeeping**
   - Provide validator setup instructions and monitor chain metrics.
   - Use Prometheus and Grafana for node and module monitoring with alerts for collateral discrepancies. The `monitoring` package must expose a metrics server with tests covering its gauges.
   - Implement secure custody and key-management procedures for mint and burn operations.
   - Track treasury balances and validator behavior in real time with automated alerts.
   - Schedule routine external audits of reserves and code; publish attestation reports using the `audit` module and its test suite.
   - Maintain relationships with banking partners to streamline collateral deposits and redemptions.
   - Enforce KYC/AML checks for all mint and redemption operations. Keep guidance for institutions in `docs/kyc_aml.md`.

8. **Launch Sequence**
   - Complete module implementations and unit/integration tests.
   - Deploy a public testnet for community testing and validator arbitration flows.
   - Conduct a thorough security audit.
   - Create the genesis block, allocate initial tokens, and launch the mainnet.
   - Publish documentation for institutions on minting and redemption with KYC/AML guidelines.

## Libraries and Tools

- Go 1.21+
- cosmos-sdk (latest stable release)
- cometbft for consensus
- IBC modules for token transfers
- GitHub Actions for CI

## Code Structure

```
/chain         Cosmos SDK app source
  /cmd         Main entrypoint
  /x           Custom modules
    /stablecoin
    /escrow
    /validator
    /yield
/docs          Documentation
/scripts       Helper scripts (start chain, integration tests)
/tests         Go test suites
```

Each module defines `keeper`, `types`, `msgs`, and `client` packages according to Cosmos SDK conventions. Use clear interfaces and keep logic in keepers.

## Testcase Plan

1. **Unit Tests**: `go test ./...` covers all modules.
2. **Integration Tests**: launch a local chain inside tests (`go test -tags=integration ./tests/integration`). Mock M^0 deposits, validator votes, and the four approval scripts.
3. **M^0 Tests**: verify collateral ratio and yield distribution.
4. Continuous integration must run formatting and tests on every commit.

## Coding Standards

- Use `gofmt` and `go vet` before committing.
- Add Go doc comments for all exported identifiers.
- Keep functions under 50 lines where possible and favor readability over micro-optimizations.
- Follow Cosmos SDK best practices for module organization and error handling.
- Commit messages should be concise and in present tense, e.g., "Add stablecoin keeper".

## Continuous Integration

GitHub Actions workflow should:

1. Set up Go 1.21.
2. Run `go mod download`.
3. Run `gofmt -w` on all Go files and fail if changes are needed.
4. Run `go vet ./...`.
5. Run `go test ./...` (and integration tests as a separate job).

