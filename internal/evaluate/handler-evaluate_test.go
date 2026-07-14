package evaluate

import "testing"

func TestValidateRequest_SubtotalOnlyPath(t *testing.T) {
	tests := []struct {
		name    string
		req     EvaluateRequest
		wantErr bool
	}{
		{
			name: "valid subtotal-only request",
			req: EvaluateRequest{
				Customer: CustomerRequest{ID: "cust_1"},
				Cart:     CartRequest{Subtotal: 50, Currency: "EUR"},
			},
			wantErr: false,
		},
		{
			name: "missing customer id",
			req: EvaluateRequest{
				Cart: CartRequest{Subtotal: 50, Currency: "EUR"},
			},
			wantErr: true,
		},
		{
			name: "missing currency",
			req: EvaluateRequest{
				Customer: CustomerRequest{ID: "cust_1"},
				Cart:     CartRequest{Subtotal: 50},
			},
			wantErr: true,
		},
		{
			name: "negative subtotal",
			req: EvaluateRequest{
				Customer: CustomerRequest{ID: "cust_1"},
				Cart:     CartRequest{Subtotal: -1, Currency: "EUR"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRequest(tt.req)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateRequest_ItemsPath(t *testing.T) {
	tests := []struct {
		name    string
		req     EvaluateRequest
		wantErr bool
	}{
		{
			name: "valid items",
			req: EvaluateRequest{
				Customer: CustomerRequest{ID: "cust_1"},
				Cart: CartRequest{
					Currency: "EUR",
					Items: []CartItemRequest{
						{ProductID: "prod_coffee", Quantity: 2, UnitPrice: 18},
						{ProductID: "prod_mug", Quantity: 1, UnitPrice: 14},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing product id",
			req: EvaluateRequest{
				Customer: CustomerRequest{ID: "cust_1"},
				Cart: CartRequest{
					Currency: "EUR",
					Items:    []CartItemRequest{{Quantity: 1, UnitPrice: 18}},
				},
			},
			wantErr: true,
		},
		{
			name: "zero quantity",
			req: EvaluateRequest{
				Customer: CustomerRequest{ID: "cust_1"},
				Cart: CartRequest{
					Currency: "EUR",
					Items:    []CartItemRequest{{ProductID: "prod_coffee", Quantity: 0, UnitPrice: 18}},
				},
			},
			wantErr: true,
		},
		{
			name: "negative unit price",
			req: EvaluateRequest{
				Customer: CustomerRequest{ID: "cust_1"},
				Cart: CartRequest{
					Currency: "EUR",
					Items:    []CartItemRequest{{ProductID: "prod_coffee", Quantity: 1, UnitPrice: -1}},
				},
			},
			wantErr: true,
		},
		{
			name: "items path ignores client subtotal being unset",
			req: EvaluateRequest{
				Customer: CustomerRequest{ID: "cust_1"},
				Cart: CartRequest{
					Subtotal: 0,
					Currency: "EUR",
					Items:    []CartItemRequest{{ProductID: "prod_coffee", Quantity: 1, UnitPrice: 18}},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRequest(tt.req)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestResolveCart_ComputesSubtotalFromItems(t *testing.T) {
	req := CartRequest{
		Subtotal: 999, // must be ignored in favor of items-derived subtotal
		Currency: "EUR",
		Items: []CartItemRequest{
			{ProductID: "prod_coffee", Quantity: 2, UnitPrice: 18},
			{ProductID: "prod_mug", Quantity: 1, UnitPrice: 14},
		},
	}

	cart := resolveCart(req)

	const wantSubtotal = 50.0
	if cart.Subtotal != wantSubtotal {
		t.Fatalf("cart.Subtotal = %v, want %v", cart.Subtotal, wantSubtotal)
	}

	if cart.Currency != "EUR" {
		t.Fatalf("cart.Currency = %q, want %q", cart.Currency, "EUR")
	}

	if len(cart.Items) != 2 {
		t.Fatalf("len(cart.Items) = %d, want 2", len(cart.Items))
	}
}

func TestResolveCart_FallsBackToSubtotalWhenNoItems(t *testing.T) {
	req := CartRequest{Subtotal: 50, Currency: "EUR"}

	cart := resolveCart(req)

	if cart.Subtotal != 50 {
		t.Fatalf("cart.Subtotal = %v, want 50", cart.Subtotal)
	}

	if cart.Items != nil {
		t.Fatalf("cart.Items = %v, want nil", cart.Items)
	}
}
