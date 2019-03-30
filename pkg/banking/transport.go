package banking

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/twonegatives/coinsph_challenge/pkg/entities"
)

var (
	errBadRequest = errors.New("bad request")
)

type payment struct {
	From   string          `json:"from"`
	To     string          `json:"to"`
	Amount decimal.Decimal `json:"amount"`
}

type sendPaymentBody struct {
	Payment payment `json:"payment"`
}

type account struct {
	Name string `json:"name"`
}

type createAccountBody struct {
	Account account `json:"account"`
}

func decodeCreateAccountRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body createAccountBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, errBadRequest
	}

	newAccountRequest := createAccountRequest{
		Name: body.Account.Name,
	}

	return newAccountRequest, nil
}

func decodeSendPaymentRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body sendPaymentBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, errBadRequest
	}

	paymentRequest := sendPaymentRequest{
		From:   entities.Account{Name: body.Payment.From},
		To:     entities.Account{Name: body.Payment.To},
		Amount: body.Payment.Amount,
	}

	return paymentRequest, nil
}

func encodePaymentsAsJSON(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp := response.(getPaymentsResponse)

	encoder := paymentsJSONEncoder{}
	encoded, err := encoder.encode(resp.Payments)
	if err != nil {
		return errors.Wrap(err, "Can't encode payments into bytes")
	}

	_, writeErr := w.Write(encoded)
	return errors.Wrap(writeErr, "Can't write response body")
}

func MakeHandler(svc BankingService, l log.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(errorEncoder),
		kithttp.ServerErrorLogger(l),
	}

	createAccount := kithttp.NewServer(
		MakeCreateAccountEndpoint(svc),
		decodeCreateAccountRequest,
		kithttp.EncodeJSONResponse,
		opts...,
	)

	getAccounts := kithttp.NewServer(
		MakeGetAccountsEndpoint(svc),
		kithttp.NopRequestDecoder,
		kithttp.EncodeJSONResponse,
		opts...,
	)

	getPayments := kithttp.NewServer(
		MakeGetPaymentsEndpoint(svc),
		kithttp.NopRequestDecoder,
		encodePaymentsAsJSON,
		opts...,
	)

	sendPayment := kithttp.NewServer(
		MakeSendPaymentEndpoint(svc),
		decodeSendPaymentRequest,
		kithttp.EncodeJSONResponse,
		opts...,
	)

	m := mux.NewRouter()
	m.Handle("/accounts", createAccount).Methods(http.MethodPost)
	m.Handle("/accounts", getAccounts).Methods(http.MethodGet)
	m.Handle("/payments", getPayments).Methods(http.MethodGet)
	m.Handle("/payments", sendPayment).Methods(http.MethodPost)
	m.NotFoundHandler = http.HandlerFunc(notFoundEncoder)
	return m
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch errors.Cause(err) {
	case errBadRequest:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	encodeErr := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	if encodeErr != nil {
		panic(fmt.Sprintf("Can't encode error, %s. Original error: %s", encodeErr, err))
	}
}

func notFoundEncoder(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": "page not found",
	})
}
