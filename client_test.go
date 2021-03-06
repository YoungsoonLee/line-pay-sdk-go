package linepay

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
)

func TestClient_PaymentsRequest(t *testing.T) {

	client, err := NewClient(ChannelID, ChannelSecret, &Signer{ChannelId: ChannelID}, &ClientOpts{})
	if err != nil {
		t.Errorf("New() error = %v", err.Error())
		return
	}

	data := PaymentsRequest{
		Amount:   100,
		Currency: "TWD",
		OrderID:  "test_order_15",
		Packages: []PaymentsPackageRequest{
			PaymentsPackageRequest{
				ID:     "pkg_id_1",
				Amount: 100,
				Name:   "pkg_name_1",
				Products: []PaymentsPackageProductRequest{
					PaymentsPackageProductRequest{
						Name:     "prod_1",
						Quantity: 1,
						Price:    100,
					},
				},
			},
		},
		RedirectUrls: PaymentsRedirectUrlsRequest{
			ConfirmURLType: PaymentsConfirmUrlTypeClient,
			ConfirmURL:     CallbackHost + "/confirm",
			CancelURL:      CallbackHost + "/cancel",
		},
	}

	res, err := client.PaymentsRequest(context.Background(), &data)
	if err != nil {
		t.Errorf("Test PaymentsRequest failed: %s", err.Error())
	}

	t.Logf("Dump PaymentsRequest body: '%+v'", res)

}

func TestClient_PaymentsRequestAndConfirm(t *testing.T) {

	t.Parallel()

	client, err := NewClient(ChannelID, ChannelSecret, &Signer{ChannelId: ChannelID}, &ClientOpts{})
	// client, err := NewClient(ChannelID, ChannelSecret, APIHostSandbox, "http://localhost:9988")
	if err != nil {
		t.Errorf("New() error = %v", err.Error())
		return
	}

	data := PaymentsRequest{
		Amount:   100,
		Currency: "TWD",
		OrderID:  "test_order_16",
		Packages: []PaymentsPackageRequest{
			PaymentsPackageRequest{
				ID:     "pkg_id_1",
				Amount: 100,
				Name:   "pkg_name_1",
				Products: []PaymentsPackageProductRequest{
					PaymentsPackageProductRequest{
						Name:     "prod_1",
						Quantity: 1,
						Price:    100,
					},
				},
			},
		},
		RedirectUrls: PaymentsRedirectUrlsRequest{
			// ConfirmUrl: path.Join(CallbackHost, "/confirm"),
			// CancelUrl:  path.Join(CallbackHost, "/cancel"),
			ConfirmURLType: PaymentsConfirmUrlTypeClient,
			ConfirmURL:     CallbackHost + "/confirm",
			CancelURL:      CallbackHost + "/cancel",
		},
	}

	res, err := client.PaymentsRequest(context.Background(), &data)
	if err != nil {
		t.Errorf("Test PaymentsRequest failed: %s", err.Error())
	}

	fmt.Println("================")
	fmt.Printf("Open link in Web: '%s'\n", res.Info.PaymentURL.Web)
	fmt.Printf("Open link in App: '%s'\n", res.Info.PaymentURL.App)
	fmt.Println("================")

	t.Logf("Dump PaymentsRequest body: '%+v'", res)

	if res.ReturnCode != ApiReturnCodeSuccess {
		t.Errorf("return code is '%s' not '%s'", res.ReturnCode, ApiReturnCodeSuccess)
	}

}

func TestClient_PaymentsConfirm(t *testing.T) {

	t.Parallel()

	client, err := NewClient(ChannelID, ChannelSecret, &Signer{ChannelId: ChannelID}, &ClientOpts{})
	if err != nil {
		t.Errorf("New() error = %v", err.Error())
		return
	}

	data2 := PaymentsConfirmRequest{
		Amount:   100,
		Currency: "TWD",
	}

	res2, err := client.PaymentsConfirm(context.Background(), 2020010800227854310, &data2)
	if err != nil {
		t.Errorf("Test PaymentsRequest failed: %s", err.Error())
	}

	t.Logf("Dump PaymentsConfirm response body: '%+v'", res2)

}

