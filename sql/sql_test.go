package sql

var (
	// Make sure DB adheres to Interface
	_ Interface = new(DB)

	// Make sure Tx adheres to Interface
	_ Interface = new(Tx)
)
