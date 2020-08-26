package constants

const (
	HeaderContentTypeJson = "application/json"

	KeyPassword = "1234567890"

	UriKeyCreate   = "/keys"
	UriAccountInfo = "/auth/accounts/%v"      // format is /auth/accounts/{address}
	UriTransfer    = "/bank/accounts/%s/send" // format is /bank/accounts/{address}/transfers
	//UriTxs           = "/txs?action=send&sender=%s&recipient=%s"
	UriTxs         = "/txs/%s"
	UriTxBroadcast = "/txs"
	UriQueryTx     = "/txs?%s" //  txs?action=send&search_request_page=1&search_request_size=1

	// http status code
	StatusCodeOk       = 200
	StatusCodeConflict = 409

	//
	MockDefaultGas = "20000"
	MockDefaultFee = "5iris"
	Denom          = "iris"
)
