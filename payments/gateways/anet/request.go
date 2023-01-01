package anet

/******************************************************************
	Authorize.net API is not based on REST.
	It accepts XML or JSON content types, but is natively XML.
	All requests (POST) are serviced at one endpoint.
	(One for production and one for sandbox.)
	Its functions and their resulting response body vary
	per request-payload content and order.

	https://developer.authorize.net/api/reference/index.html
******************************************************************/

// ====================================================================
// FORK OF https://github.com/badu/authorize.net
// GNU General Public License v3.0
// https://github.com/badu/authorize.net/blob/master/LICENSE
// ====================================================================

// ============================================================
// Request* structs adapt api-function requests to their type.
// ============================================================

type RequestReadinessTest struct {
	ReadinessTest *CreateTransactionRequest `json:"authenticateTestRequest,omitempty"`
}
type RequestChargeCustomerProfile struct {
	ChargeCustomerProfile *CreateTransactionRequest `json:"createTransactionRequest,omitempty"`
}
type RequestAcceptPayment struct {
	AcceptPaymentTransaction *CreateTransactionRequest `json:"createTransactionRequest,omitempty"`
}
type RequestChargePerCardData struct {
	ChargePerCardData *CreateTransactionRequest `json:"createTransactionRequest,omitempty"`
}
type RequestVoidTransaction struct {
	VoidTransaction *CreateTransactionRequest `json:"createTransactionRequest,omitempty"`
}

// =======================
// Request constituents
// =======================

type CreateTransactionRequest struct {
	MerchantAuthentication *MerchantAuthentication `json:"merchantAuthentication,omitempty"`
	RefId                  string                  `json:"refId,omitempty"` // Transaction.ID
	TransactionRequest     *TransactionRequest     `json:"transactionRequest,omitempty"`
}
type MerchantAuthentication struct {
	Name           string `json:"name,omitempty"`
	TransactionKey string `json:"transactionKey,omitempty"`
}
type TransactionRequest struct {
	TransactionType TransactionTypeEnum `json:"transactionType,omitempty"`
	RefTransId      string              `json:"refTransId,omitempty"`
	Amount          string              `json:"amount,omitempty"`
	Profile         *Profile            `json:"profile,omitempty"`
	Payment         *Payment            `json:"payment,omitempty"`
	Customer        *Customer           `json:"customer,omitempty"` // Email
	BillTo          *Customer           `json:"billTo,omitempty"`   // Zip @ AcceptPayment
	ShipTo          *Customer           `json:"shipTo,omitempty"`   // Zip @ ChargeProfile
	CustomerIP      string              `json:"customerIP,omitempty"`
}

type Payment struct {
	CreditCard  *CreditCard  `json:"creditCard,omitempty"`
	BankAccount *BankAccount `json:"bankAccount,omitempty"`
	//TrackData          *CreditCardTrack    `json:"trackData,omitempty"`
	//EncryptedTrackData *EncryptedTrackData `json:"encryptedTrackData,omitempty"`
	//PayPal             *PayPal             `json:"payPal,omitempty"`
	OpaqueData *OpaqueData `json:"opaqueData,omitempty"`
	//Emv                *PaymentEmv         `json:"emv,omitempty"`
	DataSource string `json:"dataSource,omitempty"`
}
type Customer struct {
	FirstName string           `json:"firstName,omitempty"` // Required @ BillTo{} @ EU processor
	LastName  string           `json:"lastName,omitempty"`  // Required @ BillTo{} @ EU processor
	Company   string           `json:"company,omitempty"`
	Type      CustomerTypeEnum `json:"type,omitempty"`  // "individual" | "business"
	ID        string           `json:"id,omitempty"`    // x.PayerID; max=20
	Email     string           `json:"email,omitempty"` // Required @ Customer{} @ EU processor
	Address   string           `json:"address,omitempty"`
	City      string           `json:"city,omitempty"`  // Required @ BillTo{} @ EU processor
	State     string           `json:"state,omitempty"` // Required @ BillTo{} @ EU processor
	Zip       string           `json:"zip,omitempty"`
	Country   string           `json:"country,omitempty"`
}

type OpaqueData struct {
	DataDescriptor string `json:"dataDescriptor"`    // use "COMMON.VCO.ONLINE.PAYMENT" for Visa checkout transactions
	DataValue      string `json:"dataValue"`         // Base64 encoded data that contains encrypted payment data.
	DataKey        string `json:"dataKey,omitempty"` // The encryption key used to encrypt the payment data.
	//ExpirationTimeStamp time.Time `json:"expirationTimeStamp,omitempty"`
}
type CreditCard struct {
	CardNumber         string `json:"cardNumber"`
	ExpirationDate     string `json:"expirationDate"`
	CardCode           string `json:"cardCode,omitempty"`
	IsPaymentToken     bool   `json:"isPaymentToken,omitempty"`
	Cryptogram         string `json:"cryptogram,omitempty"`
	TokenRequestorName string `json:"tokenRequestorName,omitempty"`
	TokenRequestorId   string `json:"tokenRequestorId,omitempty"`
	TokenRequestorEci  string `json:"tokenRequestorEci,omitempty"`
}
type BankAccount struct {
	AccountType   BankAccountTypeEnum `json:"accountType,omitempty"`
	RoutingNumber string              `json:"routingNumber"`
	AccountNumber string              `json:"accountNumber"`
	NameOnAccount string              `json:"nameOnAccount"`
	EcheckType    EcheckTypeEnum      `json:"echeckType,omitempty"`
	BankName      string              `json:"bankName,omitempty"`
	CheckNumber   string              `json:"checkNumber,omitempty"`
}
