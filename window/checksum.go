package window

// checksum just XOR's all bytes togeather
func checksum(p []byte) byte {
	var c byte
	for _, b := range p {
		c = c ^ b
	}
	return c
}
