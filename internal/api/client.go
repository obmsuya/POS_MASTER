package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

const productionBaseURL = "https://backend.wapangaji.com/api/v1"

var BaseURL = resolveBaseURL()

func resolveBaseURL() string {
	overrideURL := os.Getenv("FALTASI_API_BASE_URL")
	if overrideURL != "" {
		return overrideURL
	}
	return productionBaseURL
}

type Client struct {
	http         *http.Client
	AccessToken  string
	RefreshToken string
}

func New() *Client {
	return &Client{http: &http.Client{Timeout: 20 * time.Second}}
}

type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string { return e.Message }

type User struct {
	ID          string `json:"id"`
	PhoneNumber string `json:"phone_number"`
	FullName    string `json:"full_name"`
	IsSuperuser bool   `json:"is_superuser"`
	UserType    string `json:"user_type"`
}

type loginResponse struct {
	User   User `json:"user"`
	Tokens struct {
		Access  string `json:"access"`
		Refresh string `json:"refresh"`
	} `json:"tokens"`
}

func (c *Client) Login(phone, password string) (*User, error) {
	body, err := json.Marshal(map[string]string{"phone_number": phone, "password": password})
	if err != nil {
		return nil, err
	}

	var parsed loginResponse
	if err := c.do(http.MethodPost, "/auth/login/", body, false, &parsed); err != nil {
		return nil, err
	}

	c.AccessToken = parsed.Tokens.Access
	c.RefreshToken = parsed.Tokens.Refresh
	return &parsed.User, nil
}

func (c *Client) RestoreSession(access, refresh string) {
	c.AccessToken = access
	c.RefreshToken = refresh
}

func (c *Client) RefreshAccessToken() error {
	body, err := json.Marshal(map[string]string{"refresh": c.RefreshToken})
	if err != nil {
		return err
	}
	var parsed struct {
		Access string `json:"access"`
	}
	if err := c.do(http.MethodPost, "/auth/token/refresh/", body, false, &parsed); err != nil {
		return err
	}
	c.AccessToken = parsed.Access
	return nil
}

type LicenseInfo struct {
	LicenseKey  string `json:"license_key"`
	Phone       string `json:"phone"`
	Package     string `json:"package"`
	ExpiresAt   string `json:"expires_at"`
	DaysGranted int    `json:"days_granted"`
	MaxDevices  int    `json:"max_devices"`
	DeviceCount int    `json:"device_count"`
}

type lookupResponse struct {
	Found bool        `json:"found"`
	Data  LicenseInfo `json:"data"`
}

func (c *Client) Lookup(hardwareID, phone string) (*LicenseInfo, error) {
	query := ""
	if hardwareID != "" {
		query = "?hardware_id=" + urlEscape(hardwareID)
	} else {
		query = "?phone=" + urlEscape(phone)
	}

	var parsed lookupResponse
	if err := c.do(http.MethodGet, "/payments/balce/admin/lookup/"+query, nil, true, &parsed); err != nil {
		return nil, err
	}
	if !parsed.Found {
		return nil, nil
	}
	return &parsed.Data, nil
}

type Package struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Price       string `json:"price"`
	DaysGranted int    `json:"days_granted"`
	MaxDevices  int    `json:"max_devices"`
}

type packagesResponse struct {
	Data []Package `json:"data"`
}

func (c *Client) ListPackages() ([]Package, error) {
	var parsed packagesResponse
	if err := c.do(http.MethodGet, "/payments/balce/packages/", nil, false, &parsed); err != nil {
		return nil, err
	}
	return parsed.Data, nil
}

type CreatePackageInput struct {
	Name        string
	Price       float64
	DaysGranted int
	MaxDevices  int
}

