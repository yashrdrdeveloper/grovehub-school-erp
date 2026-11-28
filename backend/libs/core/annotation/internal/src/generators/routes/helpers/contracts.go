package helpers

type route struct {
	method string
	verb string
	path string
	limit int
}

type routeGroup struct {
	name   string
	tier   string
	resource string
	routes []route
}

var verbs = map[string]bool{
	"GET": true, "POST": true, "PUT": true,
	"PATCH": true, "DELETE": true,
}

var accessTiers = []string{"public", "protected", "restricted"}
