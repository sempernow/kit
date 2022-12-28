// Package badu/authorize.net extracts.
package authorize_net

import (
	"context"
	"fmt"
	"gd9/prj3/kit/convert"
	"log"
	"os"
)

var (
	cardExample        = "4007000000027"
	cardExpirationDate = "11/22"
	cardCode           = "999"

	ctx    = context.Background()
	client = NewAPIClient(nil, log.New(os.Stdout, "INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile),
	)
)
var (
	apiName = os.Getenv("ANET_API_NAME")
	apiKey  = os.Getenv("ANET_API_KEY")
	mode    = os.Getenv("ANET_API_MODE")
)
var (
	api   = SetAPI(apiName, apiKey, mode)
	endpt = api.URL
)

func ExampleAuthenticate() {

	payload := AuthenticateTestRequest{
		ANetApiRequest: ANetApiRequest{
			MerchantAuthentication: MerchantAuthentication{
				Name:           apiName,
				TransactionKey: apiKey,
			},
		},
	}
	//fmt.Println(convert.PrettyPrint(payload))
	// {
	//         "authenticateTestRequest": {
	//                 "merchantAuthentication": {
	//                         "name": "7uTw62hP9b",
	//                         "transactionKey": "5cy8Izk48nFf5N9a"
	//                 }
	//         }
	// }

	var resp ANetApiResponse

	if err := client.Send(ctx, endpt, &payload, &resp); err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}
	if resp.Messages.ResultCode == "Error" {
		fmt.Printf("%#v", resp)
	}
	//fmt.Printf("%+v", resp)
	// Output: {RefId: Messages:{ResultCode:Ok Messages:[{Code:Successful. Text:Successful. Description:}]} SessionToken:}
	fmt.Println(convert.PrettyPrint(resp))
	// {
	//         "messages": {
	//                 "resultCode": "Ok",
	//                 "message": [
	//                         {
	//                                 "code": 1,
	//                                 "text": "Successful."
	//                         }
	//                 ]
	//         }
	// }
}

// Note that credit card information and bank account information are mutually exclusive, so you should not submit both.
func ExampleChargeCreditCard() {

	payload := CreateTransactionRequest{
		CreateTransactionRequest: TransactionRequest{
			MerchantAuthentication: MerchantAuthentication{
				Name:           apiName,
				TransactionKey: apiKey,
			},
			Transaction: Transaction{
				TransactionType: AuthCaptureTransaction,
				Amount:          11.3,
				Payment: &Payment{
					CreditCard: &CreditCard{
						CardNumber:     cardExample,
						ExpirationDate: cardExpirationDate,
						CardCode:       cardCode,
					},
				},
				BillTo: &CustomerAddress{
					NameAndAddress: NameAndAddress{
						Zip: "32457",
					},
				},
			},
		},
	}

	var resp CreateTransactionResponse

	if err := client.Send(ctx, endpt, &payload, &resp); err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}

	if resp.Messages.ResultCode == "Error" {
		//fmt.Printf("%#v", resp)
		fmt.Println(convert.PrettyPrint(resp))
	}
	fmt.Printf("%+v", convert.PrettyPrint(resp.TransactionResponse))
	// {
	//         "responseCode": "1",
	//         "authCode": "PF72MJ",
	//         "avsResultCode": "Y",
	//         "cvvResultCode": "P",
	//         "cavvResultCode": "2",
	//         "transId": "40072455021",
	//         "testRequest": "0",
	//         "accountNumber": "XXXX0027",
	//         "accountType": "Visa",
	//         "messages": [
	//                 {
	//                         "code": 159,
	//                         "text": "",
	//                         "description": "This transaction has been approved."
	//                 }
	//         ],
	//         "networkTransId": "GUZ18GIRQHVK6TWRWSN396R"
	// }

}

// AcceptPayment scheme
func ExampleAcceptPayment() {

	payload := CreateTransactionRequest{
		CreateTransactionRequest: TransactionRequest{
			MerchantAuthentication: MerchantAuthentication{
				Name:           apiName,
				TransactionKey: apiKey,
			},
			Transaction: Transaction{
				TransactionType: AuthCaptureTransaction,
				Amount:          11,
				Payment: &Payment{
					OpaqueData: &OpaqueData{
						DataDescriptor: "dataDescriptor",
						DataValue:      "dataValue",
					},
				},
				BillTo: &CustomerAddress{
					NameAndAddress: NameAndAddress{
						Zip: "32457",
					},
				},
			},
		},
	}

	var resp CreateTransactionResponse

	if err := client.Send(ctx, endpt, &payload, &resp); err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}

	if resp.Messages.ResultCode == "Error" {
		fmt.Println(convert.PrettyPrint(resp))
	}
	fmt.Printf("%+v", convert.PrettyPrint(resp.TransactionResponse))
	// {
	//         "responseCode": "1",
	//         "authCode": "PF72MJ",
	//         "avsResultCode": "Y",
	//         "cvvResultCode": "P",
	//         "cavvResultCode": "2",
	//         "transId": "40072455021",
	//         "testRequest": "0",
	//         "accountNumber": "XXXX0027",
	//         "accountType": "Visa",
	//         "messages": [
	//                 {
	//                         "code": 159,
	//                         "text": "",
	//                         "description": "This transaction has been approved."
	//                 }
	//         ],
	//         "networkTransId": "GUZ18GIRQHVK6TWRWSN396R"
	// }

}

func ExampleRefund() {

	voidPayload := CreateTransactionRequest{
		CreateTransactionRequest: TransactionRequest{
			MerchantAuthentication: MerchantAuthentication{
				Name:           apiName,
				TransactionKey: apiKey,
			},
			Transaction: Transaction{
				TransactionType: RefundTransaction,
				Amount:          123123,
				Payment: &Payment{
					CreditCard: &CreditCard{
						CardNumber:     cardExample,
						ExpirationDate: cardExpirationDate,
						CardCode:       cardCode,
					},
				},
			},
		},
	}

	refundRequest, err := client.prepareRequest(context.Background(), endpt, &voidPayload)
	if err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}
	var resp CreateTransactionResponse
	if err := client.Do(refundRequest, &resp, false); err != nil {
		fmt.Printf("Error : %v", err)
		return
	}
	if resp.Messages.ResultCode == "Error" {
		fmt.Printf("%#v", resp)
	}
	fmt.Printf("%+v", resp.TransactionResponse.Messages)
	// Output: [{Code:This transaction has been approved. Text: Description:This transaction has been approved.}]

}

func ExampleVoid() {

	amount := 10.31
	payload := CreateTransactionRequest{
		CreateTransactionRequest: TransactionRequest{
			MerchantAuthentication: MerchantAuthentication{
				Name:           apiName,
				TransactionKey: apiKey,
			},
			Transaction: Transaction{
				TransactionType: AuthCaptureTransaction,
				Amount:          amount,
				Payment: &Payment{
					CreditCard: &CreditCard{
						CardNumber:     cardExample,
						ExpirationDate: cardExpirationDate,
						CardCode:       cardCode,
					},
				},
			},
		},
	}

	req, err := client.prepareRequest(context.Background(), endpt, &payload)
	if err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}

	var resp CreateTransactionResponse
	if err := client.Do(req, &resp, false); err != nil {
		fmt.Printf("Error : %v", err)
		return
	}

	voidPayload := CreateTransactionRequest{
		CreateTransactionRequest: TransactionRequest{
			MerchantAuthentication: MerchantAuthentication{
				Name:           apiName,
				TransactionKey: apiKey,
			},
			Transaction: Transaction{
				TransactionType: VoidTransaction,
				Amount:          amount,
				Payment: &Payment{
					CreditCard: &CreditCard{
						CardNumber:     cardExample,
						ExpirationDate: cardExpirationDate,
						CardCode:       cardCode,
					},
				},
				RefTransId: resp.TransactionResponse.TransId,
			},
		},
	}

	voidRequest, err := client.prepareRequest(context.Background(), endpt, &voidPayload)
	if err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}
	var voidResponse CreateTransactionResponse
	if err := client.Do(voidRequest, &voidResponse, false); err != nil {
		fmt.Printf("Error : %v", err)
		return
	}
	if resp.Messages.ResultCode == "Error" {
		fmt.Printf("%#v", resp)
	}
	fmt.Printf("%+v", voidResponse.TransactionResponse.Messages)
	// Output: [{Code:This transaction has been approved. Text: Description:This transaction has been approved.}]

}

