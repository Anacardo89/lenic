package helpers

import (
	"encoding/base64"
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func OrderUUIDs(u1, u2 uuid.UUID) (uuid.UUID, uuid.UUID) {
	if u1.String() > u2.String() {
		return u2, u1
	}
	return u1, u2
}

func MakePasswordRecoverMail(host, port, user, token string) string {
	encoded := base64.URLEncoding.EncodeToString([]byte(user))
	return fmt.Sprintf("https://%s:%s/recover-password/%s?token=%s", host, port, encoded, token)
}

func MakeActivateUserLink(host, port, user string) string {
	encoded := base64.URLEncoding.EncodeToString([]byte(user))
	return fmt.Sprintf("https://%s:%s/action/activate/%s", host, port, encoded)
}

func ParseAtString(input string) []string {
	re := regexp.MustCompile(`@[\w]+`)
	return re.FindAllString(input, -1)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
