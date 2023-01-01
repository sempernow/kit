// Package anet is a modified subset of the badu/authorize.net code base.
package anet

import "time"

// ====================================================================
// FORK OF https://github.com/badu/authorize.net
// GNU General Public License v3.0
// https://github.com/badu/authorize.net/blob/master/LICENSE
// ====================================================================

// =====================================
// Types common to Request and Response
// =====================================

type Profile struct {
	CustomerProfileId string          `json:"customerProfileId,omitempty"`
	PaymentProfile    *PaymentProfile `json:"paymentProfile,omitempty"`
}
type PaymentProfile struct {
	PaymentProfileId string `json:"paymentProfileId,omitempty"`
}

type ARBProfile struct {
	CustomerProfileID        string `json:"customerProfileId,omitempty"`
	CustomerPaymentProfileID string `json:"customerPaymentProfileId,omitempty"`
	CustomerAddressID        string `json:"customerAddressId,omitempty"`
}

type ANetApiRequest struct {
	MerchantAuthentication MerchantAuthentication `json:"merchantAuthentication"` // Contains merchant authentication information.
	ClientId               string                 `json:"clientId,omitempty"`
	RefId                  string                 `json:"refId,omitempty"` // Merchant-assigned reference ID for the request. If included in the request, this value is included in the response. This feature might be especially useful for multi-threaded applications.
}
type ANetApiResponse struct {
	RefId        string   `json:"refId,omitempty"` // Merchant-assigned reference ID for the request. If included in the request, this value will be included in the response. This feature might be especially useful for multi-threaded applications.
	Messages     Messages `json:"messages"`
	SessionToken string   `json:"sessionToken,omitempty"`
}

type Paging struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}
type Sorting struct {
	OrderBy         string `json:"orderBy"`
	OrderDescending bool   `json:"orderDescending"`
}

const (
	// RFC3339 a subset of the ISO8601 timestamp format. e.g 2014-04-29T18:30:38Z
	ISO8601TimeFormat = "2006-01-02T15:04:05Z"
	// same as above, but no ”Z”
	ISO8601NoZTimeFormat = "2006-01-02T15:04:05"
	// 2018-12-27T11:28:57.467
	ISO8601TimeNano = "2006-01-02T15:04:05.999999999"
)

type Date struct {
	time.Time
}

func NowDate() Date {
	return Date{time.Now()}
}
func (t Date) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.Time.UTC().Format("2006-01-02") + `"`), nil
}
func (t *Date) UnmarshalJSON(data []byte) error {
	var err error
	if t.Time, err = time.Parse(ISO8601TimeFormat, string(data[1:len(data)-1])); err != nil {
		return err
	}
	return nil
}
