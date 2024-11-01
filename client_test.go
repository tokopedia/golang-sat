package sat

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/jsonapi"
	"github.com/tokopedia/digital-b2b-client-library/golang-sat/signature"
)

const (
	PrivateKeyDummy = "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDiciTyLTBJJXZQ\nAV1CvlmcJY5Jgn+PmxWRu2/EBbJszE4LPiqoDG/xjYgIC1vejgyeePpKPcMHP8rC\nvyG4RLXiVnBwKIXwEqPmJ0V+VmO0XweBN0ipC+DquK5SjqNObLCNvCHLQfvFAWlv\n9otlqrHG/LTQzrU0lFvBZhtwRSYs0T26unI+MxWjXQRrO2zW0VTGkpP+YW7oRhTp\nIEap/Nrydyrh2tv+/N9HfR7bgCVNVc8n5YcbC++FxZruzwTrwma3eKrDw/XDWIjM\npKVugvNaTIZLC2TAX60KKRJ4BelLc+YaCrh86z88DEUsbV76xzsit8HzyqKX7B27\nMCsHXOrBAgMBAAECggEAEHdHW3rQt4jrVP78ZpWL05Bhi9Pa7bjTtTChfGoDouiq\nRiQDmwuoejKV8SvORt0iasWWQZ7DFzxaxJV8YLdSWH57l5RCxQW9+EbjxT+H6X49\nf/ZiqLQt6zN5rZQkqNe7cNr8xBhss9MZ9SPC2CY03ijTBxn40DV3hJUlqqDEmV6M\nrk5h7E60MNs2R/sIxhtAUX4K6vxPkgdD3o2T5W51k1j6SCVqGu8V7nW/D//IeBlJ\nVNTBb7lIFgKMcJiDNqKhq2tkULeRkwVc0TCCxXDZYozRw/oVn0TYfO18btSFz1uL\nnrtD9Ht0kJeIc46vuQHCQHbQ91nIWK6Po5Ve+8wUNwKBgQD1nCH7JGHk+7oYFcW2\nJ2S1nhdHIPxyQbnjVsG+ONEg9dXVS11NsjY96MJxtbO5OcI1RpzFE5yyosVUKFW3\njFLStpKXOLmcJibIOWRf8H+QbUZnM49yRisijV5jH/kwR6zetbyEGc7WqzxvjJIk\nWR08ykT80iKS+4l4EC6GNQ5wjwKBgQDsBnjdqvyz8SE8sBFMXe45UeFDuyLUEv05\njxGAhJkgMs2sOwympJjQnSRCvjxZLG9YHwDQL2q8UggtCTIpPRBRM11mxop7bPUq\nWJJiPBfPJiugpaklkoq4UoQUBigjkt3xALc+NFpfDgbQYcaxctJZjJRSUvUjRius\nxy6F19z3rwKBgHQeBNK/OKkRecG5SWf859gVjdvK9I7wE/ovIhnUsspqb1YP82Sw\nRISwbn1j8jw32mFlqOhjhUnPOou3Jg9JAD8uoc9suhPg1aUDvTi+cxDNGOPhtIfK\nNMp5G46xpxX0TP5d3Wp26RsEieYTB2S33OLIniUJE995nFxvCg/ZNaJxAoGAMIzw\nReDLVJRwWtR46nWT8FSIeu8+rdMuJa3pUr9z5CyvJBONeaX4DUmV0Oji7xD14nGW\nMDzgvtY8+k6e896swZdISkDi8ZqrH8fSbMShvSnD5arODX2EbYADzT6q+Q5X+yBD\nkVchk9YFzs2eGphc7rC9PeX0qQnhKAxc5IlP2d8CgYEA346YKTyutyTDxb11QKXW\nzwTqK3gPx9xryIHsUkbaIK1+esZAgeM/mznwGR/XpAlzP8rvy0b00N0Gmd4bGJls\n5VYu/O0NTMbPA8aLcDIxVq+S4pNKhDLcbvaDMOld9ieaX1O21Mbg0yTsS+XRBrm6\nwfhoYqMZIQ3pxDKn5NJhUs8=\n-----END PRIVATE KEY-----\n"
	PublicKeyDummy  = "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA4nIk8i0wSSV2UAFdQr5Z\nnCWOSYJ/j5sVkbtvxAWybMxOCz4qqAxv8Y2ICAtb3o4Mnnj6Sj3DBz/Kwr8huES1\n4lZwcCiF8BKj5idFflZjtF8HgTdIqQvg6riuUo6jTmywjbwhy0H7xQFpb/aLZaqx\nxvy00M61NJRbwWYbcEUmLNE9urpyPjMVo10Eazts1tFUxpKT/mFu6EYU6SBGqfza\n8ncq4drb/vzfR30e24AlTVXPJ+WHGwvvhcWa7s8E68Jmt3iqw8P1w1iIzKSlboLz\nWkyGSwtkwF+tCikSeAXpS3PmGgq4fOs/PAxFLG1e+sc7IrfB88qil+wduzArB1zq\nwQIDAQAB\n-----END PUBLIC KEY-----\n"
)

