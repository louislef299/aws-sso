package os

import "encoding/base64"

func Decode(s string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func Encode(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}
