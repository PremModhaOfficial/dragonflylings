package main

func BestPrimitive(scenario string) string {
	switch scenario {
	case "batch-100-independent-writes":
		return "pipeline"
	case "atomic-debit-credit":
		return "transaction"
	case "server-side-conditional-set":
		return "lua"
	}
	return "unknown"
}

func main() {}