func TestNewClientMethodPing(t *testing.T) {
	ctx := context.Background()
	c := &http.Client{
		Timeout: 10 * time.Millisecond,
	}

	oauthServer := httptest.NewServer(func() http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			payload := `{
					"access_token": "c:xxxxxxxxxxxxx",
					"event_code": "",
					"expires_in": 86400,
					"last_login_type": "8",
					"sq_check": false,
					"token_type": "Bearer"
				}`
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(payload))
		}
	}())
	defer oauthServer.Close()

	sat := httptest.NewServer(func() http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == PING_PATH {
				payload := &PingResponse{
					Buildhash: "abc",
					Sandbox:   true,
					Status:    "ok",
				}
				res, errM := json.Marshal(payload)
				if errM != nil {
					w.Write([]byte(errM.Error()))
				}
				w.Write(res)
			}
		}
	}())
	defer sat.Close()

	cln, err := NewClient(
		"abc",
		"cde",
		PrivateKeyDummy,
		WithServerPublicKeyString(PublicKeyDummy),
		WithPaddingType(signature.PaddingTypePSS),
		WithIsDebug(true),
		WithHTTPClient(c),
		WithAccessTokenURL(oauthServer.URL+"/token"),
		WithSatBaseURL(sat.URL),
	)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(oauthServer.URL)
	fmt.Println(sat.URL)

	resp, err := cln.Ping(ctx)

	var errR APIResponseError
	ok := errors.As(err, &errR)
	if ok {
		fmt.Println("error: ", errR.Error())
	}

	if err != nil {
		fmt.Println(err)
	}

	want := &PingResponse{
		Buildhash: "abc",
		Sandbox:   true,
		Status:    "ok",
	}

	if !reflect.DeepEqual(resp, want) {
		t.Errorf("Ping() got = %v, want %v", resp, want)
	}

}

func TestNewClientMethodAccount(t *testing.T) {
	ctx := context.TODO()
	c := &http.Client{
		Timeout: 10 * time.Millisecond,
	}

	oauthServer := httptest.NewServer(func() http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			payload := `{
					"access_token": "c:xxxxxxxxxxxxx",
					"event_code": "",
					"expires_in": 86400,
					"last_login_type": "8",
					"sq_check": false,
					"token_type": "Bearer"
				}`
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(payload))
		}
	}())
	defer oauthServer.Close()

	sat := httptest.NewServer(func() http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == ACCOUNT_PATH {
				payload := &Account{
					ID:    123,
					Saldo: 100000,
				}
				errM := jsonapi.MarshalPayload(w, payload)
				if errM != nil {
					w.Write([]byte(errM.Error()))
				}
			}
		}
	}())
	defer sat.Close()

	cln, err := NewClient(
		"abc",
		"cde",
		PrivateKeyDummy,
		WithServerPublicKeyString(PublicKeyDummy),
		WithPaddingType(signature.PaddingTypePSS),
		WithIsDebug(true),
		WithHTTPClient(c),
		WithAccessTokenURL(oauthServer.URL+"/token"),
		WithSatBaseURL(sat.URL),
	)
	if err != nil {
		fmt.Println(err)
	}

	resp, err := cln.Account(ctx)

	var errR APIResponseError
	ok := errors.As(err, &errR)
	if ok {
		fmt.Println("error: ", errR.Error())
	}

	if err != nil {
		fmt.Println(err)
	}

	want := &Account{
		ID:    123,
		Saldo: 100000,
	}

	if !reflect.DeepEqual(resp, want) {
		t.Errorf("Account() got = %v, want %v", resp, want)
	}

}

