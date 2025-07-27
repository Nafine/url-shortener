package random

import "math/rand"

func NewRandomPath(size int) string {
	alphabet := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-._~" // unreserved characters (rfc3986)

	path := make([]uint8, size)

	for i := 0; i < size; i++ {
		path[i] = alphabet[rand.Intn(len(alphabet))]
	}

	return string(path)
}
