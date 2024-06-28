//go:build darwin
// +build darwin

package teec

// Dummy UUID struct for macOS
type UUID struct {
	TimeLow          uint32
	TimeMid          uint16
	TimeHiAndVersion uint16
	ClockSeqAndNode  [8]byte
}

// Dummy TEEC struct for macOS
type TEEC struct{}

// Dummy New function
func New() *TEEC {
	return &TEEC{}
}

// Dummy Initialize function
func (t *TEEC) Initialize() error {
	return nil
}

// Dummy OpenSession function
func (t *TEEC) OpenSession(destination UUID) error {
	return nil
}

// Dummy InvokeCommand function
func (t *TEEC) InvokeCommand(commandID uint32, operation *Operation, origin *uint32) error {
	return nil
}

// Dummy CloseSession function
func (t *TEEC) CloseSession() error {
	return nil
}

// Dummy Finalize function
func (t *TEEC) Finalize() error {
	return nil
}

// Dummy Operation struct for macOS
type Operation struct {
	ParamTypes [4]ParameterTypes
	Params     [4]Parameter
}

// Dummy ParameterTypes and Parameter for macOS
type ParameterTypes int
type Parameter struct {
	Buffer []byte
	Size   uint32
}

const (
	MEMREF_TEMP_INPUT ParameterTypes = iota
	MEMREF_TEMP_INOUT
	NONE
)