// Profile creation
func ExampleCreateCustomerProfile() {

	payload := CreateCustomerProfileRequest{
		Payload: CreateCustomerProfilePayload{
			ANetApiRequest: ANetApiRequest{
				MerchantAuthentication: MerchantAuthentication{
					Name:           apiName,
					TransactionKey: apiKey,
				},
			},
			Profile: &CustomerProfile{
				CustomerProfileBase: CustomerProfileBase{
					MerchantCustomerId: apiName,
				},
				PaymentProfiles: []CustomerPaymentProfile{
					{
						CustomerPaymentProfileBase: CustomerPaymentProfileBase{
							CustomerType: Individual,
							BillTo: &CustomerAddress{
								NameAndAddress: NameAndAddress{
									FirstName: "Popescu",
									LastName:  "Grigore",
									Country:   "Romania",
								},
								PhoneNumber: "030030020",
								Email:       "b@gmail.com",
							},
						},
						Payment: &Payment{
							CreditCard: &CreditCard{
								CardNumber:     cardExample,
								ExpirationDate: cardExpirationDate,
								CardCode:       cardCode,
								IsPaymentToken: false,
							},
						},
						DefaultPaymentProfile: true,
					},
				},
				ProfileType: CustomerProfileRegular,
			},
			ValidationMode: ValidationModeTestMode,
		},
	}

	req, err := client.prepareRequest(context.Background(), endpt, &payload)
	if err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}

	var resp CreateCustomerProfileResponse
	if err := client.Do(req, &resp, false); err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}
	if resp.Messages.ResultCode == "Error" {
		fmt.Printf("%#v", resp)
	}
	//fmt.Printf("%+v", resp)
	fmt.Println(convert.PrettyPrint(resp))
}

// Create Customer Profile from Transaction ID
func ExampleCreateCustomerProfileFromTransaction() {

	payload := CreateCustomerProfileFromTransactionRequest{
		Payload: CreateCustomerProfileFromTransactionPayload{
			ANetApiRequest: ANetApiRequest{
				MerchantAuthentication: MerchantAuthentication{
					Name:           apiName,
					TransactionKey: apiKey,
				},
			},
			Customer: &CustomerProfileBase{
				MerchantCustomerId: "19328",
				//Description:        "1233",
				//Email:              "123123",
			},
			TransId:                "40072455021",
			DefaultPaymentProfile:  true,
			DefaultShippingAddress: true,
			//ProfileType:            CustomerProfileGuest,
			//... sans defaults to regular, which is what we want.
		},
	}

	var resp CreateCustomerProfileFromTransactionResponse

	if err := client.Send(ctx, endpt, &payload, &resp); err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}

	if resp.Messages.ResultCode == "Error" {
		fmt.Printf("%#v", resp)
	}
	//fmt.Printf("%#v", resp)
	// Output: authorize_net.CreateCustomerProfileFromTransactionResponse{ANetApiResponse:authorize_net.ANetApiResponse{RefId:"", Messages:authorize_net.Messages{ResultCode:"Ok", Messages:[]authorize_net.ErrMessage{authorize_net.ErrMessage{Code:0x1, Text:"Successful.", Description:""}}}, SessionToken:""}, CustomerProfileID:"501099822", CustomerPaymentProfileIDList:[]string{"501838004"}, CustomerShippingAddressIDList:[]string{}, ValidationDirectResponseList:[]string{}}
	fmt.Println(convert.PrettyPrint(resp))
	// {
	//         "messages": {
	//                 "resultCode": "Ok",
	//                 "message": [
	//                         {
	//                                 "code": 1,
	//                                 "text": "Successful."
	//                         }
	//                 ]
	//         },
	//         "customerProfileId": "501099822",
	//         "customerPaymentProfileIdList": [
	//                 "501838004"
	//         ],
	//         "customerShippingAddressIdList": [],
	//         "validationDirectResponseList": []
	// }
}

// Transaction per Profile IDs (Customer/Payment IDs)
// DOES NOT fit the current API of Authorize.net
// FAILs to return TransactionResponse.ResponseCode
func ExampleCreateCustomerProfileTransaction() {

	payload := CreateCustomerProfileTransactionRequest{
		Payload: CreateCustomerProfileTransactionPayload{
			ANetApiRequest: ANetApiRequest{
				MerchantAuthentication: MerchantAuthentication{
					Name:           apiName,
					TransactionKey: apiKey,
				},
			},
			Transaction: &ProfileTransaction{
				ProfileTransAuthCapture: &ProfileTransOrder{
					ProfileTransAmount: ProfileTransAmount{
						Amount: 11,
						Tax: &ExtendedAmount{
							Amount:      11,
							Name:        "foo",
							Description: "bar",
						},
						Shipping: &ExtendedAmount{
							Amount:      11,
							Name:        "",
							Description: "",
						},
						Duty: &ExtendedAmount{
							Amount:      11,
							Name:        "",
							Description: "",
						},
						LineItems: nil,
					},
					CustomerProfileId:        "501099822",
					CustomerPaymentProfileId: "501838004",
					Order: &OrderEx{
						Order: Order{
							InvoiceNumber:                  "INV001",
							Description:                    "",
							DiscountAmount:                 10,
							TaxIsAfterDiscount:             true,
							TotalTaxTypeCode:               "",
							PurchaserVATRegistrationNumber: "",
							MerchantVATRegistrationNumber:  "",
							VatInvoiceReferenceNumber:      "",
							PurchaserCode:                  "1234",
							SummaryCommodityCode:           "",
							PurchaseOrderDateUTC:           Now(),
							SupplierOrderReference:         "",
							AuthorizedContactName:          "",
							CardAcceptorRefNumber:          "",
							AmexDataTAA1:                   "",
							AmexDataTAA2:                   "",
							AmexDataTAA3:                   "",
							AmexDataTAA4:                   "",
						},
						PurchaseOrderNumber: "",
					},
					TaxExempt:        true,
					RecurringBilling: false,
					CardCode:         cardCode,
					//SplitTenderId:    "",
					ProcessingOptions: &ProcessingOptions{
						IsFirstRecurringPayment: false,
						IsFirstSubsequentAuth:   false,
						IsSubsequentAuth:        false,
						IsStoredCredentials:     false,
					},
				},
			},
		},
	}

	var resp CreateCustomerProfileTransactionResponse

	if err := client.Send(ctx, endpt, &payload, &resp); err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}

	if resp.Messages.ResultCode == "Error" {
		fmt.Printf("%#v", resp)
	}
	//fmt.Printf("%#v", resp)
	// Output: authorize_net.CreateCustomerProfileTransactionResponse{ANetApiResponse:authorize_net.ANetApiResponse{RefId:"", Messages:authorize_net.Messages{ResultCode:"Ok", Messages:[]authorize_net.ErrMessage{authorize_net.ErrMessage{Code:0x1, Text:"Successful.", Description:""}}}, SessionToken:""}, TransactionResponse:(*authorize_net.TransactionResponse)(nil), DirectResponse:"1,1,1,This transaction has been approved.,N0SU5L,Y,60126781190,INV001,,10.00,CC,auth_capture,,John,Smith,,,,,,,,,,,,,,,,,,10.00,10.00,10.00,TRUE,,,P,2,,,,,,,,,,,XXXX1111,Visa,,,,,,,,,,,,,,,,,"}
	fmt.Println(convert.PrettyPrint(resp))
	// {
	//         "messages": {
	//                 "resultCode": "Ok",
	//                 "message": [
	//                         {
	//                                 "code": 1,
	//                                 "text": "Successful."
	//                         }
	//                 ]
	//         },
	//         "directResponse": "1,1,1,This transaction has been approved.,O350KJ,Y,40072456492,INV001,,11.00,CC,auth_capture,19328,,,,,,,32457,,,,123123,,,,,,,,,11.00,11.00,11.00,TRUE,,,P,2,,,,,,,,,,,XXXX0027,Visa,,,,,,,BDQOVFRAHR3YRCH9RA0DDPQ,,,,,,,,,,"
	// }
}

