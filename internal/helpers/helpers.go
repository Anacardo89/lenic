package helpers

import (
	"regexp"

	"github.com/google/uuid"
)

func OrderUUIDs(u1, u2 uuid.UUID) (uuid.UUID, uuid.UUID) {
	if u1.String() > u2.String() {
		return u2, u1
	}
	return u1, u2
}

func ParseAtString(input string) []string {
	re := regexp.MustCompile(`@[\w]+`)
	return re.FindAllString(input, -1)
}
