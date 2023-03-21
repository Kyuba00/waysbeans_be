package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	dto "waysbeans_be/dto/result"
	transactiondto "waysbeans_be/dto/transaction"
	"waysbeans_be/models"
	"waysbeans_be/repositories"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
	"gopkg.in/gomail.v2"
)

var c = coreapi.Client{
	ServerKey: os.Getenv("SERVER_KEY"),
	ClientKey: os.Getenv("CLIENT_KEY"),
}

type handlerTransaction struct {
	TransactionRepository repositories.TransactionRepository
}

func HandlerTransaction(TransactionRepository repositories.TransactionRepository) *handlerTransaction {
	return &handlerTransaction{TransactionRepository}
}

func (h *handlerTransaction) FindTransactions(c echo.Context) error {
	// GET USER ROLE FROM JWT TOKEN
	userInfo := c.Get("userInfo").(jwt.MapClaims)
	userRole := userInfo["role"]

	// CHECK ROLE ADMIN
	if userRole != "admin" {
		return c.JSON(http.StatusUnauthorized, dto.ErrorResult{Code: http.StatusUnauthorized, Message: "You're not admin"})
	}

	// RUN REPOSITORY FIND TRANSACTIONS
	transaction, err := h.TransactionRepository.FindTransactions()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	// WRITE RESPONSE
	return c.JSON(http.StatusOK, dto.SuccessResult{Code: "success", Data: transaction})
}

func (h *handlerTransaction) GetUserTransactionByUserID(c echo.Context) error {
	// GET USER ID FROM JWT TOKEN
	userInfo := c.Get("userInfo").(jwt.MapClaims)
	userID := int(userInfo["id"].(float64))

	// RUN REPOSITORY GET TRANSACTION BY USER ID
	transactions, err := h.TransactionRepository.GetUserTransactionByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	// WRITE RESPONSE
	return c.JSON(http.StatusOK, dto.SuccessResult{Code: "success", Data: transactions})
}

func (h *handlerTransaction) UpdateTransaction(c echo.Context) error {
	// GET USER ID FROM JWT TOKEN
	userInfo := c.Get("userInfo").(jwt.MapClaims)
	userID := int(userInfo["id"].(float64))

	// GET REQUEST AND DECODING JSON
	request := new(transactiondto.TransactionRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	// RUN REPOSITORY GET TRANSACTION BY USER ID
	transaction, err := h.TransactionRepository.GetTransactionByUserID(userID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: "Cart Failed!"})
	}

	// CHECK UPDATE VALUE
	if request.Name != "" {
		transaction.Name = request.Name
	}

	if request.Email != "" {
		transaction.Email = request.Email
	}

	if request.Phone != "" {
		transaction.Phone = request.Phone
	}

	if request.Address != "" {
		transaction.Address = request.Address
	}

	transaction.Status = "pending"
	transaction.Total = request.Total
	transaction.UpdateAt = time.Now()

	// RUN REPOSITORY UPDATE TRANSACTION
	_, err = h.TransactionRepository.UpdateTransaction(transaction)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	// SETUP FOR MIDTRANS
	DataSnap, _ := h.TransactionRepository.GetTransactionNotification(int(transaction.ID))

	var s = snap.Client{}
	s.New(os.Getenv("SERVER_KEY"), midtrans.Sandbox)

	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  strconv.Itoa(int(DataSnap.ID)),
			GrossAmt: int64(DataSnap.Total),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: DataSnap.User.Name,
			Email: DataSnap.User.Email,
		},
	}

	// RUN MIDTRANS SNAP
	snapResp, _ := s.CreateTransaction(req)

	// WRITE RESPONSE
	return c.JSON(http.StatusOK, dto.SuccessResult{Code: "Success", Data: snapResp})
}