func ExampleGetCustomerProfile() {

	payload := GetCustomerProfileRequest{
		Payload: GetCustomerProfilePayload{
			ANetApiRequest: ANetApiRequest{
				MerchantAuthentication: MerchantAuthentication{
					Name:           apiName,
					TransactionKey: apiKey,
				},
			},
			CustomerProfileId: "1920441543",
			IncludeIssuerInfo: true,
		},
	}
	req, err := client.prepareRequest(context.Background(), endpt, &payload)
	if err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}

	var resp GetCustomerProfileResponse
	if err := client.Do(req, &resp, false); err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}
	if resp.Messages.ResultCode == "Error" {
		fmt.Printf("%#v", resp)
	}
	fmt.Printf("%s\n", resp.Profile.CustomerProfileId)
	for _, pp := range resp.Profile.PaymentProfiles {
		fmt.Printf("%s", pp.CustomerPaymentProfileId)
	}
	// Output: 1920441543
	// 1833416626
}

func ExampleUpdateCustomerProfile() {

	payload := UpdateCustomerProfileRequest{
		Payload: UpdateCustomerProfilePayload{
			ANetApiRequest: ANetApiRequest{
				MerchantAuthentication: MerchantAuthentication{
					Name:           apiName,
					TransactionKey: apiKey,
				},
			},
			Profile: &CustomerProfileInfoEx{
				CustomerProfileEx: CustomerProfileEx{
					CustomerProfileBase: CustomerProfileBase{
						MerchantCustomerId: "custId123",
						Description:        "description",
						Email:              "some@invalid@email",
					},
					CustomerProfileId: "38157432",
				},
				ProfileType: CustomerProfileRegular,
			},
		},
	}
	req, err := client.prepareRequest(context.Background(), endpt, &payload)
	if err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}

	var resp UpdateCustomerProfileResponse
	if err := client.Do(req, &resp, false); err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}
	if resp.Messages.ResultCode == "Error" {
		fmt.Printf("%#v", resp)
	}
	fmt.Printf("%#v", resp)
	// Output: authorize_net.UpdateCustomerProfileResponse{ANetApiResponse:authorize_net.ANetApiResponse{RefId:"", Messages:authorize_net.Messages{ResultCode:"Ok", Messages:[]authorize_net.ErrMessage{authorize_net.ErrMessage{Code:0x1, Text:"Successful.", Description:""}}}, SessionToken:""}}

}
func ExampleDeleteCustomerProfile() {

	payload := DeleteCustomerProfileRequest{
		Payload: DeleteCustomerProfilePayload{
			ANetApiRequest: ANetApiRequest{
				MerchantAuthentication: MerchantAuthentication{
					Name:           apiName,
					TransactionKey: apiKey,
				},
			},
			CustomerProfileId: "38157410",
		},
	}
	req, err := client.prepareRequest(context.Background(), endpt, &payload)
	if err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}

	var resp DeleteCustomerProfileResponse
	if err := client.Do(req, &resp, false); err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}
	if resp.Messages.ResultCode == "Error" {
		fmt.Printf("%#v", resp)
	}
	fmt.Printf("%#v\n", resp)
	// Output: authorize_net.DeleteCustomerProfileResponse{ANetApiResponse:authorize_net.ANetApiResponse{RefId:"", Messages:authorize_net.Messages{ResultCode:"Ok", Messages:[]authorize_net.ErrMessage{authorize_net.ErrMessage{Code:0x1, Text:"Successful.", Description:""}}}, SessionToken:""}}
}

func ExampleGetCustomerPaymentProfile() {

	payload := GetCustomerPaymentProfileRequest{
		Payload: GetCustomerPaymentProfilePayload{
			ANetApiRequest: ANetApiRequest{
				MerchantAuthentication: MerchantAuthentication{
					Name:           apiName,
					TransactionKey: apiKey,
				},
			},
			CustomerProfileId:        "1920441543",
			CustomerPaymentProfileId: "1833416626",
			IncludeIssuerInfo:        true,
		},
	}
	req, err := client.prepareRequest(context.Background(), endpt, &payload)
	if err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}

	var resp GetCustomerPaymentProfileResponse
	if err := client.Do(req, &resp, false); err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}
	if resp.Messages.ResultCode == "Error" {
		fmt.Printf("%#v", resp)
	}
	fmt.Printf("%s\n", resp.PaymentProfile.CustomerProfileId)
	fmt.Printf("%s\n", resp.PaymentProfile.CustomerPaymentProfileId)
	// Output: 1920441543
	// 1833416626
}

func ExampleGetCustomerPaymentProfileList() {

	payload := GetCustomerPaymentProfileListRequest{
		Payload: GetCustomerPaymentProfileListPayload{
			ANetApiRequest: ANetApiRequest{
				MerchantAuthentication: MerchantAuthentication{
					Name:           apiName,
					TransactionKey: apiKey,
				},
			},
			SearchType: "cardsExpiringInMonth", // always the same
			Month:      "2020-12",
			Sorting: Sorting{
				OrderBy:         "id",
				OrderDescending: false,
			},
			Paging: Paging{
				Limit:  10,
				Offset: 1,
			},
		},
	}
	req, err := client.prepareRequest(context.Background(), endpt, &payload)
	if err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}

	var resp GetCustomerPaymentProfileListResponse
	if err := client.Do(req, &resp, false); err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}
	if resp.Messages.ResultCode == "Error" {
		fmt.Printf("%#v", resp)
	}
	fmt.Printf("%d\n", resp.TotalNumInResultSet)
	for _, pp := range resp.PaymentProfiles {
		fmt.Printf("%d\n", pp.CustomerProfileId)
	}

	// Output:118569
	// 37821321
	// 38155405
	// 38155971
	// 38157410
	// 38157423
	// 38157432
	// 38157450
	// 38157457
	// 38157567
	// 38157570
}

