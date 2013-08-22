package timeout_controller

type MockTimeoutController struct {
	timeout chan bool
}

func MakeMockTimeoutController() *MockTimeoutController {
	return &MockTimeoutController{
		timeout: make(chan bool, 1),
	}
}

func (c *MockTimeoutController) Timeout() chan bool {
	return c.timeout
}

func (c *MockTimeoutController) Countdown() error {
	return nil
}

func (c *MockTimeoutController) Restart() {
}

func (c *MockTimeoutController) Stop() {
}
