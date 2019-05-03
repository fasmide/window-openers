package window

func checksum(p []byte) byte {
	var c byte
	for _, b := range p {
		c = c ^ b
	}
	return c
}
