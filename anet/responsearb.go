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
