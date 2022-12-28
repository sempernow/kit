package anet

// ====================================================================
// FORK OF https://github.com/badu/authorize.net
// GNU General Public License v3.0
// https://github.com/badu/authorize.net/blob/master/LICENSE
// ====================================================================

// ===============
// Responses
// ===============

type Response struct {
	TransactionResponse *TransactionResponse `json:"transactionResponse,omitempty"`
	ARBSubID            string               `json:"subscriptionId,omitempty"`
	ARBProfile          *ARBProfile          `json:"profile,omitempty"`
	RefId               string               `json:"refId,omitempty"` // Transaction.ID
	Messages            *Messages            `json:"messages,omitempty"`
	ProfileResponse     *ProfileResponse     `json:"profileResponse,omitempty"`

	// Amalgum of TransactionResponse.Errors and Messages (former takes precedent)
	ErrorCode int
	ErrorText string

	// Method int `db:"method" json:"-"` // Auth (1), ARB (2), Charge (3), Void (4)
	// Mode   int `db:"mode" json:"-"`   // {Charge, Arb}{Bare, Accept, AcceptUI, Profile} (1-8)

}

// =======================
// Response constituents
// =======================

type TransactionResponse struct {
	ResponseCode   string       `json:"responseCode,omitempty"`
	AuthCode       string       `json:"authCode,omitempty"`
	AvsResultCode  string       `json:"avsResultCode,omitempty"`
	CvvResultCode  string       `json:"cvvResultCode,omitempty"`
	CavvResultCode string       `json:"cavvResultCode,omitempty"`
	TransId        string       `json:"transId,omitempty"`
	RefTransID     string       `json:"refTransID,omitempty"`
	TransHash      string       `json:"transHash,omitempty"`
	TestRequest    string       `json:"testRequest,omitempty"`
	AccountNumber  string       `json:"accountNumber,omitempty"`
	AccountType    string       `json:"accountType,omitempty"`
	Messages       []ErrMessage `json:"messages,omitempty"`
	TransHashSha2  string       `json:"transHashSha2,omitempty"`
	Profile        *Profile     `json:"profile,omitempty"`

	SupplementalDataQualificationIndicator int `json:"supplementalDataQualificationIndicator,omitempty"`

	NetworkTransId string  `json:"networkTransId,omitempty"`
	Errors         []Error `json:"errors,omitempty"` //min=0
}

type ProfileResponse struct {
	Messages *Messages `json:"messages,omitempty"`
}

type Messages struct {
	ResultCode string       `json:"resultCode,omitempty"` // "Ok" or "Error"
	Message    []ErrMessage `json:"message,omitempty"`
}
type ErrMessage struct {
	Code        ErrCodeEnum `json:"code"`
	Text        string      `json:"text"`
	Description string      `json:"description,omitempty"`
}

type Error struct {
	ErrorCode ErrCodeEnum `json:"errorCode,omitempty"`
	ErrorText string      `json:"errorText,omitempty"`
}
