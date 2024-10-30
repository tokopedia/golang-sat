package sat

// Account is a schema related with account entity
type Account struct {
	ID    int64 `jsonapi:"primary,account"`
	Saldo int64 `jsonapi:"attr,saldo"`
}