func ExampleGetCustomerProfileIds() {

	payload := GetCustomerProfileIdsRequest{
		ANetApiRequest: ANetApiRequest{
			MerchantAuthentication: MerchantAuthentication{
				Name:           apiName,
				TransactionKey: apiKey,
			},
		},
	}
	req, err := client.prepareRequest(context.Background(), endpt, &payload)
	if err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}

	var resp GetCustomerProfileIdsResponse
	if err := client.Do(req, &resp, false); err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}
	if resp.Messages.ResultCode == "Error" {
		fmt.Printf("%#v", resp)
	}
	for idx, id := range resp.Ids {
		if idx > 5 {
			break
		}
		fmt.Printf("%s\n", id)
	}
	// Output: 36152127
	// 36152166
	// 36152181
	// 36596285
	// 36715809
	// 36763407
}

func ExampleGetMerchantDetails() {

	payload := GetMerchantDetailsRequest{
		ANetApiRequest: ANetApiRequest{
			MerchantAuthentication: MerchantAuthentication{
				Name:           apiName,
				TransactionKey: apiKey,
			},
		},
	}

	var resp GetMerchantDetailsResponse

	if err := client.Send(ctx, endpt, &payload, &resp); err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}

	if resp.Messages.ResultCode == "Error" {
		fmt.Printf("%#v", resp)
	}
	//fmt.Printf("%#v", resp)
	// Output: authorize_net.GetMerchantDetailsResponse{ANetApiResponse:authorize_net.ANetApiResponse{RefId:"", Messages:authorize_net.Messages{ResultCode:"Ok", Messages:[]authorize_net.ErrMessage{authorize_net.ErrMessage{Code:0x1, Text:"Successful.", Description:""}}}, SessionToken:""}, IsTestMode:false, Processors:[]authorize_net.Processor{authorize_net.Processor{Name:"First Data Nashville", Id:2, CardTypes:[]string{"A", "D", "M", "V"}}}, MerchantName:"Test Developer", GatewayId:"482527", MarketTypes:[]string{"eCommerce"}, ProductCodes:[]string{"CNP"}, PaymentMethods:[]string{"AmericanExpress", "Discover", "Echeck", "MasterCard", "Paypal", "Visa", "VisaCheckout", "ApplePay", "AndroidPay"}, Currencies:[]string{"USD"}, PublicClientKey:"4cEB6WrfHGS76gIW3v7btBCE3HuuAuke9Pj96Ztfn5R32G5ep42vne7MCWZtAVyd", BusinessInformation:(*authorize_net.CustomerAddress)(0xc0003e4000), MerchantTimeZone:"America/Los_Angeles", ContactDetails:[]authorize_net.ContactDetail{authorize_net.ContactDetail{Email:"bmcmanus@visa.com", FirstName:"Sandbox", LastName:"Default"}}}
	fmt.Println(convert.PrettyPrint(resp))
	// 	{
	//         "messages": {
	//                 "resultCode": "Ok",
	//                 "message": [
	//                         {
	//                                 "code": 1,
	//                                 "text": "Successful."
	//                         }
	//                 ]
	//         },
	//         "processors": [
	//                 {
	//                         "name": "First Data Nashville",
	//                         "id": 2,
	//                         "cardTypes": [
	//                                 "A",
	//                                 "D",
	//                                 "J",
	//                                 "M",
	//                                 "V"
	//                         ]
	//                 }
	//         ],
	//         "merchantName": "Foo Bar",
	//         "gatewayId": "482527",
	//         "marketTypes": [
	//                 "eCommerce"
	//         ],
	//         "productCodes": [
	//                 "CNP"
	//         ],
	//         "paymentMethods": [
	//                 "AmericanExpress",
	//                 "Discover",
	//                 "Echeck",
	//                 "JCB",
	//                 "MasterCard",
	//                 "Visa"
	//         ],
	//         "currencies": [
	//                 "USD"
	//         ],
	//         "publicClientKey": "4rJ7h7m6d3ZfFUU8AUVW5m9wZB85c8d6BbSx8rLUe9Lw9HzdU6zGAcHQ456GeK9J",
	//         "businessInformation": {
	//                 "company": "Foo Bar LLC",
	//                 "address": "1 Main Street",
	//                 "city": "Bellevue",
	//                 "state": "WA",
	//                 "zip": "98004",
	//                 "country": "US",
	//                 "phoneNumber": "425-555-1212"
	//         },
	//         "merchantTimeZone": "America/Los_Angeles",
	//         "contactDetails": [
	//                 {
	//                         "email": "foo@bar.com",
	//                         "firstName": "TestFirstName",
	//                         "lastName": "TestLastName"
	//                 }
	//         ]
	// }

}

func ExampleUpdateMerchantDetails() {

	payload := UpdateMerchantDetailsRequest{
		Payload: UpdateMerchantDetailsPayload{
			ANetApiRequest: ANetApiRequest{
				MerchantAuthentication: MerchantAuthentication{
					Name:           apiName,
					TransactionKey: apiKey,
				},
			},
			IsTestMode: true,
		},
	}
	req, err := client.prepareRequest(context.Background(), endpt, &payload)
	if err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}

	var resp UpdateMerchantDetailsResponse
	if err := client.Do(req, &resp, false); err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}
	if resp.Messages.ResultCode == "Error" {
		fmt.Printf("%#v", resp)
	}
	//fmt.Printf("%#v", resp)
	// Output: authorize_net.UpdateMerchantDetailsResponse{ANetApiResponse:authorize_net.ANetApiResponse{RefId:"", Messages:authorize_net.Messages{ResultCode:"Error", Messages:[]authorize_net.ErrMessage{authorize_net.ErrMessage{Code:0x47, Text:"The authentication type is not allowed for this method call.", Description:""}}}, SessionToken:""}}
	fmt.Println(convert.PrettyPrint(resp))
}

func ExampleGetCustomerShippingAddress() {

	payload := GetCustomerShippingAddressRequest{
		Payload: GetCustomerShippingAddressPayload{
			ANetApiRequest: ANetApiRequest{
				MerchantAuthentication: MerchantAuthentication{
					Name:           apiName,
					TransactionKey: apiKey,
				},
			},
			CustomerProfileId: "1920672921",
			//CustomerAddressId: "1877745863",// optional, if you created the record yourself
		},
	}
	req, err := client.prepareRequest(context.Background(), endpt, &payload)
	if err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}

	var resp GetCustomerShippingAddressResponse
	if err := client.Do(req, &resp, false); err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}
	if resp.Messages.ResultCode == "Error" {
		fmt.Printf("%#v", resp)
	}
	fmt.Printf("%#v", resp.Address)
	// Output: &authorize_net.CustomerAddressEx{CustomerAddress:authorize_net.CustomerAddress{NameAndAddress:authorize_net.NameAndAddress{FirstName:"Newfirstname", LastName:"Doe", Company:"", Address:"123 Main St.", City:"Bellevue", State:"WA", Zip:"98004", Country:"USA"}, PhoneNumber:"000-000-0000", FaxNumber:"", Email:""}, CustomerAddressId:"1810861269"}
}

