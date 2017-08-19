/*
* CODE GENERATED AUTOMATICALLY WITH github.com/ernesto-jimenez/goautomock
* THIS FILE MUST NEVER BE EDITED MANUALLY
 */

package logger

// MarshallerMock mock
type MarshallerMock struct {
	calls map[string]int

	MarshalFunc func(*Entry) ([]byte, error)
}

func NewMarshallerMock() *MarshallerMock {
	return &MarshallerMock{
		calls: make(map[string]int),
	}
}

// Marshal mocked method
func (m *MarshallerMock) Marshal(p0 *Entry) ([]byte, error) {
	if m.MarshalFunc == nil {
		panic("unexpected call to mocked method Marshal")
	}
	m.calls["Marshal"]++
	return m.MarshalFunc(p0)
}

// MarshalCalls returns the amount of calls to the mocked method Marshal
func (m *MarshallerMock) MarshalTotalCalls() int {
	return m.calls["Marshal"]
}
