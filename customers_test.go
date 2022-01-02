
package housecall 

import (
	"github.com/stretchr/testify/assert"

	"testing"
	"context"
	"time"
)

func TestThirdCustomers (t *testing.T) {
	hc, cfg := newHouseCall (t)

	ctx, cancel := context.WithTimeout (context.Background(), time.Minute) // this should take < 1 minute
	defer cancel()

	// get our list of customers
	customers, err := hc.SearchCustomers (ctx, cfg.Token, "05445")
	if err != nil { t.Fatal (err) }

	assert.Equal (t, true, len(customers) > 0, "expecting at least 1 customer")
	assert.Equal (t, true, len(customers[0].Id) > 0)
	assert.Equal (t, true, len(customers[0].Addresses) > 0)

	// try a specific one
	customers, err = hc.SearchCustomers (ctx, cfg.Token, "2 COMMON WAY")
	if err != nil { t.Fatal (err) }

	assert.Equal (t, 1, len(customers), "expecting 1 customer")
	assert.Equal (t, true, len(customers[0].Id) > 0)
	assert.Equal (t, true, len(customers[0].Addresses) > 0)
	assert.Equal (t, "Louisa", customers[0].FirstName)
	assert.Equal (t, "Adams", customers[0].LastName)

}

func TestThirdCustomerCreate (t *testing.T) {
	hc, cfg := newHouseCall (t)

	ctx, cancel := context.WithTimeout (context.Background(), time.Minute) // this should take < 1 minute
	defer cancel()

	customer := &Customer {
		FirstName: "Mayor",
		LastName: "Burlington",
	}

	customer.Addresses = append (customer.Addresses, Address {
		Street: "149 Church St",
		City: "Burlington",
		State: "VT",
		Zip: "05401",
	})

	err := hc.CreateCustomer (ctx, cfg.Token, customer)
	if err != nil { t.Fatal (err) }

	assert.Equal (t, true, len(customer.Id) > 0)
}