func ExampleCreateCustomerShippingAddress() {

	payload := CreateCustomerShippingAddressRequest{
		Payload: CreateCustomerShippingAddressPayload{
			ANetApiRequest: ANetApiRequest{
				MerchantAuthentication: MerchantAuthentication{
					Name:           apiName,
					TransactionKey: apiKey,
				},
			},
			CustomerProfileId: "1920672921",
			Address: &CustomerAddress{
				NameAndAddress: NameAndAddress{
					FirstName: "FirstName",
					LastName:  "LastName",
					Company:   "SomeCompany LTD",
					Address:   "Hope Street, ground floor",
					City:      "Oz",
					State:     "Oz",
					Zip:       "300200",
					Country:   "Norway",
				},
				PhoneNumber: "1131131131",
				FaxNumber:   "what?",
				Email:       "36152127@authorize.net@authorize.net",
			},
		},
	}
	req, err := client.prepareRequest(context.Background(), endpt, &payload)
	if err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}

	var resp CreateCustomerShippingAddressResponse
	if err := client.Do(req, &resp, false); err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}
	if resp.Messages.ResultCode == "Error" {
		fmt.Printf("%#v", resp)
	}
	fmt.Printf("%#v", resp)
	// Output: authorize_net.CreateCustomerShippingAddressResponse{ANetApiResponse:authorize_net.ANetApiResponse{RefId:"", Messages:authorize_net.Messages{ResultCode:"Ok", Messages:[]authorize_net.ErrMessage{authorize_net.ErrMessage{Code:0x1, Text:"Successful.", Description:""}}}, SessionToken:""}, CustomerProfileId:"1920672921", CustomerAddressId:"1877745883"}
}

func ExampleUpdateCustomerShippingAddress() {

	payload := UpdateCustomerShippingAddressRequest{
		Payload: UpdateCustomerShippingAddressPayload{
			ANetApiRequest: ANetApiRequest{
				MerchantAuthentication: MerchantAuthentication{
					Name:           apiName,
					TransactionKey: apiKey,
				},
			},
			CustomerProfileId: "1920672921",
			Address: &CustomerAddressEx{
				CustomerAddressId: "1877745883",
				CustomerAddress: CustomerAddress{
					NameAndAddress: NameAndAddress{
						FirstName: "FirstName",
						LastName:  "LastName",
						Company:   "SomeCompany LTD",
						Address:   "Hope Street, ground floor",
						City:      "Oz",
						State:     "Oz",
						Zip:       "300200",
						Country:   "China",
					},
					PhoneNumber: "1131131131",
					FaxNumber:   "what?",
					Email:       "36152127@novalidation@authorize.net",
				},
			},
		},
	}
	req, err := client.prepareRequest(context.Background(), endpt, &payload)
	if err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}

	var resp UpdateCustomerShippingAddressResponse
	if err := client.Do(req, &resp, false); err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}
	if resp.Messages.ResultCode == "Error" {
		fmt.Printf("%#v", resp)
	}
	fmt.Printf("%#v", resp)
	// Output: authorize_net.UpdateCustomerShippingAddressResponse{ANetApiResponse:authorize_net.ANetApiResponse{RefId:"", Messages:authorize_net.Messages{ResultCode:"Ok", Messages:[]authorize_net.ErrMessage{authorize_net.ErrMessage{Code:0x1, Text:"Successful.", Description:""}}}, SessionToken:""}}

}

func ExampleDeleteCustomerShippingAddress() {

	payload := DeleteCustomerShippingAddressRequest{
		Payload: DeleteCustomerShippingAddressPayload{
			ANetApiRequest: ANetApiRequest{
				MerchantAuthentication: MerchantAuthentication{
					Name:           apiName,
					TransactionKey: apiKey,
				},
			},
			CustomerProfileId: "1920672921",
			CustomerAddressId: "1877745863",
		},
	}
	req, err := client.prepareRequest(context.Background(), endpt, &payload)
	if err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}

	var resp DeleteCustomerShippingAddressResponse
	if err := client.Do(req, &resp, false); err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}
	if resp.Messages.ResultCode == "Error" {
		fmt.Printf("%#v", resp)
	}
	fmt.Printf("%#v", resp)
	// Output: authorize_net.DeleteCustomerShippingAddressResponse{ANetApiResponse:authorize_net.ANetApiResponse{RefId:"", Messages:authorize_net.Messages{ResultCode:"Ok", Messages:[]authorize_net.ErrMessage{authorize_net.ErrMessage{Code:0x1, Text:"Successful.", Description:""}}}, SessionToken:""}}

}

func ExampleCreateCustomerPaymentProfile() {

	payload := CreateCustomerPaymentProfileRequest{
		Payload: CreateCustomerPaymentProfilePayload{
			ANetApiRequest: ANetApiRequest{
				MerchantAuthentication: MerchantAuthentication{
					Name:           apiName,
					TransactionKey: apiKey,
				},
			},
			CustomerProfileId: "1920672070",
			PaymentProfile: &CustomerPaymentProfile{
				CustomerPaymentProfileBase: CustomerPaymentProfileBase{
					CustomerType: 0,
					BillTo: &CustomerAddress{
						NameAndAddress: NameAndAddress{
							FirstName: "Johnny",
							LastName:  "Doe",
							Country:   "USA",
						},
						Email: "gmail@chuck.norris",
					},
				},
				Payment: &Payment{
					CreditCard: &CreditCard{
						CardNumber:     cardExample,
						ExpirationDate: cardExpirationDate,
						CardCode:       cardCode,
					},
				},
				DefaultPaymentProfile: true,
			},
			ValidationMode: ValidationModeTestMode,
		},
	}
	req, err := client.prepareRequest(context.Background(), endpt, &payload)
	if err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}

	var resp CreateCustomerPaymentProfileResponse
	if err := client.Do(req, &resp, false); err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}
	if resp.Messages.ResultCode == "Error" {
		fmt.Printf("%#v", resp)
	}
	fmt.Printf("%#v", resp)
	// Output: authorize_net.CreateCustomerPaymentProfileResponse{ANetApiResponse:authorize_net.ANetApiResponse{RefId:"", Messages:authorize_net.Messages{ResultCode:"Ok", Messages:[]authorize_net.ErrMessage{authorize_net.ErrMessage{Code:0x1, Text:"Successful.", Description:""}}}, SessionToken:""}, CustomerProfileId:"1920672070", CustomerPaymentProfileId:"1833657059", ValidationDirectResponse:"1,1,1,(TESTMODE) This transaction has been approved.,000000,P,0,none,Test transaction for ValidateCustomerPaymentProfile.,1.00,CC,auth_only,none,Johnny,Doe,,,,,,USA,,,email@example.com,,,,,,,,,0.00,0.00,0.00,FALSE,none,,,,,,,,,,,,,,XXXX0015,MasterCard,,,,,,,,,,,,,,,,,"}

}

func ExampleValidateCustomerPaymentProfile() {

	payload := ValidateCustomerPaymentProfileRequest{
		Payload: ValidateCustomerPaymentProfilePayload{
			ANetApiRequest: ANetApiRequest{
				MerchantAuthentication: MerchantAuthentication{
					Name:           apiName,
					TransactionKey: apiKey,
				},
			},
			CustomerProfileId:        "1920672070",
			CustomerPaymentProfileId: "1833657059",
			ValidationMode:           ValidationModeTestMode,
		},
	}
	req, err := client.prepareRequest(context.Background(), endpt, &payload)
	if err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}

	var resp ValidateCustomerPaymentProfileResponse
	if err := client.Do(req, &resp, false); err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}
	if resp.Messages.ResultCode == "Error" {
		fmt.Printf("%#v", resp)
	}
	fmt.Printf("%#v", resp)
	// Output: authorize_net.ValidateCustomerPaymentProfileResponse{ANetApiResponse:authorize_net.ANetApiResponse{RefId:"", Messages:authorize_net.Messages{ResultCode:"Ok", Messages:[]authorize_net.ErrMessage{authorize_net.ErrMessage{Code:0x1, Text:"Successful.", Description:""}}}, SessionToken:""}, DirectResponse:"1,1,1,(TESTMODE) This transaction has been approved.,000000,P,0,none,Test transaction for ValidateCustomerPaymentProfile.,1.00,CC,auth_only,jdoe4801,Johnny,Doe,,,,,,USA,,,4953@mail.com,,,,,,,,,0.00,0.00,0.00,FALSE,none,,,,,,,,,,,,,,XXXX0015,MasterCard,,,,,,,,,,,,,,,,,"}

}

