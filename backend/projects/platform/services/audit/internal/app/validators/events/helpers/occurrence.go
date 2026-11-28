package helpers

import client "grovehub/services/audit/internal/app/dto/producer"

func RegisteredOccurrence(data client.Activity) bool {
	switch data.Source {
	case "identity":
		if data.AssetType != "ACCOUNT" || data.AssetVersion != 0 || data.PersonID != "" {
			return false
		}
		if data.Operation == "LOGIN_FAILED" {
			return data.ActorID == ""
		}
		if data.ActorID != data.AssetID {
			return false
		}
		switch data.Operation {
		case "LOGIN_SUCCEEDED", "PASSWORD_RESET", "PASSWORD_CHANGED", "LOGOUT", "LOGOUT_ALL", "SESSION_REVOKED", "REAUTHENTICATED", "REAUTHENTICATION_FAILED":
			return true
		}
	case "foundation":
		if data.AssetType != "PERSON" || data.AssetVersion < 1 || data.ActorID == "" || data.PersonID != data.AssetID {
			return false
		}
		switch data.Operation {
		case "CREATED", "UPDATED", "DEACTIVATED":
			return true
		}
	case "integration":
		if data.AssetType != "CONNECTION" || data.AssetVersion != 0 || !ValidID(data.PersonID) {
			return false
		}
		switch data.Operation {
		case "CREATED", "ACTIVATED", "DISCONNECTED", "TESTED", "SYSTEM_MAIL_SENDER_CONFIGURED":
			return true
		}
	}

	return false
}
