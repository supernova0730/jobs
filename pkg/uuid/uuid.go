package uuid

import "github.com/gofrs/uuid"

func Generate() string {
	result, _ := uuid.NewV4()
	return result.String()
}