func ExampleUpdateCustomerPaymentProfile() {

	payload := UpdateCustomerPaymentProfileRequest{
		Payload: UpdateCustomerPaymentProfilePayload{
			ANetApiRequest: ANetApiRequest{
				MerchantAuthentication: MerchantAuthentication{
					Name:           apiName,
					TransactionKey: apiKey,
				},
			},
			CustomerProfileId: "1920672070",
			ValidationMode:    ValidationModeNone,
			PaymentProfile: &CustomerPaymentProfileEx{
				CustomerPaymentProfileId: "1833657059",
				CustomerPaymentProfile: CustomerPaymentProfile{
					CustomerPaymentProfileBase: CustomerPaymentProfileBase{
						CustomerType: Business,
					},
					Payment: &Payment{

						BankAccount: &BankAccount{
							AccountType:   BankAccountTypeSavings,
							RoutingNumber: "133563585",
							AccountNumber: "0123456789",
							NameOnAccount: "SALSA",
							EcheckType:    EcheckCCD,
							BankName:      "SALSA",
							CheckNumber:   "123",
						},
					},
					DriversLicense: &DriversLicense{
						Number:      "12388381813",
						State:       "CA",
						DateOfBirth: "2000-10-10",
					},
					DefaultPaymentProfile: true,
				},
			},
		},
	}
	req, err := client.prepareRequest(context.Background(), endpt, &payload)
	if err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}

	var resp UpdateCustomerPaymentProfileResponse
	if err := client.Do(req, &resp, false); err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}
	if resp.Messages.ResultCode == "Error" {
		fmt.Printf("%#v", resp)
	}
	fmt.Printf("%#v", resp)
	// Output: authorize_net.UpdateCustomerPaymentProfileResponse{ANetApiResponse:authorize_net.ANetApiResponse{RefId:"", Messages:authorize_net.Messages{ResultCode:"Ok", Messages:[]authorize_net.ErrMessage{authorize_net.ErrMessage{Code:0x1, Text:"Successful.", Description:""}}}, SessionToken:""}, ValidationDirectResponse:""}

}

func ExampleDeleteCustomerPaymentProfile() {

	payload := DeleteCustomerPaymentProfileRequest{
		Payload: DeleteCustomerPaymentProfilePayload{
			ANetApiRequest: ANetApiRequest{
				MerchantAuthentication: MerchantAuthentication{
					Name:           apiName,
					TransactionKey: apiKey,
				},
			},
			CustomerProfileId:        "1920672070",
			CustomerPaymentProfileId: "1833657059",
		},
	}
	req, err := client.prepareRequest(context.Background(), endpt, &payload)
	if err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}

	var resp DeleteCustomerPaymentProfileResponse
	if err := client.Do(req, &resp, false); err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}
	if resp.Messages.ResultCode == "Error" {
		fmt.Printf("%#v", resp)
	}
	fmt.Printf("%#v", resp)
	// Output: authorize_net.DeleteCustomerPaymentProfileResponse{ANetApiResponse:authorize_net.ANetApiResponse{RefId:"", Messages:authorize_net.Messages{ResultCode:"Ok", Messages:[]authorize_net.ErrMessage{authorize_net.ErrMessage{Code:0x1, Text:"Successful.", Description:""}}}, SessionToken:""}}

}

// NOT the AcceptPayment scheme; this regards 3rd-party tokenization-service provider
func ExampleDecryptPaymentData() {

	payload := DecryptPaymentDataRequest{
		Payload: DecryptPaymentDataPayload{
			ANetApiRequest: ANetApiRequest{
				MerchantAuthentication: MerchantAuthentication{
					Name:           apiName,
					TransactionKey: apiKey,
				},
			},
			OpaqueData: &OpaqueData{
				DataDescriptor: "COMMON.VCO.ONLINE.PAYMENT", // always send this
				DataValue:      "+s6+Xn9fCF8kfA+5TAAMpX3dejnEg2u5HUZMIVj43M7JpjuIeCahMVXZoHJIQyVaAaWqsWlwgZSwYTh+gSOx8PT1G70wplfnDMU1Wcv1WDYnqz+fI1NMfShqUorC2GDRt2uLksAHH7Zst7GPmsoKa2lY3Do86s+fxwC7ZK86NRVW7Y7JrvEmEI7fxTUai6/vn8kHGq3OOEMSrqnzJERgDx6wdxSLZeFBGubVwdJUPqXUisukGS+QG3ipqL8+kSHmBsbad375scFGF1AeZU7H+8kO7Wzr1QkPnx3pQsiMQGxghmkGX1wmZ0ilmr344ytfKx0Gr3v8JNCw6B6LhXscu18KUXsidLW33pxtjqLQRAaFHgNh0QNMJZkbLEHmtywY16m4NWnDFzlFIk+y+iaonbfZlrEfdkmZhePlXW3N6UhwFnozI5vMsl7E1cqlVO6TJO1ocfEqlnKdFBCTqdeuzaXuUbSi7IUPEtEWFOb8WoKDU+0Ae5LjXVH1jNVN8XC4S9HGIibV69xHKaE155DU3rZjrMFfBcQIufdbbQI3qBXVK1e3J9B1FLMtAbYxn2ZtGCyWxjQ13wq5OECfjR0u1xrrTd0VPzwBQhQ+aqDCTPkQRoYdKU+p9GaKy5NKiqxLWqeTu3bGRQpffps2jEZDIUrdJXc2t5VNA8F+KtN933PCUuVZROC7ADSoSY3mFN5PQEknkg5GQpXqJZNFBRITUyleTBIFBV0sELUOFy81DBfkhSnsjSb+X8TNez0qaG2NkMVhF7oOIBYawT7NSUvxwbquZw7GEbSss4yNl94zDKOK1CR2cAZsYtTGlWGIhiwCsFCiKVuGUkioF3M5gXAkPhkw20V69ed1DE7DJu8PF89U2FQI7p6UPE7XvHVbSrXHl0m7gP8w+QPYvJBNtJATgjqDD+xCQDa0csh+1p17UtaTR8I3XDPQWNl7NV2msWzhyZxMi1wLR+nlnYEn7K97uihRqcoAjnIFpXAsoY+tCXiig7xkcbe4W4rvl610Zr7QuVQ0xovXUX48+PLYb1uGzR+RQ0hR914syxpe72HPS/VrZXwmEmtnlviedx1cjVThwKSQ59vv5+zVNgcQJt51snuT+3xpSzfPTjX5UiTZFYZ09tPbWiCJ88y3UujFI8zhn4EjNQNKzDMXvh84EImty46uMaB0Ehjfg50s2FXSMO18VaCA7VTUuj0dvQkL8Zg2aNCWlJEjiNI9AUnyBZxJgM9elX4RJUhuurNSAq+OizMx9E+pNwUvF+L270TNUBQ4RAGpmP2QSB4dw8rJ0yL7V10TD2gq4J1lHhpOOD94IVH0XrwmusRieFt44rakA7rw4zipEkprqP6UO6q6OR9cgm4wphBS2lBIYyvexsMh+Y6J0sH6Ixu8QEsRSWKlv+aLULu0c42K86czgFDYkJnNlbyHbTXFXzAHuWlAvSoMyj9m4eK5vZp9JjzAa2duyofOlpMAvjbnaV8UeAGIPI+QK3imm8D6+VAXKBqTVpnpqQsRIDJb8Pxu7pmDBUyzwO9NCdhmNLomhROpnNOqp65p4neqhyp4kHsLq4vTLBvx6pfAsyrwqdy4QyHRCIDnX6wcy6J4MOW+gfM6Hm1cNm6AqcOaiTafXRR3TdOShOzXBUm3gr6o4dYi3l+oQx9LwqFgoD8Bg+3u0PWEVZODkIrcLETyKWjao4s5YRratWclo0In3mfYZiO3kSxoDRAQoi8BWVwiWEnstNZhx0edLcIseWNCQ8GHDXPWQzs+NMHWRQT8yFwPC34iMOrRPk0TlM1CWySBC5LBaD8ZNZ3R+al2XwGP6wwXbJtAvwR1a3Wqqb+vGFdggp0ISPHQTI7I1w4Kp0ijXA36rTvmZ9xin0sN+ayOtoNfvBo4blj9FHKgDoqWimfzxsrwOcFwWU1i0Xfd4wvriv73Z6gqXhjmS19S7zuVk2+TtPlUPKiHj7fpPON18euYuVe8jszM3xcMNiSDmG+HWjuUI78kyYQEEsxMmGQHt+nenScbjBu9JkZP7KrH380G34zMJeJF7A7GzzrErrCsEsMHPFuJzJkogwnfL5zHy24UdYOBDc8+4aTaUmAxu9shdvoWCj3i1Q==",
				DataKey:        "qGl+69iAWXJ13+18jHgBO2zHCuekawWcApfLhGnbYXD4GsI9EOM2Y5V8zGXvr1lF3hjGMT0NoD2thxzR7CrAvgfgw7lAJzlGIACOZnEkx70ywnJHAR3sGO8hyq9f1Fk",
			},
			CallId: "2186595692635007317",
		},
	}
	req, err := client.prepareRequest(context.Background(), endpt, &payload)
	if err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}

	var resp DecryptPaymentDataResponse
	if err := client.Do(req, &resp, false); err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}
	if resp.Messages.ResultCode == "Error" {
		fmt.Printf("%#v", resp)
	}
	// Note : for some reason it always return an error
	fmt.Printf("%#v", resp)
	// Output: authorize_net.DecryptPaymentDataResponse{ANetApiResponse:authorize_net.ANetApiResponse{RefId:"", Messages:authorize_net.Messages{ResultCode:"Error", Messages:[]authorize_net.ErrMessage{authorize_net.ErrMessage{Code:0xd, Text:"An error occurred during processing. Please try again.", Description:""}}}, SessionToken:""}, ShippingInfo:(*authorize_net.CustomerAddress)(nil), BillingInfo:(*authorize_net.CustomerAddress)(nil), CardInfo:(*authorize_net.CreditCardMasked)(nil), PaymentDetails:(*authorize_net.PaymentDetails)(nil)}

}

