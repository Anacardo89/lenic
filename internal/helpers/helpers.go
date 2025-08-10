package helpers

import (
	"encoding/base64"
	"fmt"

	"github.com/Anacardo89/tpsi25_blog/internal/server"
	"github.com/google/uuid"
)

func OrderUUIDs(u1, u2 uuid.UUID) (uuid.UUID, uuid.UUID) {
	if u1.String() > u2.String() {
		return u2, u1
	}
	return u1, u2
}

func MakePasswordRecoverMail(user string, token string) string {
	encoded := base64.URLEncoding.EncodeToString([]byte(user))
	return fmt.Sprintf("https://%s:%s/recover-password/%s?token=%s", server.Server.Host, server.Server.HttpsPORT, encoded, token)
}

func MakeActivateUserLink(user string) string {
	encoded := base64.URLEncoding.EncodeToString([]byte(user))
	return fmt.Sprintf("https://%s:%s/action/activate/%s", server.Server.Host, server.Server.HttpsPORT, encoded)
}