func TestNewClientMethodInquiry(t *testing.T) {

	ctx := context.TODO()
	c := &http.Client{
		Timeout: 3000 * time.Millisecond,
	}

	oauthServer := httptest.NewServer(func() http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			payload := `{
					"access_token": "c:xxxxxxxxxxxxx",
					"event_code": "",
					"expires_in": 86400,
					"last_login_type": "8",
					"sq_check": false,
					"token_type": "Bearer"
				}`
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(payload))
		}
	}())
	defer oauthServer.Close()

	sat := httptest.NewServer(func() http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == INQUIRY_PATH {
				payload := &InquiryResponse{
					ID:          "client_number",
					ProductCode: "pln-prepaid-token-100k",
					SalesPrice:  12000,
					Fields: Fields{
						{
							Name:  "Name",
							Value: "Value",
						},
					},
					InquiryResult: nil,
					BasePrice:     10000,
					AdminFee:      2000,
					ClientName:    "client name",
					ClientNumber:  "client_number",
					MeterID:       "",
					RefID:         "",
					MaxPayment:    0,
					MinPayment:    0,
				}
				errM := jsonapi.MarshalPayload(w, payload)
				if errM != nil {
					w.Write([]byte(errM.Error()))
				}
			}
		}
	}())
	defer sat.Close()

	cln, err := NewClient(
		"abc",
		"cde",
		PrivateKeyDummy,
		WithServerPublicKeyString(PublicKeyDummy),
		WithPaddingType(signature.PaddingTypePSS),
		WithIsDebug(true),
		WithHTTPClient(c),
		WithAccessTokenURL(oauthServer.URL+"/token"),
		WithSatBaseURL(sat.URL),
	)
	if err != nil {
		fmt.Println(err)
	}

	resp, err := cln.Inquiry(ctx, &InquiryRequest{
		ID:           "client_number",
		ProductCode:  "pln-prepaid-token-100k",
		ClientNumber: "client_number",
		Amount:       0,
		Fields: Fields{
			{
				Name:  "Name",
				Value: "value",
			},
		},
	})

	var errR APIResponseError
	ok := errors.As(err, &errR)
	if ok {
		fmt.Println("error: ", errR.Error())
	}

	if err != nil {
		fmt.Println(err)
	}

	want := &InquiryResponse{
		ID:          "client_number",
		ProductCode: "pln-prepaid-token-100k",
		SalesPrice:  12000,
		Fields: Fields{
			{
				Name:  "Name",
				Value: "Value",
			},
		},
		InquiryResult: nil,
		BasePrice:     10000,
		AdminFee:      2000,
		ClientName:    "client name",
		ClientNumber:  "client_number",
		MeterID:       "",
		RefID:         "",
		MaxPayment:    0,
		MinPayment:    0,
	}

	if !reflect.DeepEqual(resp, want) {
		t.Errorf("Inquiry() got = %v, want %v", resp, want)
	}

}