func TestClient_PaymentsDetails(t *testing.T) {

	client, err := NewClient(ChannelID, ChannelSecret, &Signer{ChannelId: ChannelID}, &ClientOpts{})
	if err != nil {
		t.Errorf("New() error = %v", err.Error())
		return
	}

	data := PaymentsDetailsRequest{
		TransactionIDs: []int64{2020011300254002010, 2020010900231782310, 2020010900229878210},
		OrderIDs:       []string{"order_0be9807d-88cf-42fe-bf69-75a51f1ad83f", "order_9583d466-6c47-488b-813f-894c0a26d7e8", "order_d776f2dd-eb7a-4611-b8cc-53242b9d7e71"},
		Fields:         PaymentsDetailsFieldsDefault,
	}

	res2, err := client.PaymentsDetails(context.Background(), &data)
	if err != nil {
		t.Errorf("Test PaymentsDetails failed: %s", err.Error())
	}

	t.Logf("Dump PaymentsDetails response body: '%+v'", res2)

}

func TestClient_PaymentsCapture_1_Request(t *testing.T) {

	client, err := getClient()
	if err != nil {
		t.Errorf("New() error = %v", err.Error())
		return
	}

	data := PaymentsRequest{
		Amount:   100,
		Currency: "TWD",
		OrderID:  "test_order_29",
		Packages: []PaymentsPackageRequest{
			PaymentsPackageRequest{
				ID:     "pkg_id_1",
				Amount: 100,
				Name:   "pkg_name_1",
				Products: []PaymentsPackageProductRequest{
					PaymentsPackageProductRequest{
						Name:     "prod_1",
						Quantity: 1,
						Price:    100,
					},
				},
			},
		},
		RedirectUrls: PaymentsRedirectUrlsRequest{
			ConfirmURLType: PaymentsConfirmUrlTypeClient,
			ConfirmURL:     CallbackHost + "/confirm",
			CancelURL:      CallbackHost + "/cancel",
		},
		Options: PaymentsOptionsRequest{
			Payment: PaymentsOptionsPaymentRequest{
				Capture: false, // flag as false to go Capture API flow
			},
		},
	}

	a, _ := json.Marshal(data)
	fmt.Printf("\ndump PaymentRequest body: %s\n", string(a))
	fmt.Printf("\ndump PaymentRequest capture: %+v\n", data.Options.Payment.Capture)

	res, err := client.PaymentsRequest(context.Background(), &data)
	if err != nil {
		t.Errorf("Test PaymentsRequest failed: %s", err.Error())
	}

	printRequestInfo(res, true)
}

func TestClient_PaymentsCapture_2_Capture(t *testing.T) {

	client, err := NewClient(ChannelID, ChannelSecret, &Signer{ChannelId: ChannelID}, &ClientOpts{})
	if err != nil {
		t.Errorf("New() error = %v", err.Error())
		return
	}

	data2 := PaymentsCaptureRequest{
		Amount:   100,
		Currency: "TWD",
	}

	res2, err := client.PaymentsCapture(context.Background(), 2020011500264285210, &data2)
	if err != nil {
		t.Errorf("Test PaymentsRequest failed: %s", err.Error())
	}

	t.Logf("Dump PaymentsCapture response body: '%+v'", res2)

}

func getClient() (client *Client, err error) {
	return NewClient(ChannelID, ChannelSecret, &Signer{ChannelId: ChannelID}, &ClientOpts{})
}

func printRequestInfo(res *PaymentsResponse, dumpBody bool) {

	fmt.Println("================")

	if dumpBody {

		fmt.Printf("Dump body: '%+v'\n", res)
	}

	fmt.Printf("Open link in Web: '%s'\n", res.Info.PaymentURL.Web)
	fmt.Printf("Open link in App: '%s'\n", res.Info.PaymentURL.App)

	fmt.Println("================")
}
