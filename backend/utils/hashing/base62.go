package hashing

const base62Alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func EncodeBase62(num uint64) string {
	if num == 0 {
		return "0"
	}

	buf := make([]byte, 0, 11) // max length for uint64

	for num > 0 {
		rem := num % 62
		buf = append(buf, base62Alphabet[rem])
		num /= 62
	}

	// reverse
	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		buf[i], buf[j] = buf[j], buf[i]
	}

	return string(buf)
}

func DecodeBase62(s string) uint64 {
	var num uint64
	for _, c := range s {
		num = num*62 + uint64(indexInBase62Alphabet(byte(c)))
	}
	return num
}

func indexInBase62Alphabet(character byte) int {
	for i := range base62Alphabet {
		if base62Alphabet[i] == character {
			return i
		}
	}
	return -1
}
