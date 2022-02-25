package util

type HttpResp struct {
	Data      []byte
	Cookie    *map[string]string
	Localtion string
	Error     error
}
