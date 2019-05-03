package window

import "fmt"

// Various command bits we need
const openCmd = 136 // 10001000
const stopCmd = 72  // 01001000
const closeCmd = 40 // 00101000

// Window represents a window and allows for opening, stopping and closing
type Window struct {
	// ID identifies the window in question
	// and should only be 5 bytes long
	ID []byte

	Radio func([]byte) error
}

// VerifyID checks for a correct ID
func (w *Window) VerifyID() error {
	if len(w.ID) != 5 {
		return fmt.Errorf("wrong ID length: %d, should be exactly 5 bytes", len(w.ID))
	}

	return nil
}

// Pair is just a convenience method of calling "Open"
func (w *Window) Pair() error {
	return w.Open()
}

// Open opens a window
// Example protocol:
//   <----------- device id --------------------> <-cmd -> <- XOR >
//   177      167      68       228      115      136      77
//   10110001 10100111 01000100 11100100 01110011 10001000 01001101
func (w *Window) Open() error {
	payload := w.ID
	payload = append(payload, openCmd)
	payload = append(payload, checksum(payload))

	return w.Radio(payload)
}

// Stop stops a window from whatever its doing (closing or opening)
// Example protocol:
//   <----------- device id --------------------> <-cmd -> <- XOR >
//   177      167      68       228      115      72       141
//   10110001 10100111 01000100 11100100 01110011 01001000 10001101
func (w *Window) Stop() error {
	payload := w.ID
	payload = append(payload, stopCmd)
	payload = append(payload, checksum(payload))

	return w.Radio(payload)
}

// Close closes a window
// Example protocol:
//   <----------- device id --------------------> <-cmd -> <- XOR >
//   177      167      68       228      115      40       237
//   10110001 10100111 01000100 11100100 01110011 00101000 11101101
func (w *Window) Close() error {
	payload := w.ID
	payload = append(payload, closeCmd)
	payload = append(payload, checksum(payload))

	return w.Radio(payload)
}
