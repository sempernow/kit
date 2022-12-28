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

type RequestCreateSubscription struct {
	ARBCreateSubscriptionRequest *ARBCreateSubscriptionRequest `json:"ARBCreateSubscriptionRequest,omitempty"`
}

// =======================
// Request constituents
// =======================

type ARBCreateSubscriptionRequest struct {
	MerchantAuthentication *MerchantAuthentication `json:"merchantAuthentication,omitempty"`
	RefId                  string                  `json:"refId,omitempty"` // Transaction.ID
	Subscription           *Subscription           `json:"subscription,omitempty"`
}
type Subscription struct {
	Name            string           `json:"name,omitempty"`
	PaymentSchedule *PaymentSchedule `json:"paymentSchedule,omitempty"`
	Amount          float64          `json:"amount,omitempty"`
	TrialAmount     float64          `json:"trialAmount"`
	Payment         *Payment         `json:"payment,omitempty"`
	BillTo          *Customer        `json:"billTo,omitempty"`
	Profile         *Profile         `json:"profile,omitempty"`
}
type PaymentSchedule struct {
	Interval         *Interval `json:"interval,omitempty"`
	StartDate        Date      `json:"startDate,omitempty"`
	TotalOccurrences int       `json:"totalOccurrences"`
	TrialOccurrences int       `json:"trialOccurrences"`
}
type Interval struct {
	Length int                     `json:"length"`
	Unit   ARBSubscriptionUnitEnum `json:"unit"`
}

type CustomerProfileEx struct {
	CustomerProfileBase
	CustomerProfileId string `json:"customerProfileId,omitempty"`
}