func (h *handlerTransaction) Notification(c echo.Context) error {
	var notificationPayload map[string]interface{}

	err := c.Bind(&notificationPayload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResult{Code: http.StatusBadRequest, Message: err.Error()})
	}

	transactionStatus := notificationPayload["transaction_status"].(string)
	fraudStatus := notificationPayload["fraud_status"].(string)
	orderID := notificationPayload["order_id"].(string)

	transaction, _ := h.TransactionRepository.GetTransactionMidtrans(orderID)

	if transactionStatus == "capture" {
		if fraudStatus == "challenge" {
			// TODO set transaction status on your database to 'challenge'
			// e.g: 'Payment status challenged. Please take action on your Merchant Administration Portal
			h.TransactionRepository.UpdateTransactionMidtrans("pending", transaction.ID)
		} else if fraudStatus == "accept" {
			// TODO set transaction status on your database to 'success'
			SendMail("success", transaction)
			h.TransactionRepository.UpdateTransactionMidtrans("success", transaction.ID)
		}
	} else if transactionStatus == "settlement" {
		// TODO set transaction status on your databaase to 'success'
		SendMail("success", transaction)
		h.TransactionRepository.UpdateTransactionMidtrans("success", transaction.ID)
	} else if transactionStatus == "deny" {
		// TODO you can ignore 'deny', because most of the time it allows payment retries
		// and later can become success
		SendMail("failed", transaction)
		h.TransactionRepository.UpdateTransactionMidtrans("failed", transaction.ID)
	} else if transactionStatus == "cancel" || transactionStatus == "expire" {
		// TODO set transaction status on your databaase to 'failure'
		SendMail("failed", transaction) // Call SendMail function ...
		h.TransactionRepository.UpdateTransactionMidtrans("failed", transaction.ID)
	} else if transactionStatus == "pending" {
		// TODO set transaction status on your databaase to 'pending' / waiting payment
		h.TransactionRepository.UpdateTransactionMidtrans("pending", transaction.ID)
	}

	return c.JSON(http.StatusOK, dto.SuccessResult{Code: "Success", Data: notificationPayload})
}

func SendMail(status string, transaction models.Transaction) {

	if status != transaction.Status && (status == "success") {
		var CONFIG_SMTP_HOST = "smtp.gmail.com"
		var CONFIG_SMTP_PORT = 587
		var CONFIG_SENDER_NAME = "WaysBeans <akanime1@gmail.com>"
		var CONFIG_AUTH_EMAIL = os.Getenv("EMAIL_SYSTEM")
		var CONFIG_AUTH_PASSWORD = os.Getenv("PASSWORD_SYSTEM")

		var productName = transaction.User.Name
		var price = strconv.Itoa(int(transaction.Total))

		mailer := gomail.NewMessage()
		mailer.SetHeader("From", CONFIG_SENDER_NAME)
		mailer.SetHeader("To", transaction.User.Email)
		mailer.SetHeader("Subject", "Transaction Status")
		mailer.SetBody("text/html", fmt.Sprintf(`<!DOCTYPE html>
	  <html lang="en">
		<head>
		<meta charset="UTF-8" />
		<meta http-equiv="X-UA-Compatible" content="IE=edge" />
		<meta name="viewport" content="width=device-width, initial-scale=1.0" />
		<title>Document</title>
		<style>
		  h1 {
		  color: brown;
		  }
		</style>
		</head>
		<body>
		<h2>Product payment :</h2>
		<ul style="list-style-type:none;">
		  <li>Name : %s</li>
		  <li>Total payment: Rp.%s</li>
		  <li>Status : <b>%s</b></li>
		</ul>
		</body>
	  </html>`, productName, price, status))

		dialer := gomail.NewDialer(
			CONFIG_SMTP_HOST,
			CONFIG_SMTP_PORT,
			CONFIG_AUTH_EMAIL,
			CONFIG_AUTH_PASSWORD,
		)

		err := dialer.DialAndSend(mailer)
		if err != nil {
			log.Fatal(err.Error())
		}

		log.Println("Mail sent! to " + transaction.User.Email)

	}
}