func ExampleGetTransactionListForCustomer() {

	payload := GetTransactionListForCustomerRequest{
		Payload: GetTransactionListForCustomerPayload{
			ANetApiRequest: ANetApiRequest{
				MerchantAuthentication: MerchantAuthentication{
					Name:           apiName,
					TransactionKey: apiKey,
				},
			},
			Sorting: Sorting{
				OrderBy:         "submitTimeUTC",
				OrderDescending: true,
			},
			Paging: Paging{
				Offset: 1,
				Limit:  30,
			},
			CustomerProfileId: "1920672070",
		},
	}
	req, err := client.prepareRequest(context.Background(), endpt, &payload)
	if err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}

	var resp GetTransactionListResponse
	if err := client.Do(req, &resp, false); err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}
	if resp.Messages.ResultCode == "Error" {
		fmt.Printf("%#v", resp)
	}
	for _, trans := range resp.Transactions {
		fmt.Printf("%#v\n", trans.TransId)
	}
	// Output: "60126763648"
	//"60126763647"
	//"60126763644"
}

func ExampleSendCustomerTransactionReceipt() {

	payload := SendCustomerTransactionReceiptRequest{
		Payload: SendCustomerTransactionReceiptPayload{
			ANetApiRequest: ANetApiRequest{
				MerchantAuthentication: MerchantAuthentication{
					Name:           apiName,
					TransactionKey: apiKey,
				},
			},
			TransId:       "60126763648",
			CustomerEmail: "put_valid_email_here",
		},
	}
	req, err := client.prepareRequest(context.Background(), endpt, &payload)
	if err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}

	var resp SendCustomerTransactionReceiptResponse
	if err := client.Do(req, &resp, false); err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}
	if resp.Messages.ResultCode == "Error" {
		fmt.Printf("%#v", resp)
	}
	fmt.Printf("%v", resp)
	// Output:
}

func ExampleGetTransactionDetails() {

	payload := GetTransactionDetailsRequest{
		Payload: GetTransactionDetailsPayload{
			ANetApiRequest: ANetApiRequest{
				MerchantAuthentication: MerchantAuthentication{
					Name:           apiName,
					TransactionKey: apiKey,
				},
			},
			TransId: "60126763648",
		},
	}
	req, err := client.prepareRequest(context.Background(), endpt, &payload)
	if err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}

	var resp GetTransactionDetailsResponse
	if err := client.Do(req, &resp, false); err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}
	if resp.Messages.ResultCode == "Error" {
		fmt.Printf("%#v", resp)
	}
	fmt.Printf("%#v", resp.Transaction.Batch.BatchId)
	// Output: "9697004"
}

func ExampleGetTransactionList() {

	payload := GetTransactionListRequest{
		Payload: GetTransactionListPayload{
			ANetApiRequest: ANetApiRequest{
				MerchantAuthentication: MerchantAuthentication{
					Name:           apiName,
					TransactionKey: apiKey,
				},
			},
			Paging: Paging{
				Limit:  50,
				Offset: 1,
			},
			Sorting: Sorting{
				OrderBy:         "submitTimeUTC",
				OrderDescending: false,
			},
			BatchId: "9697004",
		},
	}
	req, err := client.prepareRequest(context.Background(), endpt, &payload)
	if err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}

	var resp GetTransactionListResponse
	if err := client.Do(req, &resp, false); err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}
	if resp.Messages.ResultCode == "Error" {
		fmt.Printf("%#v", resp)
	}
	fmt.Printf("%#v", resp.TotalNumInResultSet)
	// Output: 50

}

func ExampleGetUnsettledTransactionList() {

	payload := GetUnsettledTransactionListRequest{
		Payload: GetUnsettledTransactionListPayload{
			ANetApiRequest: ANetApiRequest{
				MerchantAuthentication: MerchantAuthentication{
					Name:           apiName,
					TransactionKey: apiKey,
				},
			},
			Paging: Paging{
				Limit:  50,
				Offset: 1,
			},
			Sorting: Sorting{
				OrderBy:         "submitTimeUTC",
				OrderDescending: false,
			},
		},
	}
	req, err := client.prepareRequest(context.Background(), endpt, &payload)
	if err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}

	var resp GetUnsettledTransactionListResponse
	if err := client.Do(req, &resp, false); err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}
	if resp.Messages.ResultCode == "Error" {
		fmt.Printf("%#v", resp)
	}
	fmt.Printf("%#v", resp.TotalNumInResultSet)
	// Output: 50
}

