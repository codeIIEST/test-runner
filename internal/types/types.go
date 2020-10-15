package types

// TestData contains all test data related to the program
// to be evaluated
type TestData struct {
	ID         string
	Lang       string
	Filename   string
	Code       string
	TestCount  int
	InputData  []string
	OutputData []string
	TimeLimit  int64
	MemLimit   int64
}

// TestResult stores the evaluation data of the program
type TestResult struct {
	ID        string
	Status    []string
	Time      []float64
	Memory    []float64
	Error     string
	TestError []string
}
