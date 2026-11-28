package helpers

import (
	"strings"
)

func generatedRegistrationName(group routeGroup) string {
	tier := ""
	switch group.tier {
	case "public":
		tier = "Public"
	case "protected":
		tier = "Protected"
	case "restricted":
		tier = "Restricted"
	}
	return "register" + tier + strings.TrimSuffix(group.name, "Controller") + "Routes"
}

func groupImportAlias(group routeGroup) string {
	pkg:=group.tier
	if pkg=="" {pkg="controller"}
	if group.resource!="" {return group.resource+pkg}
	return pkg
}
