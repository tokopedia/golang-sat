package sat

// ProductStatus is a product status type
type ProductStatus int

const (
	// ProductStatusActive is for product active
	ProductStatusActive ProductStatus = 1
	// ProductStatusInactive is for product inactive for some reason, like business reason
	ProductStatusInactive ProductStatus = 2
	// ProductTempInactive is for product inactive caused by some problem on supplier
	ProductTempInactive ProductStatus = 3
)

// Product is a schema for product will be used by partner to identify how many products they have and allowed
type Product struct {
	Name         string        `jsonapi:"attr,product_name"`
	Code         string        `jsonapi:"primary,product"`
	OperatorName string        `jsonapi:"attr,operator_name,omitempty"`
	CategoryName string        `jsonapi:"attr,category_name,omitempty"`
	IsInquiry    bool          `jsonapi:"attr,is_inquiry"`
	SalesPrice   int64         `jsonapi:"attr,price"`
	Status       ProductStatus `jsonapi:"attr,status"`
	ClientNumber string        `jsonapi:"attr,client_number,omitempty"`
}
