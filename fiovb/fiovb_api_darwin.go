//go:build darwin

package fiovb

// Define a dummy FIOVB struct and methods for macOS
type FIOVB struct{}

func New() *FIOVB {
	return &FIOVB{}
}

func (fiovb *FIOVB) Initialize() error {
	return nil
}

func (fiovb *FIOVB) Write(name, value string) error {
	return nil
}

func (fiovb *FIOVB) Read(name string) (string, error) {
	return "", nil
}

func (fiovb *FIOVB) Delete(name string) error {
	return nil
}

func (fiovb *FIOVB) Finalize() error {
	return nil
}
