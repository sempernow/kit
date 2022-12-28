package authorize_net

type AuthenticateTestRequest struct {
	ANetApiRequest `json:"authenticateTestRequest"`
}

type LogoutResponse struct {
	ANetApiResponse
}

type API struct {
	MerchantName   string
	MerchantSecret string
	TestMode       bool
	URL            string
}

func SetAPI(name string, key string, mode string) API {
	api := API{}
	api.MerchantName = name
	api.MerchantSecret = key
	if mode == "live" {
		api.TestMode = true
		api.URL = EndpointLive
	} else {
		api.TestMode = false
		api.URL = EndpointSandbox
	}
	return api
}
