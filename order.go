package sat

import "time"

// OrderRequest contains order request payload
type OrderRequest struct {
	RequestID    string `jsonapi:"primary,order"`
	ProductCode  string `jsonapi:"attr,product_code"`
	ClientNumber string `jsonapi:"attr,client_number"`
	Amount       int64  `jsonapi:"attr,amount"`
	Fields       Fields `jsonapi:"attr,fields"`
	DownlineID   string `jsonapi:"attr,downline_id"`
}

// OrderDetail contains order detail information
type OrderDetail struct {
	RequestID         string     `jsonapi:"primary,order"`
	Fields            Fields     `jsonapi:"attr,fields"`
	FulfillmentResult Fields     `jsonapi:"attr,fulfillment_result"`
	FulfilledAt       *time.Time `jsonapi:"attr,fulfilled_at,iso8601"`
	ErrorCode         string     `jsonapi:"attr,error_code"`
	ErrorDetail       string     `jsonapi:"attr,error_detail"`
	ProductCode       string     `jsonapi:"attr,product_code"`
	Status            string     `jsonapi:"attr,status"`
	PartnerFee        int64      `jsonapi:"attr,partner_fee"`
	SalesPrice        int64      `jsonapi:"attr,sales_price"`
	AdminFee          int64      `jsonapi:"attr,admin_fee"`
	ClientName        string     `jsonapi:"attr,client_name"`
	ClientNumber      string     `jsonapi:"attr,client_number"`
	VoucherCode       string     `jsonapi:"attr,voucher_code"`
	SerialNumber      string     `jsonapi:"attr,serial_number"`
}
