package main

import (
	"context"
	"fmt"

	"github.com/tokopedia/digital-b2b-client-library/golang-sat"
)

type callbackExample struct{}

func (c *callbackExample) Do(ctx context.Context, request *sat.OrderDetail) error {
	fmt.Println("CALLBACK Payload: ", request)
	// Do something
	return nil
}
