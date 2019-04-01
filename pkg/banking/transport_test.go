package banking_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/twonegatives/coinsph_challenge/pkg/banking"
	"github.com/twonegatives/coinsph_challenge/pkg/entities"
	"github.com/twonegatives/coinsph_challenge/pkg/mocks"
)

var (
	ErrSvc = errors.New("service error")
)

type dependencies struct {
	TestServer *httptest.Server
	Service    *mocks.MockBankingService
}

func setupServer(t *testing.T) (dependencies, func()) {
	mockCtrl := gomock.NewController(t)

	svc := mocks.NewMockBankingService(mockCtrl)

	router := banking.MakeHandler(svc, mocks.TestLogger{T: t})
	srv := httptest.NewServer(router)
	srv.Client()
	return dependencies{TestServer: srv, Service: svc}, func() {
		mockCtrl.Finish()
		srv.Close()
	}
}

type createAccountResponse struct {
	Account entities.Account `json:"account"`
}

type getAccountsListResponse struct {
	Accounts []entities.Account `json:"accounts"`
}

type payment struct {
	Account     string             `json:"account"`
	Amount      decimal.Decimal    `json:"amount"`
	Direction   entities.Direction `json:"direction"`
	Currency    entities.Currency  `json:"currency"`
	ToAccount   *string            `json:"to_account"`
	FromAccount *string            `json:"from_account"`
}

type getPaymentsListResponse struct {
	Payments []payment `json:"payments"`
}

func TestCreateAccountRoute(t *testing.T) {
	t.Run("renders new account", func(t *testing.T) {
		dep, cleanUp := setupServer(t)
		client := dep.TestServer.Client()
		defer cleanUp()

		name := "barry"
		expectedBody := createAccountResponse{
			Account: entities.Account{
				Name:     name,
				Balance:  decimal.New(19, 0),
				Currency: entities.USD,
			},
		}

		dep.Service.EXPECT().CreateAccount(gomock.Any(), name).Return(expectedBody.Account, nil)

		requestBody := `{"account": {"name": "barry"}}`
		resp, err := client.Post(dep.TestServer.URL+"/accounts", "application/json", strings.NewReader(requestBody))
		defer resp.Body.Close()
		require.NoError(t, err)

		var actualBody createAccountResponse
		assert.NoError(t, json.NewDecoder(resp.Body).Decode(&actualBody))

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))
		assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))
		assert.Equal(t, expectedBody, actualBody)
	})

	t.Run("returns 500 on server error", func(t *testing.T) {
		dep, cleanUp := setupServer(t)
		client := dep.TestServer.Client()
		defer cleanUp()

		name := "barry"
		dep.Service.EXPECT().CreateAccount(gomock.Any(), name).Return(entities.Account{}, ErrSvc)

		requestBody := `{"account": {"name": "barry"}}`
		resp, err := client.Post(dep.TestServer.URL+"/accounts", "application/json", strings.NewReader(requestBody))
		defer resp.Body.Close()
		require.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))
	})
}

func TestGetAccountsListRoute(t *testing.T) {
	t.Run("renders accounts list", func(t *testing.T) {
		dep, cleanUp := setupServer(t)
		client := dep.TestServer.Client()
		defer cleanUp()

		expectedBody := getAccountsListResponse{
			Accounts: []entities.Account{
				{
					Name:     "ben",
					Balance:  decimal.New(19, 0),
					Currency: entities.USD,
				},
			},
		}

		dep.Service.EXPECT().GetAccountsList(gomock.Any()).Return(expectedBody.Accounts, nil)

		resp, err := client.Get(dep.TestServer.URL + "/accounts")
		defer resp.Body.Close()
		require.NoError(t, err)

		var actualBody getAccountsListResponse
		assert.NoError(t, json.NewDecoder(resp.Body).Decode(&actualBody))

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))
		assert.Equal(t, expectedBody, actualBody)
	})

	t.Run("returns 500 on server error", func(t *testing.T) {
		dep, cleanUp := setupServer(t)
		client := dep.TestServer.Client()
		defer cleanUp()

		dep.Service.EXPECT().GetAccountsList(gomock.Any()).Return(nil, ErrSvc)

		resp, err := client.Get(dep.TestServer.URL + "/accounts")
		defer resp.Body.Close()
		require.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))
	})
}

