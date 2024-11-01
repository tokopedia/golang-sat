package main

import (
	"context"
	"fmt"

	sat "github.com/tokopedia/golang-sat"
)

type callbackExample struct{}

func (c *callbackExample) Do(ctx context.Context, request *sat.OrderDetail) error {
	fmt.Println("CALLBACK Payload: ", request)
	// Do something
	return nil
}
