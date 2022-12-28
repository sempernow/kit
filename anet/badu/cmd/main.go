// Package anet/badu/cmd provides an example usage of badu/authorize.net code base.
package main

import (
	anet "gd9/prj3/kit/anet/badu"
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