func TestNewClientMethodCheckout(t *testing.T) {

	ctx := context.TODO()
	c := &http.Client{
		Timeout: 10 * time.Millisecond,
	}

	oauthServer := httptest.NewServer(func() http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			payload := `{
					"access_token": "c:xxxxxxxxxxxxx",
					"event_code": "",
					"expires_in": 86400,
					"last_login_type": "8",
					"sq_check": false,
					"token_type": "Bearer"
				}`
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(payload))
		}
	}())
	defer oauthServer.Close()

	sat := httptest.NewServer(func() http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == CHECKOUT_PATH {
				payload := &OrderDetail{
					RequestID: "request_id",
					Fields: Fields{
						{
							Name:  "Name",
							Value: "Value",
						},
					},
					FulfillmentResult: nil,
					FulfilledAt:       nil,
					ErrorCode:         "",
					ErrorDetail:       "",
					ProductCode:       "pln-prepaid-token-100k",
					Status:            "Pending",
					PartnerFee:        1000,
					SalesPrice:        12000,
					AdminFee:          2000,
					ClientName:        "",
					ClientNumber:      "",
					VoucherCode:       "",
					SerialNumber:      "",
				}
				errM := jsonapi.MarshalPayload(w, payload)
				if errM != nil {
					w.Write([]byte(errM.Error()))
				}
			}
		}
	}())
	defer sat.Close()

	cln, err := NewClient(
		"abc",
		"cde",
		PrivateKeyDummy,
		WithServerPublicKeyString(PublicKeyDummy),
		WithPaddingType(signature.PaddingTypePSS),
		WithIsDebug(true),
		WithHTTPClient(c),
		WithAccessTokenURL(oauthServer.URL+"/token"),
		WithSatBaseURL(sat.URL),
	)
	if err != nil {
		fmt.Println(err)
	}

	resp, err := cln.Checkout(ctx, &OrderRequest{
		RequestID:    "request_id",
		ProductCode:  "pln-prepaid-token-100k",
		ClientNumber: "client_number",
		Amount:       0,
		Fields: Fields{
			{
				Name:  "Name",
				Value: "value",
			},
		},
		DownlineID: "",
	})

	var errR APIResponseError
	ok := errors.As(err, &errR)
	if ok {
		fmt.Println("error: ", errR.Error())
	}

	if err != nil {
		fmt.Println(err)
	}

	want := &OrderDetail{
		RequestID: "request_id",
		Fields: Fields{
			{
				Name:  "Name",
				Value: "Value",
			},
		},
		FulfillmentResult: nil,
		FulfilledAt:       nil,
		ErrorCode:         "",
		ErrorDetail:       "",
		ProductCode:       "pln-prepaid-token-100k",
		Status:            "Pending",
		PartnerFee:        1000,
		SalesPrice:        12000,
		AdminFee:          2000,
		ClientName:        "",
		ClientNumber:      "",
		VoucherCode:       "",
		SerialNumber:      "",
	}

	if !reflect.DeepEqual(resp, want) {
		t.Errorf("Checkout() got = %v, want %v", resp, want)
	}

}

func TestNewClientMethodCheckStatus(t *testing.T) {
	ctx := context.TODO()
	c := &http.Client{
		Timeout: 10 * time.Millisecond,
	}

	oauthServer := httptest.NewServer(func() http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			payload := `{
					"access_token": "c:xxxxxxxxxxxxx",
					"event_code": "",
					"expires_in": 86400,
					"last_login_type": "8",
					"sq_check": false,
					"token_type": "Bearer"
				}`
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(payload))
		}
	}())
	defer oauthServer.Close()

	sat := httptest.NewServer(func() http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == fmt.Sprintf(CHECK_STATUS_PATH, "request_id") {
				payload := &OrderDetail{
					RequestID: "request_id",
					Fields: Fields{
						{
							Name:  "Name",
							Value: "Value",
						},
					},
					FulfillmentResult: Fields{
						{
							Name:  "Client Name",
							Value: "Mr X",
						},
						{
							Name:  "Admin Fee",
							Value: "Rp2000",
						},
						{
							Name:  "Total Bayar",
							Value: "Rp12.000",
						},
					},
					FulfilledAt:  nil,
					ErrorCode:    "",
					ErrorDetail:  "",
					ProductCode:  "pln-prepaid-token-100k",
					Status:       "Success",
					PartnerFee:   1000,
					SalesPrice:   12000,
					AdminFee:     2000,
					ClientName:   "",
					ClientNumber: "",
					VoucherCode:  "",
					SerialNumber: "",
				}
				sgn := signature.Init(signature.Options{
					PrivateKeyString: PrivateKeyDummy,
					PublicKeyString:  PublicKeyDummy,
					PaddingType:      0,
				})

				b := &bytes.Buffer{}
				payloadstr, err := jsonapi.Marshal(payload)
				if err != nil {
					fmt.Println("error: ", err)
				}

				err = json.NewEncoder(b).Encode(payloadstr)
				if err != nil {
					fmt.Println("error: ", err)
				}

				signt, err := sgn.Sign(b.Bytes())
				if err != nil {
					fmt.Println("error: ", err)
				}
				w.Header().Set(SIGNATURE_HEADER_KEY, signt)

				errM := jsonapi.MarshalPayload(w, payload)
				if errM != nil {
					w.Write([]byte(errM.Error()))
				}
			}
		}
	}())
	defer sat.Close()

	cln, err := NewClient(
		"abc",
		"cde",
		PrivateKeyDummy,
		WithServerPublicKeyString(PublicKeyDummy),
		WithPaddingType(signature.PaddingTypePSS),
		WithIsDebug(true),
		WithHTTPClient(c),
		WithAccessTokenURL(oauthServer.URL+"/token"),
		WithSatBaseURL(sat.URL),
	)
	if err != nil {
		fmt.Println(err)
	}

	resp, err := cln.CheckStatus(ctx, "request_id")

	var errR APIResponseError
	ok := errors.As(err, &errR)
	if ok {
		fmt.Println("error: ", errR.Error())
	}

	if err != nil {
		fmt.Println(err)
	}

	want := &OrderDetail{
		RequestID: "request_id",
		Fields: Fields{
			{
				Name:  "Name",
				Value: "Value",
			},
		},
		FulfillmentResult: Fields{
			{
				Name:  "Client Name",
				Value: "Mr X",
			},
			{
				Name:  "Admin Fee",
				Value: "Rp2000",
			},
			{
				Name:  "Total Bayar",
				Value: "Rp12.000",
			},
		},
		FulfilledAt:  nil,
		ErrorCode:    "",
		ErrorDetail:  "",
		ProductCode:  "pln-prepaid-token-100k",
		Status:       "Success",
		PartnerFee:   1000,
		SalesPrice:   12000,
		AdminFee:     2000,
		ClientName:   "",
		ClientNumber: "",
		VoucherCode:  "",
		SerialNumber: "",
	}

	if !reflect.DeepEqual(resp, want) {
		t.Errorf("CheckStatus() got = %v, want %v", resp, want)
	}

}

