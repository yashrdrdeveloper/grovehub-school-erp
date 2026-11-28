package helpers

import (
	"gorm.io/gorm"

	coreid "grovehub/libs/core/id"
	"grovehub/services/identity/internal/app/constants/feature"
	"grovehub/services/identity/internal/app/constants/permission"
	"grovehub/services/identity/internal/app/repository/registry"
	"grovehub/services/identity/internal/app/repository/registry/constants"
)

type registeredDefinition struct {
	ID        string
	Module    string
	FeatureID *string
}

func RegisterFeature(db *gorm.DB, definition feature.Definition) (int, error) {
	identifier, err := coreid.NewUUIDv7()
	if err != nil {
		return 0, err
	}
	result := db.Exec(constants.InsertFeature, identifier, definition.Code, constants.Module, definition.Name)
	if result.Error != nil {
		return 0, result.Error
	}
	if _, err := findFeature(db, definition.Code); err != nil {
		return 0, err
	}
	return int(result.RowsAffected), nil
}

func findFeature(db *gorm.DB, code feature.Code) (string, error) {
	var stored registeredDefinition
	result := db.Raw(constants.FindFeature, code).Scan(&stored)
	if result.Error != nil {
		return "", result.Error
	}
	if result.RowsAffected == 0 {
		return "", registry.ErrFeatureMissing
	}
	if result.RowsAffected != 1 || stored.ID == "" || stored.Module != constants.Module {
		return "", registry.ErrConflict
	}
	return stored.ID, nil
}

func RegisterPermission(db *gorm.DB, definition permission.Definition) (int, error) {
	var featureID *string
	if definition.Feature != "" {
		identifier, err := findFeature(db, definition.Feature)
		if err != nil {
			return 0, err
		}
		featureID = &identifier
	}
	identifier, err := coreid.NewUUIDv7()
	if err != nil {
		return 0, err
	}
	result := db.Exec(constants.InsertPermission, identifier, definition.Code, constants.Module, definition.Name, featureID)
	if result.Error != nil {
		return 0, result.Error
	}
	var stored registeredDefinition
	found := db.Raw(constants.FindPermission, definition.Code).Scan(&stored)
	if found.Error != nil {
		return 0, found.Error
	}
	if found.RowsAffected != 1 || stored.ID == "" || stored.Module != constants.Module || !sameBinding(stored.FeatureID, featureID) {
		return 0, registry.ErrConflict
	}
	return int(result.RowsAffected), nil
}

func sameBinding(actual *string, expected *string) bool {
	if actual == nil || expected == nil {
		return actual == nil && expected == nil
	}
	return *actual == *expected
}
