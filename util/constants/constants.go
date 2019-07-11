package constants

const (
	HeaderContentTypeJson = "application/json"

	KeyPassword   = "1234567890"

	UriKeyCreate     = "/keys"
	UriAccountInfo   = "/bank/accounts/%v"           // format is /auth/accounts/{address}
	UriTransfer      = "/bank/accounts/%s/send" // format is /bank/accounts/{address}/transfers
	//UriTxs           = "/txs?action=send&sender=%s&recipient=%s"
	UriTxs           = "/txs/%s"
	UriTxBroadcast   = "/tx/broadcast"

	// http status code
	StatusCodeOk       = 200
	StatusCodeConflict = 409

	//
	MockDefaultGas     = "20000"
	MockDefaultFee     = "5iris"
	Denom              = "iris"
)
