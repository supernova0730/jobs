package uuid

import "github.com/gofrs/uuid"

func Generate() string {
	return uuid.Must(uuid.NewV4()).String()
}
