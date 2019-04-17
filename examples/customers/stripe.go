package customers

import (
	"encoding/json"
	"fmt"
	"github.com/autom8ter/goproxy"
	"github.com/autom8ter/goproxy/config"
	"github.com/autom8ter/goproxy/util"
	"io/ioutil"
	"net/http"
	"os"
)

var BaseURL = "https://api.stripe.com/v1/customers"

var proxy = goproxy.NewGoProxy(&config.Config{
	TargetUrl:           BaseURL,
	Headers:             nil,

	FormValues:          nil,
	FlushInterval:       0,
	ResponseCallbackURL: os.Getenv("CALLBACK"),
})

func CreateCustomer(w http.ResponseWriter, r *http.Request) {
	proxy.ServeHTTP(w, r)
}

func CustomerCallback(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	bits, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := &Response{}
	err = json.Unmarshal(bits, resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "%s", util.Handle.MarshalJSON(resp.ID))
	return
}

type Response struct {
	ID              string      `json:"id"`
	Object          string      `json:"object"`
	AccountBalance  int         `json:"account_balance"`
	Created         int         `json:"created"`
	Currency        string      `json:"currency"`
	DefaultSource   interface{} `json:"default_source"`
	Delinquent      bool        `json:"delinquent"`
	Description     interface{} `json:"description"`
	Discount        interface{} `json:"discount"`
	Email           interface{} `json:"email"`
	InvoicePrefix   string      `json:"invoice_prefix"`
	InvoiceSettings struct {
		CustomFields         interface{} `json:"custom_fields"`
		DefaultPaymentMethod interface{} `json:"default_payment_method"`
		Footer               interface{} `json:"footer"`
	} `json:"invoice_settings"`
	Livemode bool `json:"livemode"`
	Metadata struct {
	} `json:"metadata"`
	Shipping interface{} `json:"shipping"`
	Sources  struct {
		Object     string        `json:"object"`
		Data       []interface{} `json:"data"`
		HasMore    bool          `json:"has_more"`
		TotalCount int           `json:"total_count"`
		URL        string        `json:"url"`
	} `json:"sources"`
	Subscriptions struct {
		Object     string        `json:"object"`
		Data       []interface{} `json:"data"`
		HasMore    bool          `json:"has_more"`
		TotalCount int           `json:"total_count"`
		URL        string        `json:"url"`
	} `json:"subscriptions"`
	TaxInfo             interface{} `json:"tax_info"`
	TaxInfoVerification interface{} `json:"tax_info_verification"`
}
