package req

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

func ParseUUIDOrGenerate(id string) (uuid.UUID, error) {
	if id == "" {
		return uuid.New(), nil
	}
	parsed, err := uuid.Parse(id)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("неверный формат UUID: %w", err)
	}
	return parsed, nil
}

func ParseUUIDPvz(id string) (uuid.UUID, error) {
	if id == "" {
		return uuid.New(), fmt.Errorf("неверный формат UUID")
	}
	parsed, err := uuid.Parse(id)
	if err != nil {
		return uuid.New(), fmt.Errorf("неверный формат UUID: %w", err)
	}
	return parsed, nil
}

func ParseTimeOrNow(value string) (time.Time, error) {
	if value == "" {
		return time.Now(), nil
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("неверный формат даты: %w", err)
	}
	return parsed, nil
}
