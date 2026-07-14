package evaluate

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/squeakycheese75/open-incentives/internal/auth"
	"github.com/squeakycheese75/open-incentives/internal/domain"
	"github.com/squeakycheese75/open-incentives/internal/httputil"
)

type EvaluateRequest struct {
	Customer CustomerRequest `json:"customer"`
	Cart     CartRequest     `json:"cart"`
}

type CustomerRequest struct {
	ID      string `json:"id"`
	Country string `json:"country"`
	Tier    string `json:"tier"`
}

type CartRequest struct {
	Subtotal float64           `json:"subtotal"`
	Currency string            `json:"currency"`
	Items    []CartItemRequest `json:"items,omitempty"`
}

type CartItemRequest struct {
	ProductID string  `json:"productId"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unitPrice"`
}

type EvaluateResponse struct {
	Decision  DecisionResponse   `json:"decision"`
	Cart      CartResponse       `json:"cart"`
	Discounts []DiscountResponse `json:"discounts"`
}

type DecisionResponse struct {
	Matched          bool `json:"matched"`
	CampaignsMatched int  `json:"campaignsMatched"`
}

type CartResponse struct {
	Subtotal      float64 `json:"subtotal"`
	DiscountTotal float64 `json:"discountTotal"`
	Total         float64 `json:"total"`
	Currency      string  `json:"currency"`
}

type DiscountResponse struct {
	CampaignID   string  `json:"campaignId"`
	CampaignName string  `json:"campaignName"`
	RuleID       string  `json:"ruleId"`
	Type         string  `json:"type"`
	Amount       float64 `json:"amount"`
}

func (h *Handler) Evaluate(w http.ResponseWriter, r *http.Request) {
	scope, ok := auth.EvalAuthFromContext(r.Context())
	if !ok {
		httputil.WriteJSON(w, http.StatusUnauthorized, map[string]any{
			"error": "unauthorized",
		})
		return
	}

	var req EvaluateRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&req); err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"error": "invalid_json",
		})
		return
	}

	if err := validateRequest(req); err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, map[string]any{
			"error":   "invalid_request",
			"message": err.Error(),
		})
		return
	}

	output, err := h.evalContainer.EvaluateUsecase(scope.OrgID).Execute(r.Context(), domain.EvaluateUsecaseInput{
		ProjectID: scope.ProjectID,
		Customer: domain.Customer{
			ID:      req.Customer.ID,
			Country: req.Customer.Country,
			Tier:    req.Customer.Tier,
		},
		Cart: resolveCart(req.Cart),
	})
	if err != nil {
		slog.ErrorContext(
			r.Context(),
			"evaluation failed",
			"error", err,
		)

		httputil.WriteJSON(w, http.StatusInternalServerError, map[string]any{
			"error": "evaluation_failed",
		})
		return
	}

	httputil.WriteJSON(
		w,
		http.StatusOK,
		toEvaluateResponse(output),
	)
}

func validateRequest(req EvaluateRequest) error {
	if req.Customer.ID == "" {
		return errors.New("customer.id is required")
	}

	if req.Cart.Currency == "" {
		return errors.New("cart.currency is required")
	}

	if len(req.Cart.Items) > 0 {
		for i, item := range req.Cart.Items {
			if item.ProductID == "" {
				return fmt.Errorf("cart.items[%d].productId is required", i)
			}

			if item.Quantity <= 0 {
				return fmt.Errorf("cart.items[%d].quantity must be greater than zero", i)
			}

			if item.UnitPrice < 0 {
				return fmt.Errorf("cart.items[%d].unitPrice must not be negative", i)
			}
		}

		return nil
	}

	if req.Cart.Subtotal < 0 {
		return errors.New("cart.subtotal must not be negative")
	}

	return nil
}

func resolveCart(req CartRequest) domain.Cart {
	if len(req.Items) == 0 {
		return domain.Cart{
			Subtotal: req.Subtotal,
			Currency: req.Currency,
		}
	}

	items := make([]domain.CartItem, 0, len(req.Items))
	subtotal := 0.0

	for _, item := range req.Items {
		items = append(items, domain.CartItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		})

		subtotal += float64(item.Quantity) * item.UnitPrice
	}

	return domain.Cart{
		Subtotal: subtotal,
		Currency: req.Currency,
		Items:    items,
	}
}

func toEvaluateResponse(
	output domain.EvaluateUsecaseOutput,
) EvaluateResponse {
	discounts := make(
		[]DiscountResponse,
		0,
		len(output.Discounts),
	)

	for _, discount := range output.Discounts {
		discounts = append(discounts, DiscountResponse{
			CampaignID:   discount.CampaignID,
			CampaignName: discount.CampaignName,
			RuleID:       discount.RuleID,
			Type:         discount.Type,
			Amount:       discount.Amount,
		})
	}

	return EvaluateResponse{
		Decision: DecisionResponse{
			Matched:          output.Decision.Matched,
			CampaignsMatched: output.Decision.CampaignsMatched,
		},
		Cart: CartResponse{
			Subtotal:      output.Cart.Subtotal,
			DiscountTotal: output.Cart.DiscountTotal,
			Total:         output.Cart.Total,
			Currency:      output.Cart.Currency,
		},
		Discounts: discounts,
	}
}
