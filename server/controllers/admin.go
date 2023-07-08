package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mmdaz/lime/license"
	"github.com/mmdaz/lime/server/middleware"
	"github.com/mmdaz/lime/server/models"
)

// MainHandler is a ...
func MainHandler(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(middleware.IdentityKey)

	if user == nil {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "🔑 Autorisation",
		})
	} else {
		customersList := models.CustomersList()
		c.HTML(http.StatusOK, "customers.html", gin.H{
			"title":     "🐱 Customers",
			"customers": customersList,
		})
	}

}

func CreateCustomer(c *gin.Context) {
	// Bind request body to struct
	request := &requestCreateCustomer{}
	err := c.BindJSON(request)
	if err != nil {
		respondJSON(c, 400, err.Error())
		return
	}

	customer := models.Customer{
		Name:   request.Name,
		Status: request.Status,
	}
	cos, err := customer.SaveCustomer()
	if err != nil {
		respondJSON(c, 500, err.Error())
		return
	}

	// create subscription
	subs := &models.Subscription{
		CustomerID: cos.ID,
		Status:     true,
	}
	_, err = subs.SaveSubscription()
	if err != nil {
		respondJSON(c, 500, err.Error())
		return
	}

	respondJSON(c, http.StatusOK, "Customer Created.")

}

// CustomerSubscrptionList is a ...
func CustomerSubscrptionList(c *gin.Context) {
	customerID := c.Param("id")
	action := c.Param("action")
	subscriptionsList := models.SubscriptionsByCustomerID(customerID)
	subscriptionID := strconv.Itoa(int((*subscriptionsList)[0].ID))
	licensesList := models.LicensesListBySubscriptionID(subscriptionID)

	switch action {
	case "/":
	case "/new":
		month := time.Hour * 24 * 31
		hwID := c.PostForm("hw_id")
		metadata := []byte(fmt.Sprintf(`{"hw_id": "%s"}`, hwID))
		_license := &license.License{
			Iss: (*subscriptionsList)[0].CustomerName,
			Cus: (*subscriptionsList)[0].CustomerID,
			Sub: (*subscriptionsList)[0].TariffID,
			Typ: "Module Name",
			Dat: metadata,
			Exp: time.Now().UTC().Add(month).Unix(),
			Iat: time.Now().UTC().Unix(),
		}
		encoded, _ := _license.Encode(license.GetPrivateKey())

		hash := md5.Sum([]byte(encoded))
		licenseHash := hex.EncodeToString(hash[:])

		models.DeactivateLicenseBySubID((*subscriptionsList)[0].ID)
		key := &models.License{
			SubscriptionID: (*subscriptionsList)[0].ID,
			License:        encoded,
			Hash:           licenseHash,
			Status:         true,
		}
		key.SaveLicense()

		c.Redirect(http.StatusFound, "/admin/subscription/"+customerID+"/")

	default:
		c.Redirect(http.StatusFound, "/admin/subscription/"+customerID+"/")
	}

	c.HTML(http.StatusOK, "subscriptions.html", gin.H{
		"title":         "🧩 Subscription and Licenses by " + (*subscriptionsList)[0].CustomerName,
		"customerID":    customerID,
		"Subscriptions": subscriptionsList,
		"Licensies":     licensesList,
	})

}


// DownloadLicense is a ...
func DownloadLicense(c *gin.Context) {
	licenseID := c.Param("id")
	license := models.License{}

	uid, err := strconv.ParseUint(licenseID, 10, 32)
	if err != nil {
		panic(err)
	}

	license.FindLicenseByID(uint32(uid))

	body := string(license.License)
	reader := strings.NewReader(body)
	contentLength := int64(len(body))
	contentType := "application/octet-stream"
	extraHeaders := map[string]string{"Content-Disposition": `attachment; filename="` + license.Hash + `"`}
	c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)
}
