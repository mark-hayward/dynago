package services

// TestService provides a testResults method
type TestService interface {
	TestResults() string
}

// NewTestService creates a new test service
func NewTestService() TestService {
	return &testService{}
}

type testService struct {
}

// testResults returns results of the testing for the Test service
func (h *testService) TestResults() string {
	return "Authentication Successful"
}
