package sat

const (
	// ACCESS_TOKEN_URL is constant full url for oauth
	ACCESS_TOKEN_URL = "https://accounts.tokopedia.com/token"

	// PLAYGROUND_SAT_BASE_URL is constant of base URL SAT Server Playground
	// Production usage should use different base URL, please override it using this function WithSatBaseURL
	PLAYGROUND_SAT_BASE_URL = "https://b2b-playground.tokopedia.com/api"

	// PING_PATH is constant of ping endpoint
	PING_PATH = "/ping"
	// ACCOUNT_PATH is constant of account endpoint
	ACCOUNT_PATH = "/v2/account"
	// INQUIRY_PATH is constant of inquiry endpoint
	INQUIRY_PATH = "/v2/inquiry"
	// CHECKOUT_PATH is constant of checkout endpoint
	CHECKOUT_PATH = "/v2/order"
	// CHECK_STATUS_PATH is constant of check status endpoint
	CHECK_STATUS_PATH = "/v2/order/%s"
	// PRODUCT_LIST_PATH is constant of product list endpoint
	PRODUCT_LIST_PATH = "/v2/product-list"

	// SIGNATURE_HEADER_KEY is the key name used as header http of digital signature
	SIGNATURE_HEADER_KEY = "signature"
	// SAT_SDK_VERSION is current sdk version
	SAT_SDK_VERSION = "golang-sat@v1.0.0"
)
