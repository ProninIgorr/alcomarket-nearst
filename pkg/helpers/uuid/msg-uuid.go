package uuid

import (
	"strconv"

	guid "github.com/satori/go.uuid"
)

// GenerateMsgUUID Генерирует uuid
func GenerateUUID() string {
	uuid := guid.Must(guid.NewV4(), nil)

	return uuid.String()
}

// GenerateUUIDV3 Генерирует uuid по входящей строке
func GenerateUUIDV3(base string) string {
	uuid := guid.Must(guid.NewV3(guid.NamespaceDNS, base), nil)

	return uuid.String()
}

// GenerateUUIDV3FromInt Генерирует uuid по входящему числу
func GenerateUUIDV3FromInt(base int) string {
	return GenerateUUIDV3(strconv.Itoa(base))
}

// GenerateUUIDV3FromIntWithPrefix Генерирует uuid по входящему числу и префиксу
func GenerateUUIDV3FromIntWithPrefix(prefix string, base int) string {
	return GenerateUUIDV3(prefix + strconv.Itoa(base))
}
