package types

const (
	// ModuleName defines the audit module name.
	ModuleName = "audit"
	// StoreKey defines the primary module store key.
	StoreKey = ModuleName
)

// Attestation holds a simple audit result.
type Attestation struct {
	Timestamp int64  `json:"timestamp"`
	Hash      string `json:"hash"`
}
