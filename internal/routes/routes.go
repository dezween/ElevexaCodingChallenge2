package routes

// API endpoint path templates for registration with mux (use {name}).
//
// Methods:
//
//	POST RouteCreateKey   - Create a new Kyber key pair
//	POST RouteEncrypt     - Encrypt data with Kyber
//	POST RouteDecrypt     - Decrypt data with Kyber
const (
	// POST: Create a new Kyber key pair
	RouteCreateKey = "/transit/keys/{name}"
	// POST: Encrypt data with Kyber
	RouteEncrypt = "/transit/encrypt/{name}"
	// POST: Decrypt data with Kyber
	RouteDecrypt = "/transit/decrypt/{name}"

	// Names for mux routes (used for URL building)
	RouteNameCreateKey = "createKey"
	RouteNameEncrypt   = "encrypt"
	RouteNameDecrypt   = "decrypt"
)
