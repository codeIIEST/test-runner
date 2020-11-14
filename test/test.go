package main

import (
	c "github.com/codeiiest/test-runner/test/languages/c"
	cpp "github.com/codeiiest/test-runner/test/languages/cpp"
)

func main() {
	// @TODO: add asserts to check with predefined test cases
	// Need to add large test cases and use concurrency as well
	c.Test()
	cpp.Test()
}
