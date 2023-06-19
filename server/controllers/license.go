package controllers

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mmdaz/lime/license"
	"github.com/mmdaz/lime/server/models"
)

// VerifyKey is a ...
// @Accept application/json
// @Produce application/json
// @Param
// @Success 200 {string} string "{"status":"200", "msg":""}"
// @Router /verify [post]
func VerifyKey(c *gin.Context) {
	modelLicense := models.License{}

	reques := &requestLicense{}
	c.BindJSON(&reques)

	licenseKey, err := base64.StdEncoding.DecodeString(reques.License)
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}

	_license, err := modelLicense.FindLicense(licenseKey)
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}
	if _license.ID == 0 {
		respondJSON(c, http.StatusNotFound, "License not found!")
		return
	}

	l, err := license.Decode([]byte(licenseKey), license.GetPublicKey())
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(l)

	respondJSON(c, http.StatusOK, "Active")
}

// CreateKey is a ...
// @Accept application/json
// @Produce application/json
// @Param
// @Success 200 {string} string "{"status":"200", "msg":""}"
// @Router /key [post]
func CreateKey(c *gin.Context) {
	month := time.Hour * 24 * 31
	modelSubscription := models.Subscription{}
	modelTariff := models.Tariff{}
	modelCustomer := models.Customer{}

	reques := &requestLicense{}
	c.BindJSON(&reques)

	_subscription, err := modelSubscription.FindSubscriptionByStripeID(reques.StripeID)
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}
	if _subscription.ID == 0 {
		respondJSON(c, http.StatusNotFound, "Customers not found!")
		return
	}

	_customer, _ := modelCustomer.FindCustomerByID(_subscription.CustomerID)

	_tariff, err := modelTariff.FindTariffByID(_subscription.TariffID)
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}
	if _tariff.ID == 0 {
		respondJSON(c, http.StatusNotFound, "Tariff not found!")
		return
	}

	limit := license.Limits{
		Servers:   _tariff.Servers,
		Companies: _tariff.Companies,
		Users:     _tariff.Users,
	}
	metadata := []byte(`{"message": "test message"}`)
	_license := &license.License{
		Iss: _customer.Name,
		Cus: _subscription.StripeID,
		Sub: _subscription.TariffID,
		Typ: _tariff.Name,
		Lim: limit,
		Dat: metadata,
		Exp: time.Now().UTC().Add(month).Unix(),
		Iat: time.Now().UTC().Unix(),
}

	encoded, err := _license.Encode(license.GetPrivateKey())
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}

	models.DeactivateLicenseBySubID(_subscription.ID)

	hash := md5.Sum([]byte(encoded))
	licenseHash := hex.EncodeToString(hash[:])

	key := &models.License{
		SubscriptionID: _subscription.ID,
		License:        encoded,
		Hash:           licenseHash,
		Status:         true,
	}

	_, err = key.SaveLicense()
	if err != nil {
		respondJSON(c, http.StatusNotFound, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, base64.StdEncoding.EncodeToString([]byte(encoded)))
}

// GetKey is a ...
// @Accept application/json
// @Produce application/json
// @Param
// @Success 200 {string} string "{"status":"200", "msg":""}"
// @Router /key/:customer_id [get]
func GetKey(c *gin.Context) {
	respondJSON(c, http.StatusOK, "GetKey")
}

// UpdateKey is a ...
// @accept application/json
// @Produce application/json
// @Param
// @Success 200 {string} string "{"status":"200", "msg":""}"
// @Router /key/:customer_id [PATCH]
func UpdateKey(c *gin.Context) {
	respondJSON(c, http.StatusOK, "UpdateKey")
}