func ExampleARBCreateSubscription() {

	payload := ARBCreateSubscriptionRequest{
		Payload: ARBCreateSubscriptionPayload{
			ANetApiRequest: ANetApiRequest{
				MerchantAuthentication: MerchantAuthentication{
					Name:           apiName,
					TransactionKey: apiKey,
				},
			},
			Subscription: ARBSubscription{
				Name: "Sample",
				PaymentSchedule: &PaymentSchedule{
					Interval: &PaymentScheduleTypeInterval{
						Length: 1,
						Unit:   Months,
					},
					StartDate:        NowTime(),
					TotalOccurrences: 12,
					TrialOccurrences: 1,
				},
				Amount:      10.29,
				TrialAmount: 0.00,
				Payment: &Payment{
					CreditCard: &CreditCard{
						CardNumber:     cardExample,
						ExpirationDate: cardExpirationDate,
						CardCode:       cardCode,
						IsPaymentToken: false,
					},
				},
				BillTo: &NameAndAddress{
					FirstName: "Papas",
					LastName:  "RollingStone",
				},
			},
		},
	}
	req, err := client.prepareRequest(context.Background(), endpt, &payload)
	if err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}

	var resp ARBCreateSubscriptionResponse
	if err := client.Do(req, &resp, false); err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}
	if resp.Messages.ResultCode == "Error" {
		fmt.Printf("%#v", resp)
	}
	fmt.Printf("%#v", resp.Messages.ResultCode)
	// Output: "Ok"
}

func ExampleARBCancelSubscription() {

	payload := ARBCancelSubscriptionRequest{
		Payload: ARBCancelSubscriptionPayload{
			ANetApiRequest: ANetApiRequest{
				MerchantAuthentication: MerchantAuthentication{
					Name:           apiName,
					TransactionKey: apiKey,
				},
			},
			SubscriptionId: "5982761",
		},
	}
	req, err := client.prepareRequest(context.Background(), endpt, &payload)
	if err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}

	var resp ARBCancelSubscriptionResponse
	if err := client.Do(req, &resp, false); err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}
	if resp.Messages.ResultCode == "Error" {
		fmt.Printf("%#v", resp)
	}
	fmt.Printf("%#v", resp.Messages.ResultCode)
	// Output: "Ok"
}

func ExampleARBGetSubscriptionList() {

	payload := ARBGetSubscriptionListRequest{
		Payload: ARBGetSubscriptionListPayload{
			ANetApiRequest: ANetApiRequest{
				MerchantAuthentication: MerchantAuthentication{
					Name:           apiName,
					TransactionKey: apiKey,
				},
			},
			SearchType: "subscriptionActive",
			Sorting: Sorting{
				OrderBy:         "id",
				OrderDescending: false,
			},
			Paging: Paging{
				Limit:  1,
				Offset: 100,
			},
		},
	}
	req, err := client.prepareRequest(context.Background(), endpt, &payload)
	if err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}

	var resp ARBGetSubscriptionListResponse
	if err := client.Do(req, &resp, false); err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}
	if resp.Messages.ResultCode == "Error" {
		fmt.Printf("%#v", resp)
	}
	for _, subscription := range resp.SubscriptionDetails {
		fmt.Printf("%#v", subscription.Id)
	}
	// Output: "5517694"

}

func ExampleARBGetSubscription() {

	payload := ARBGetSubscriptionRequest{
		Payload: ARBGetSubscriptionPayload{
			ANetApiRequest: ANetApiRequest{
				MerchantAuthentication: MerchantAuthentication{
					Name:           apiName,
					TransactionKey: apiKey,
				},
			},
			SubscriptionId: "5517694",
		},
	}
	req, err := client.prepareRequest(context.Background(), endpt, &payload)
	if err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}

	var resp ARBGetSubscriptionResponse
	if err := client.Do(req, &resp, false); err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}
	if resp.Messages.ResultCode == "Error" {
		fmt.Printf("%#v", resp)
	}
	fmt.Printf("%#v", resp.Subscription.Name)
	// Output: "Sample Subscription"

}

func ExampleARBGetSubscriptionStatus() {

	payload := ARBGetSubscriptionStatusRequest{
		Payload: ARBGetSubscriptionStatusPayload{
			ANetApiRequest: ANetApiRequest{
				MerchantAuthentication: MerchantAuthentication{
					Name:           apiName,
					TransactionKey: apiKey,
				},
			},
			SubscriptionId: "5517694",
		},
	}
	req, err := client.prepareRequest(context.Background(), endpt, &payload)
	if err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}

	var resp ARBGetSubscriptionStatusResponse
	if err := client.Do(req, &resp, false); err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}
	if resp.Messages.ResultCode == "Error" {
		fmt.Printf("%#v", resp)
	}
	fmt.Printf("%#v", resp.Status)
	// Output: "active"

}

func ExampleARBUpdateSubscription() {

	payload := ARBUpdateSubscriptionRequest{
		Payload: ARBUpdateSubscriptionPayload{
			ANetApiRequest: ANetApiRequest{
				MerchantAuthentication: MerchantAuthentication{
					Name:           apiName,
					TransactionKey: apiKey,
				},
			},
			SubscriptionId: "5517694",
			Subscription: ARBSubscription{
				Name: "Updated Subscription",
				Payment: &Payment{
					CreditCard: &CreditCard{
						CardNumber:     "4111111111111111",
						ExpirationDate: "2022-12",
					},
				},
			},
		},
	}
	req, err := client.prepareRequest(context.Background(), endpt, &payload)
	if err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}

	var resp ARBUpdateSubscriptionResponse
	if err := client.Do(req, &resp, false); err != nil {
		client.log.Printf("Error : %v\n", err)
		fmt.Printf("Error : %v", err)
		return
	}
	if resp.Messages.ResultCode == "Error" {
		fmt.Printf("%#v", resp)
	}
	fmt.Printf("%#v", resp.Messages.ResultCode)
	// Output: "Ok"
}

func ExampleANetApi() {

	fmt.Sprintf("%v", client)

}

// *********  NOTHING BELOW here is IMPLEMENTED  **************

func ExampleGetAUJobDetails() {

	fmt.Sprintf("%v", client)

}

func ExampleGetAUJobSummary() {

	fmt.Sprintf("%v", client)

}

func ExampleGetBatchStatistics() {

	fmt.Sprintf("%v", client)

}

func ExampleGetHostedPaymentPage() {

	fmt.Sprintf("%v", client)

}

func ExampleGetHostedProfilePage() {

	fmt.Sprintf("%v", client)

}

func ExampleGetSettledBatchList() {

	fmt.Sprintf("%v", client)

}

func ExampleIsAlive() {

	fmt.Sprintf("%v", client)

}

func ExampleLogout() {

	fmt.Sprintf("%v", client)

}

func ExampleMobileDeviceLogin() {

	fmt.Sprintf("%v", client)

}

func ExampleMobileDeviceRegistration() {

	fmt.Sprintf("%v", client)

}

func ExampleSecurePaymentContainer() {

	fmt.Sprintf("%v", client)

}

func ExampleUpdateHeldTransaction() {

	fmt.Sprintf("%v", client)

}

func ExampleUpdateSplitTenderGroup() {

	fmt.Sprintf("%v", client)

}