func TestNewClientMethodListProduct(t *testing.T) {
	ctx := context.TODO()
	c := &http.Client{
		Timeout: 10 * time.Millisecond,
	}

	oauthServer := httptest.NewServer(func() http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			payload := `{
					"access_token": "c:xxxxxxxxxxxxx",
					"event_code": "",
					"expires_in": 86400,
					"last_login_type": "8",
					"sq_check": false,
					"token_type": "Bearer"
				}`
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(payload))
		}
	}())
	defer oauthServer.Close()

	sat := httptest.NewServer(func() http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			if req.URL.Path == PRODUCT_LIST_PATH {
				payload := []*Product{
					{
						Name:         "PLN Prepaid",
						Code:         "pln-prepaid-token-100k",
						OperatorName: "Token Listrik",
						CategoryName: "Listrik PLN",
						IsInquiry:    true,
						SalesPrice:   100000,
						Status:       1,
						ClientNumber: "testing",
					},
				}
				errM := jsonapi.MarshalPayload(w, payload)
				if errM != nil {
					w.Write([]byte(errM.Error()))
				}
			}
		}
	}())
	defer sat.Close()

	cln, err := NewClient(
		"abc",
		"cde",
		PrivateKeyDummy,
		WithServerPublicKeyString(PublicKeyDummy),
		WithPaddingType(signature.PaddingTypePSS),
		WithIsDebug(true),
		WithHTTPClient(c),
		WithAccessTokenURL(oauthServer.URL+"/token"),
		WithSatBaseURL(sat.URL),
	)
	if err != nil {
		fmt.Println(err)
	}

	resp, err := cln.ListProduct(ctx, "pln-prepaid-token-100k")

	var errR APIResponseError
	ok := errors.As(err, &errR)
	if ok {
		fmt.Println("error: ", errR.Error())
	}

	if err != nil {
		fmt.Println(err)
	}

	want := []*Product{
		{
			Name:         "PLN Prepaid",
			Code:         "pln-prepaid-token-100k",
			OperatorName: "Token Listrik",
			CategoryName: "Listrik PLN",
			IsInquiry:    true,
			SalesPrice:   100000,
			Status:       1,
			ClientNumber: "testing",
		},
	}

	if !reflect.DeepEqual(resp, want) {
		t.Errorf("ListProduct() got = %v, want %v", resp, want)
	}

}

type callbackExample struct {
	t *testing.T
}

