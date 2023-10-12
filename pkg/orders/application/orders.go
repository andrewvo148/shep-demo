package application

import (
	"log"
	"shop-demo/pkg/common/price"
	"shop-demo/pkg/orders/domain/orders"

	"github.com/pkg/errors"
)

type productsService interface {
	ProductByID(id orders.ProductID) (orders.Product, error)
}

type paymentsService interface {
	InitializeOrderPayment(id orders.ID, price price.Price) error
}

type OrdersService struct {
	productsService productsService
	paymentsService paymentsService

	orderRepository orders.Repository
}

type PlaceOrderCommandAddress struct {
	Name     string
	Street   string
	City     string
	PostCode string
	Country  string
}

type PlaceOrderCommand struct {
	OrderID   orders.ID
	ProductID orders.ProductID

	Address PlaceOrderCommandAddress
}

func (s OrdersService) PlaceOrder(cmd PlaceOrderCommand) error {
	address, err := orders.NewAddress(
		cmd.Address.Name,
		cmd.Address.Street,
		cmd.Address.City,
		cmd.Address.PostCode,
		cmd.Address.Country,
	)
	if err != nil {
		return errors.Wrap(err, "invalid address")
	}

	product, err := s.productsService.ProductByID(cmd.ProductID)
	if err != nil {
		return errors.Wrap(err, "cannot get product")
	}

	newOrder, err := orders.NewOrder(cmd.OrderID, product, address)
	if err != nil {
		return errors.Wrap(err, "cannot save order")
	}

	if err := s.orderRepository.Save(newOrder); err != nil {
		return errors.Wrap(err, "cannot save order")
	}

	if err := s.paymentsService.InitializeOrderPayment(newOrder.ID(), newOrder.Product().Price()); err != nil {
		return errors.Wrap(err, "cannot initialize payment")
	}

	log.Printf("order % placed", cmd.OrderID)
	return nil
}
