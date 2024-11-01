package sat

import (
	"encoding/json"
)

// InquiryRequest to hold inquiry request
type InquiryRequest struct {
	ID           string `jsonapi:"primary,inquiry"`
	ProductCode  string `jsonapi:"attr,product_code"`
	ClientNumber string `jsonapi:"attr,client_number"`
	Amount       int64  `jsonapi:"attr,amount"`
	Fields       Fields `jsonapi:"attr,fields"`
	DownlineID   string `jsonapi:"attr,downline_id"`
}

// InquiryResponse to hold inquiry response
type InquiryResponse struct {
	ID            string `jsonapi:"primary,inquiry"`
	ProductCode   string `jsonapi:"attr,product_code"`
	SalesPrice    int64  `jsonapi:"attr,sales_price"`
	Fields        Fields `jsonapi:"attr,fields"`
	InquiryResult Fields `jsonapi:"attr,inquiry_result"`
	BasePrice     int64  `jsonapi:"attr,base_price"`
	AdminFee      int64  `jsonapi:"attr,admin_fee"`
	ClientName    string `jsonapi:"attr,client_name"`
	ClientNumber  string `jsonapi:"attr,client_number"`
	MeterID       string `jsonapi:"attr,meter_id"`
	RefID         string `jsonapi:"attr,ref_id,omitempty"`
	MaxPayment    int64  `jsonapi:"attr,max_payment,omitempty"`
	MinPayment    int64  `jsonapi:"attr,min_payment,omitempty"`
	MinAmount     int64  `jsonapi:"attr,min_amount,omitempty"`
}

// Fields contains dynamic key value
type Fields []Field

// Field contains field detail information
type Field struct {
	Name  string `json:"name" jsonapi:"attr,name"`
	Value string `json:"value" jsonapi:"attr,value"`
}

// UnmarshalJSON override unmarshal jsonapi
func (s *Field) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, s)
	if err != nil {
		return err
	}

	return nil
}