func (c *callbackExample) Do(ctx context.Context, request *OrderDetail) error {
	fmt.Println("CALLBACK: ", request)
	want := OrderDetail{
		RequestID: "request_id",
		Fields: Fields{
			{
				Name:  "Name",
				Value: "Value",
			},
		},
		FulfillmentResult: Fields{
			{
				Name:  "Client Name",
				Value: "Mr X",
			},
			{
				Name:  "Admin Fee",
				Value: "Rp2000",
			},
			{
				Name:  "Total Bayar",
				Value: "Rp12.000",
			},
		},
		FulfilledAt:  nil,
		ErrorCode:    "",
		ErrorDetail:  "",
		ProductCode:  "pln-prepaid-token-100k",
		Status:       "Success",
		PartnerFee:   1000,
		SalesPrice:   12000,
		AdminFee:     2000,
		ClientName:   "",
		ClientNumber: "",
		VoucherCode:  "",
		SerialNumber: "",
	}

	if reflect.DeepEqual(request, want) {
		c.t.Errorf("Do() got = %v, want %v", request, want)
	}

	return nil
}

func TestNewClientMethodCallback(t *testing.T) {
	ctx := context.TODO()
	c := &http.Client{
		Timeout: 10 * time.Millisecond,
	}

	oauthServer := httptest.NewServer(func() http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			payload := `{
					"access_token": "c:xxxxxxxxxxxxx",
					"event_code": "",
					"expires_in": 86400,
					"last_login_type": "8",
					"sq_check": false,
					"token_type": "Bearer"
				}`
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(200)
			w.Write([]byte(payload))
		}
	}())
	defer oauthServer.Close()

	cln, err := NewClient(
		"abc",
		"cde",
		PrivateKeyDummy,
		WithServerPublicKeyString(PublicKeyDummy),
		WithPaddingType(signature.PaddingTypePSS),
		WithIsDebug(true),
		WithHTTPClient(c),
		WithAccessTokenURL(oauthServer.URL+"/token"),
	)
	if err != nil {
		fmt.Println(err)
	}

	clbx := &callbackExample{
		t: t,
	}
	s := httptest.NewServer(cln.HandleCallback(clbx))
	defer s.Close()
	fmt.Println("Callback URL: ", s.URL)

	payload := &OrderDetail{
		RequestID: "request_id",
		Fields: Fields{
			{
				Name:  "Name",
				Value: "Value",
			},
		},
		FulfillmentResult: Fields{
			{
				Name:  "Client Name",
				Value: "Mr X",
			},
			{
				Name:  "Admin Fee",
				Value: "Rp2000",
			},
			{
				Name:  "Total Bayar",
				Value: "Rp12.000",
			},
		},
		FulfilledAt:  nil,
		ErrorCode:    "",
		ErrorDetail:  "",
		ProductCode:  "pln-prepaid-token-100k",
		Status:       "Success",
		PartnerFee:   1000,
		SalesPrice:   12000,
		AdminFee:     2000,
		ClientName:   "",
		ClientNumber: "",
		VoucherCode:  "",
		SerialNumber: "",
	}
	bd := &bytes.Buffer{}
	err = jsonapi.MarshalPayload(bd, payload)
	if err != nil {
		fmt.Println("error: ", err)
	}

	pb := bd.Bytes()

	cd := http.DefaultClient
	requestCallback, err := http.NewRequestWithContext(ctx, http.MethodPost, s.URL, bd)
	if err != nil {
		fmt.Println("error: ", err)
	}

	sgn := signature.Init(signature.Options{
		PrivateKeyString: PrivateKeyDummy,
		PublicKeyString:  PublicKeyDummy,
		PaddingType:      0,
	})

	signt, err := sgn.Sign(pb)
	if err != nil {
		fmt.Println("error: ", err)
	}

	requestCallback.Header.Add(SIGNATURE_HEADER_KEY, signt)

	resp, err := cd.Do(requestCallback)
	if err != nil {
		fmt.Println("error: ", err)
	}

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error: ", err)
	}

	want := "OK"

	if !reflect.DeepEqual(string(res), want) {
		t.Errorf("Callback() got = %v, want %v", string(res), want)
	}

}

