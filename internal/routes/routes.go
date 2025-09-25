package routes

// API endpoint path templates for registration with mux (use {name}).
const (
	RouteCreateKey = "/transit/keys/{name}"
	RouteEncrypt   = "/transit/encrypt/{name}"
	RouteDecrypt   = "/transit/decrypt/{name}"

	// Names for mux routes (used for URL building)
	RouteNameCreateKey = "createKey"
	RouteNameEncrypt   = "encrypt"
	RouteNameDecrypt   = "decrypt"
)
