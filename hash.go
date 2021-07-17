package gopar3

// func ChecksumValidate(b []byte) bool {
// 	// TODO: make sure this function does not cause race conditions due to using the table?
// 	length := len(b)
// 	if length < blockHashSize {
// 		return false
// 	}
// 	length -= blockHashSize
// 	return 0 == bytes.Compare(
// 		b[length:], Hash(b[:length]))
// }
