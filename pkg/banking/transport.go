package banking

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/log"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
	"github.com/twonegatives/coinsph_challenge/pkg/entities"
)

type GetAccountsRequest struct {
}

type GetPaymentsRequest struct {
}

type SendPaymentRequest struct {
	From   entities.Account
	To     entities.Account
	Amount decimal.Decimal
}

func decodeGetAccountsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return GetAccountsRequest{}, nil
}

func decodeGetPaymentsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return GetPaymentsRequest{}, nil
}

func decodeSendPaymentRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return SendPaymentRequest{}, nil
}

func MakeHandler(svc BankingService, l log.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		//kithttp.ServerErrorEncoder(errorEncoder),
		kithttp.ServerErrorLogger(l),
	}

	getAccounts := kithttp.NewServer(
		MakeGetAccountsEndpoint(svc),
		decodeGetAccountsRequest,
		kithttp.EncodeJSONResponse,
		opts...,
	)

	getPayments := kithttp.NewServer(
		MakeGetPaymentsEndpoint(svc),
		decodeGetPaymentsRequest,
		kithttp.EncodeJSONResponse,
		opts...,
	)

	sendPayment := kithttp.NewServer(
		MakeSendPaymentEndpoint(svc),
		decodeSendPaymentRequest,
		kithttp.EncodeJSONResponse,
		opts...,
	)

	m := mux.NewRouter()
	m.Handle("/accounts", getAccounts).Methods(http.MethodGet)
	m.Handle("/payments", getPayments).Methods(http.MethodGet)
	m.Handle("/payments", sendPayment).Methods(http.MethodPost)
	//m.NotFoundHandler = NotFoundHandler(l)
	return m
}
