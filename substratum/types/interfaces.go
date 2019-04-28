package types

// Assemblable objects can be turned into machine code
type Assemblable interface {
	Assemble() []byte
}

// Args represents actions which are common to instruction arguments
type Args interface {
	Verify() error
	Sanitize()
}
