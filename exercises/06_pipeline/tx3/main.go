package main

// EXERCISE: tx3 - Right Tool, Right Job
//
// PREDICT: Before writing any code, answer in your head:
//   When would you choose Pipeline over TxPipelined?
//   When would you choose TxPipelined over Pipeline?
//   What can a Lua script do that neither pipeline nor transaction can?
//
// TODO: Fix BestPrimitive -- all three cases return wrong values.

// BestPrimitive returns the right Redis primitive for each use case.
// Valid return values: "pipeline", "transaction", "lua"
//
// BUG: All three cases return the wrong primitive.
func BestPrimitive(scenario string) string {
	switch scenario {
	case "batch-100-independent-writes":
		// Need speed, no atomicity required -- each write is independent
		return "transaction" // BUG: should be "pipeline"

	case "atomic-debit-credit":
		// Must debit one account and credit another as a single unit
		return "pipeline" // BUG: should be "transaction"

	case "server-side-conditional-set":
		// Read a value, conditionally set another, atomically, server-side, no round trip
		return "pipeline" // BUG: should be "lua"
	}
	return "unknown"
}

func main() {}