func TestGetPaymentsListRoute(t *testing.T) {
	t.Run("renders payments list", func(t *testing.T) {
		dep, cleanUp := setupServer(t)
		client := dep.TestServer.Client()
		defer cleanUp()

		mark := entities.Account{Name: "mark", ID: 143}
		vlad := entities.Account{Name: "vlad", ID: 250}

		expectedBody := getPaymentsListResponse{
			Payments: []payment{
				{
					Account:   mark.Name,
					Direction: entities.Outgoing,
					Currency:  entities.USD,
					ToAccount: &vlad.Name,
					Amount:    decimal.New(173, 0),
				},
			},
		}

		svcResponse := []entities.Payment{
			{
				Account:      mark,
				Counterparty: vlad,
				Transaction:  entities.Transaction{},
				Direction:    entities.Outgoing,
				Amount:       decimal.New(173, 0),
				Currency:     entities.USD,
			},
		}

		dep.Service.EXPECT().GetPaymentsList(gomock.Any()).Return(svcResponse, nil)

		resp, err := client.Get(dep.TestServer.URL + "/payments")
		defer resp.Body.Close()
		require.NoError(t, err)

		var actualBody getPaymentsListResponse
		assert.NoError(t, json.NewDecoder(resp.Body).Decode(&actualBody))

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))
		assert.Equal(t, expectedBody, actualBody)
	})

	t.Run("returns 500 on server error", func(t *testing.T) {
		dep, cleanUp := setupServer(t)
		client := dep.TestServer.Client()
		defer cleanUp()

		dep.Service.EXPECT().GetPaymentsList(gomock.Any()).Return(nil, ErrSvc)

		resp, err := client.Get(dep.TestServer.URL + "/payments")
		defer resp.Body.Close()
		require.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))
	})
}

func TestSendPaymentRoute(t *testing.T) {
	t.Run("returns blank object", func(t *testing.T) {
		dep, cleanUp := setupServer(t)
		client := dep.TestServer.Client()
		defer cleanUp()

		from := entities.Account{Name: "barry"}
		to := entities.Account{Name: "wicky"}
		amount := decimal.New(1426, -2)
		dep.Service.EXPECT().SendPayment(gomock.Any(), from, to, amount).Return(nil)

		requestBody := `{"payment": {"from": "barry", "to": "wicky", "amount": 14.26}}`
		resp, err := client.Post(dep.TestServer.URL+"/payments", "application/json", strings.NewReader(requestBody))
		defer resp.Body.Close()
		require.NoError(t, err)

		var actualBody map[string]interface{}
		assert.NoError(t, json.NewDecoder(resp.Body).Decode(&actualBody))

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))
		assert.Equal(t, map[string]interface{}{}, actualBody)
	})

	t.Run("returns 500 on server error", func(t *testing.T) {
		dep, cleanUp := setupServer(t)
		client := dep.TestServer.Client()
		defer cleanUp()

		from := entities.Account{Name: "barry"}
		to := entities.Account{Name: "wicky"}
		amount := decimal.New(1426, -2)
		dep.Service.EXPECT().SendPayment(gomock.Any(), from, to, amount).Return(ErrSvc)

		requestBody := `{"payment": {"from": "barry", "to": "wicky", "amount": 14.26}}`
		resp, err := client.Post(dep.TestServer.URL+"/payments", "application/json", strings.NewReader(requestBody))
		defer resp.Body.Close()
		require.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))
	})
}
