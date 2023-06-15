package hypvsr

type MockTemplate struct {
	name string
}

func (m *MockTemplate) Load() (KindManager, error) {
	// Implement the Load method
	return &MockKindManager{}, nil
}
