package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/tokopedia/digital-b2b-client-library/golang-sat"
)

type integrationExample struct {
	client *sat.Client
}

func (i *integrationExample) Handle(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	// checking product before doing inquiry or checkout is optional.
	// this implementation is only for sample to trigger a get list of product via SDK.
	// ListProduct will be beneficial for sync product status.
	// make your product status always up to date
	resProductList, err := i.client.ListProduct(ctx, "pln-prepaid-token-50k-sat")

	var errR sat.APIResponseError
	ok := errors.As(err, &errR)
	if ok {
		fmt.Println("[PRODUCT LIST] error: ", errR.Error())
	}

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("[PRODUCT LIST] response: ", resProductList)

	resInq, err := i.client.Inquiry(ctx, &sat.InquiryRequest{
		ProductCode:  "pln-prepaid-token-50k-sat",
		ClientNumber: "102111106111",
	})

	ok = errors.As(err, &errR)
	if ok {
		fmt.Println("[INQUIRY] error: ", errR.Error())
		// handle the error
		// read the documentation
		switch errR.Code() {
		case "S00":
			// do something
		case "P00":
			// do something
		case "U00":
			// do something
		}
	}

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("[INQUIRY] response: ", resInq)

	reqID := strconv.Itoa(rand.Int())

	resOrder, err := i.client.Checkout(ctx, &sat.OrderRequest{
		ProductCode:  "pln-prepaid-token-100k",
		ClientNumber: "102111106111",
		RequestID:    reqID,
	})

	ok = errors.As(err, &errR)
	if ok {
		fmt.Println("[CHECKOUT] error: ", errR.Error())
		// handle the error
		// read the documentation
		switch errR.Code() {
		case "S00":
			// do something
		case "P00":
			// do something
		case "U00":
			// do something
		}
	}

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("[CHECKOUT] response: ", resOrder)
	fmt.Println("WAITING 3 Seconds...")
	time.Sleep(3 * time.Second)

	resOrderDetail, err := i.client.CheckStatus(ctx, reqID)

	ok = errors.As(err, &errR)
	if ok {
		fmt.Println("[CHECK STATUS] error: ", errR.Error())
		// handle the error
		// read the documentation
		switch errR.Code() {
		case "S00":
			// do something
		case "P00":
			// do something
		case "U00":
			// do something
		}
	}

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("[CHECK STATUS] response: %v", resOrderDetail)

	if resOrderDetail != nil && resOrderDetail.Status == "Failed" {
		switch resOrderDetail.ErrorCode {
		case "S00":
			// do something
		case "P00":
			// do something
		case "U00":
			// do something
		}
	}

	res, err := json.Marshal(resOrderDetail)
	if err != nil {
		fmt.Println(err)
	}

	writer.Write(res)
	return
}
