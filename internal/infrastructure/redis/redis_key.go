package redis

import "fmt"

var (
	ChatPrivateKey = func(label string) string {
		return fmt.Sprintf("chat:private:%v", label)
	}
)
