package helpers

import (
	"errors"

	"grovehub/services/identity/internal/app/entity"
)

func NextStatus(status entity.AccountStatus, action string) (entity.AccountStatus, error) {
	switch action {
	case "suspend":
		if status == entity.Active {
			return entity.Suspended, nil
		}
	case "activate":
		if status == entity.Suspended {
			return entity.Active, nil
		}
	case "deactivate":
		if status != entity.Deactivated {
			return entity.Deactivated, nil
		}
	}
	return "", errors.New("invalid account status transition")
}
