package gosss

import "math/big"

// encodeMessage function splits a message into parts of a given size and
// converts them to big.Int. It returns an error if the message cannot be
// encoded into a big.Int. If the given message is smaller than the part size,
// it returns a single part.
func encodeMessage(message []byte, partSize int) []*big.Int {
	if len(message) <= partSize {
		return []*big.Int{new(big.Int).SetBytes(message)}
	}
	var parts []*big.Int
	for i := 0; i < len(message); i += partSize {
		end := i + partSize
		if end > len(message) {
			end = len(message)
		}
		parts = append(parts, new(big.Int).SetBytes(message[i:end]))
	}
	return parts
}

// decodeMessage function converts the parts of a message to a single string.
// It returns the decoded message. It uses the bytes of the big.Int to decode
// the message, appending them to a single byte slice and converting it to a
// string.
func decodeMessage(parts []*big.Int) []byte {
	var bMessage []byte
	for _, part := range parts {
		bMessage = append(bMessage, part.Bytes()...)
	}
	return bMessage
}
