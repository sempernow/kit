// Package anet/badu/cmd provides an example usage of badu/authorize.net code base.
package main

import (
	anet "github.com/sempernow/kit/anet/badu"
)

//import anet "authorize_net"

func main() {

	//anet.ExampleAuthenticate()
	//anet.ExampleChargeCreditCard()
	//anet.ExampleCreateCustomerProfileFromTransaction()
	//anet.ExampleCreateCustomerProfileTransaction()
	//anet.ExampleGetMerchantDetails()

	anet.ExampleAcceptPayment()
	//... requires Accept.js (client)
}
