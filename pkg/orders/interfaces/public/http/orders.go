package http

import (
	"net/http"
	"shop-demo/pkg/orders/application"
	"shop-demo/pkg/orders/domain/orders"
	common_http "shop-demo/pkg/orders/pkg/common/http"

	"github.com/go-chi/render"
	"github.com/gofrs/uuid"
)

type ordersResource struct {
	service    application.OrdersService
	repository orders.Repository
}

func (o ordersResource) Post(w http.ResponseWriter, r *http.Request) {
	req := PostOrderRequest{}
	if err := render.Decode(r, &req); err != nil {
		_ = render.Render(w, r, common_http.ErrBadRequest(err))
		return
	}

	cmd := application.PlaceOrderCommand{
		OrderID:   orders.ID(uuid.NewV1().String()),
		ProductID: req.ProductID,
		Address:   application.PlaceOrderCommandAddress(req.Address),
	}

	if err := o.service.PlaceOrder(cmd); err != nil {
		_ = render.Render(w, r, common_http.ErrInternal(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, PostOrdersResponse{
		OrderID: string(cmd.OrderID),
	})
}

type PostOrderAddress struct {
	Name     string `json:"name"`
	Street   string `json:"street"`
	City     string `json:"city"`
	PostCode string `json:"postCode"`
	Country  string `json:"country"`
}

type PostOrderRequest struct {
	ProductID orders.ProductID `json:"product_id"`
	Address   PostOrderAddress `json:"address"`
}

type PostOrdersResponse struct {
	OrderID string
}
