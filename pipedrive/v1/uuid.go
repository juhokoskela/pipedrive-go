package v1

import (
	"fmt"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

func parseUUID(value string, label string) (openapi_types.UUID, error) {
	parsed, err := uuid.Parse(value)
	if err != nil {
		return openapi_types.UUID{}, fmt.Errorf("parse %s: %w", label, err)
	}
	return openapi_types.UUID(parsed), nil
}
