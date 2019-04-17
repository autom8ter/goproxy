# GoProxy
--
    import "github.com/autom8ter/goproxy"

- [goproxy](#goproxy)
  * [Usage](#usage)
      - [type GoProxy](#type-goproxy)
      - [func  NewGoProxy](#func--newgoproxy)
      - [func (*GoProxy) ServeHTTP](#func---goproxy--servehttp)
  * [Example (Stripe)](#example)
    + [Code:](#code-)
    + [Deploy:](#deploy-)
    + [Output: (customer_ID)](#output---customer-id-)

## Usage

#### type GoProxy

```go
type GoProxy struct {
}
```

GoProxy is an API Gateway/Reverse Proxy and http.ServeMux/http.Handler

#### func  NewGoProxy

```go
func NewGoProxy(config *config.Config) *GoProxy
```
NewGoProxy registers a new reverseproxy handler for each provided config with
the specified path prefix

#### func (*GoProxy) ServeHTTP

```go
func (g *GoProxy) ServeHTTP(w http.ResponseWriter, r *http.Request)
```

## Example

### Code:

#### Proxy:
```text

var BaseURL = "https://api.stripe.com/v1/customers"

var proxy = goproxy.NewGoProxy(&config.Config{
	TargetUrl:           BaseURL,
	Username:            os.Getenv("ACCOUNT"),
	//
	Password:            os.Getenv("KEY"),
	// Callback url to handle response from the Stripe customers rest API
	ResponseCallbackURL: os.Getenv("CALLBACK"),
})

func CreateCustomer(w http.ResponseWriter, r *http.Request) {
	proxy.ServeHTTP(w, r)
}

```
#### Callback:
```go

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
	//Send back just the customer id to the original requester
	fmt.Fprintf(w, "%s", util.Handle.MarshalJSON(resp.ID))
	return
}

//This is a go representation of the stripe customers response
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

```
### Deploy:

    callback:
    	gcloud functions deploy CustomerCallback --runtime go111 --trigger-http
    
    proxy:
    	gcloud functions deploy CreateCustomer --set-env-vars ACCOUNT=XXXXXX,KEY=XXXXXX--runtime go111 --trigger-http


### Output: (original client)
    cust_XXXXXXXXX