func TestNewClient(t *testing.T) {
	type args struct {
		clientID     string
		clientSecret string
		privateKey   string
		opts         []ClientOptionFunc
	}

	dummyClient := &http.Client{}

	tests := []struct {
		name    string
		args    args
		want    *Client
		wantErr bool
	}{
		{
			name: "init success",
			args: args{
				clientID:     "abc",
				clientSecret: "def",
				privateKey:   "priv key",
				opts: []ClientOptionFunc{
					WithIsDebug(true),
					WithHTTPClient(dummyClient),
					WithServerPublicKeyString("pub key"),
				},
			},
			want: &Client{
				http:           dummyClient,
				logger:         nil,
				satBaseURL:     PLAYGROUND_SAT_BASE_URL,
				accessTokenURL: ACCESS_TOKEN_URL,
				isDebug:        true,
				signature: signature.Init(signature.Options{
					PrivateKeyString: "priv key",
					PublicKeyString:  "pub key",
					PaddingType:      0,
				}),
			},
		},
		{
			name: "failed because empty client id",
			args: args{
				clientID:     "",
				clientSecret: "def",
				privateKey:   "priv key",
				opts: []ClientOptionFunc{
					WithIsDebug(true),
					WithHTTPClient(dummyClient),
					WithServerPublicKeyString("pub key"),
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "failed because empty client secret",
			args: args{
				clientID:     "abc",
				clientSecret: "",
				privateKey:   "priv key",
				opts: []ClientOptionFunc{
					WithIsDebug(true),
					WithHTTPClient(dummyClient),
					WithServerPublicKeyString("pub key"),
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "failed because empty private key",
			args: args{
				clientID:     "abc",
				clientSecret: "def",
				privateKey:   "",
				opts: []ClientOptionFunc{
					WithIsDebug(true),
					WithHTTPClient(dummyClient),
					WithServerPublicKeyString("pub key"),
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.args.clientID, tt.args.clientSecret, tt.args.privateKey, tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != nil {
				got.logger = nil
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClient() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_handleErrorResponse(t *testing.T) {
	type fields struct {
		http           *http.Client
		logger         *log.Logger
		satBaseURL     string
		accessTokenURL string
		signature      *signature.Signature
		isDebug        bool
	}
	type args struct {
		resp *http.Response
	}

	htmlError := "<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n    <meta charset=\"UTF-8\">\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n    <title>500 Internal Server Error</title>\n    <style>\n        body {\n            font-family: Arial, sans-serif;\n            text-align: center;\n            padding: 50px;\n        }\n        h1 {\n            font-size: 50px;\n        }\n        p {\n            font-size: 20px;\n        }\n    </style>\n</head>\n<body>\n    <h1>500 Internal Server Error</h1>\n    <p>Oops! Something went wrong on our end. Please try again later.</p>\n</body>\n</html>\n"

	errorResp := &http.Response{
		Status:     http.StatusText(500),
		StatusCode: 500,
		Proto:      "",
		ProtoMajor: 0,
		ProtoMinor: 0,
		Header: http.Header{
			"Content-Type": []string{"text/html"},
		},
		Body: func() io.ReadCloser {
			reader := strings.NewReader(htmlError)
			return io.NopCloser(reader)
		}(),
		ContentLength:    123,
		TransferEncoding: nil,
		Close:            false,
		Uncompressed:     false,
		Trailer:          nil,
		Request:          nil,
		TLS:              nil,
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name:   "handle error standard",
			fields: fields{},
			args: args{
				resp: &http.Response{
					Status:     http.StatusText(500),
					StatusCode: 500,
					Proto:      "",
					ProtoMajor: 0,
					ProtoMinor: 0,
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
					Body: func() io.ReadCloser {
						reader := strings.NewReader("{\"errors\":[{\"detail\":\"Internal Server Error\",\"status\":\"500\",\"code\":\"S00\"}]}")
						return io.NopCloser(reader)
					}(),
					ContentLength:    76,
					TransferEncoding: nil,
					Close:            false,
					Uncompressed:     false,
					Trailer:          nil,
					Request:          nil,
					TLS:              nil,
				},
			},
			wantErr: &ErrorResponse{Errors: []*ErrorObject{
				{
					Detail: "Internal Server Error",
					Status: "500",
					Code:   "S00",
				},
			}},
		},
		{
			name:   "handle error html page, just return as is",
			fields: fields{},
			args: args{
				resp: errorResp,
			},
			wantErr: &InternalError{resp: errorResp},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				http:           tt.fields.http,
				logger:         tt.fields.logger,
				satBaseURL:     tt.fields.satBaseURL,
				accessTokenURL: tt.fields.accessTokenURL,
				signature:      tt.fields.signature,
				isDebug:        tt.fields.isDebug,
			}

			err := c.handleErrorResponse(tt.args.resp)
			if !reflect.DeepEqual(err, tt.wantErr) {
				t.Errorf("handleErrorResponse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