func (c *Client) CreatePackage(input CreatePackageInput) (*Package, error) {
	body, err := json.Marshal(map[string]interface{}{
		"name":         input.Name,
		"price":        input.Price,
		"days_granted": input.DaysGranted,
		"max_devices":  input.MaxDevices,
	})
	if err != nil {
		return nil, err
	}
	var parsed struct {
		Data Package `json:"data"`
	}
	if err := c.do(http.MethodPost, "/payments/balce/packages/", body, true, &parsed); err != nil {
		return nil, err
	}
	return &parsed.Data, nil
}

func (c *Client) UpdatePackage(id int, input CreatePackageInput) (*Package, error) {
	body, err := json.Marshal(map[string]interface{}{
		"name":         input.Name,
		"price":        input.Price,
		"days_granted": input.DaysGranted,
		"max_devices":  input.MaxDevices,
	})
	if err != nil {
		return nil, err
	}
	var parsed struct {
		Data Package `json:"data"`
	}
	path := "/payments/balce/packages/" + strconv.Itoa(id) + "/"
	if err := c.do(http.MethodPatch, path, body, true, &parsed); err != nil {
		return nil, err
	}
	return &parsed.Data, nil
}

type RecordPaymentInput struct {
	HardwareID string
	Phone      string
	PackageID  int
	Amount     float64
	Method     string
	Reference  string
}

type recordPaymentResponse struct {
	Data LicenseInfo `json:"data"`
}

func (c *Client) RecordPayment(input RecordPaymentInput) (*LicenseInfo, error) {
	body, err := json.Marshal(map[string]interface{}{
		"hardware_id": input.HardwareID,
		"phone":       input.Phone,
		"package_id":  input.PackageID,
		"amount":      input.Amount,
		"method":      input.Method,
		"reference":   input.Reference,
	})
	if err != nil {
		return nil, err
	}
	var parsed recordPaymentResponse
	if err := c.do(http.MethodPost, "/payments/balce/admin/record-payment/", body, true, &parsed); err != nil {
		return nil, err
	}
	return &parsed.Data, nil
}

type PaymentRecord struct {
	Phone       string `json:"phone"`
	Package     string `json:"package"`
	Amount      string `json:"amount"`
	Method      string `json:"method"`
	Reference   string `json:"reference"`
	ProcessedBy string `json:"processed_by"`
	CreatedAt   string `json:"created_at"`
}

type historyResponse struct {
	Data []PaymentRecord `json:"data"`
}

func (c *Client) PaymentHistory(hardwareID string) ([]PaymentRecord, error) {
	path := "/payments/balce/admin/payments/"
	if hardwareID != "" {
		path += "?hardware_id=" + urlEscape(hardwareID)
	}
	var parsed historyResponse
	if err := c.do(http.MethodGet, path, nil, true, &parsed); err != nil {
		return nil, err
	}
	return parsed.Data, nil
}

func (c *Client) do(method, path string, body []byte, authorized bool, out interface{}) error {
	response, err := c.request(method, path, body, authorized)
	if err != nil {
		return err
	}
	defer func() { response.Body.Close() }()

	if response.StatusCode == http.StatusUnauthorized && authorized && c.RefreshToken != "" {
		if refreshErr := c.RefreshAccessToken(); refreshErr == nil {
			response.Body.Close()
			response, err = c.request(method, path, body, authorized)
			if err != nil {
				return err
			}
		}
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		var errBody struct {
			Error string `json:"error"`
		}
		message := fmt.Sprintf("server returned %d", response.StatusCode)
		if json.Unmarshal(responseBody, &errBody) == nil && errBody.Error != "" {
			message = errBody.Error
		}
		return &APIError{StatusCode: response.StatusCode, Message: message}
	}

	if out == nil || len(responseBody) == 0 {
		return nil
	}
	return json.Unmarshal(responseBody, out)
}

func (c *Client) request(method, path string, body []byte, authorized bool) (*http.Response, error) {
	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}
	req, err := http.NewRequest(method, BaseURL+path, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if authorized && c.AccessToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	}
	return c.http.Do(req)
}

func urlEscape(value string) string {
	return url.QueryEscape(value)
}
