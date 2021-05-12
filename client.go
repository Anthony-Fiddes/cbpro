package cbpro

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"
)

const (
	apiURL               = "https://api.pro.coinbase.com/"
	keyFieldName         = "CB-ACCESS-KEY"
	passphraseFieldName  = "CB-ACCESS-PASSPHRASE"
	timestampFieldName   = "CB-ACCESS-TIMESTAMP"
	sigFieldName         = "CB-ACCESS-SIGN"
	contentTypeFieldName = "Content-Type"
	contentType          = "application/json"
)

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

// Client holds the configuration information
type Client struct {
	Secret     string
	Key        string
	Passphrase string
	Doer
}

// sign creates the string that is used to sign requests to the API.
func (c *Client) sign(timestamp, method, requestPath string) (string, error) {
	if timestamp == "" || method == "" || requestPath == "" {
		return "", errors.New("cbpro: arguments to sign cannot be empty")
	}
	prehash := timestamp + method + requestPath
	secretKey, err := base64.StdEncoding.DecodeString(c.Secret)
	if err != nil {
		return "", fmt.Errorf("cbpro: error decoding secret %w", err)
	}
	mac := hmac.New(sha256.New, secretKey)
	_, err = io.WriteString(mac, prehash)
	if err != nil {
		return "", fmt.Errorf("cbpro: error hashing signature %w", err)
	}
	sig := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return sig, nil
}

// Request makes a request to the Coinbase Pro API using the given client.
func (c *Client) Request(method, requestPath string, query url.Values, body io.Reader) (*http.Response, error) {
	if c.Doer == nil {
		return nil, errors.New("cbpro: Client.Doer is nil")
	}

	u, err := url.Parse(apiURL)
	if err != nil {
		return nil, fmt.Errorf("cbpro: error parsing url %w", err)
	}
	u.Path = path.Join(u.Path, requestPath)
	u.RawQuery = query.Encode()

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, fmt.Errorf("cbpro: error generating request %w", err)
	}
	req.Header.Add(keyFieldName, c.Key)
	sig, err := c.sign(timestamp(), method, requestPath)
	if err != nil {
		return nil, err
	}
	req.Header.Add(sigFieldName, sig)
	req.Header.Add(passphraseFieldName, c.Passphrase)
	req.Header.Add(timestampFieldName, timestamp())
	req.Header.Add(contentTypeFieldName, contentType)

	result, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cbpro: error performing request %w", err)
	}
	return result, nil
}

func unmarshalResponse(body io.ReadCloser, v interface{}) error {
	j := json.NewDecoder(body)
	err := j.Decode(v)
	if err != nil {
		return fmt.Errorf("cbpro: error decoding json %w", err)
	}
	return nil
}

// GetProducts returns a list of all of the available Coinbase Pro currency pairs.
func (c *Client) GetProducts() ([]Product, error) {
	result, err := c.Request("GET", "/products", nil, nil)
	if err != nil {
		return nil, err
	}
	pp := make([]Product, 0, approxNumProducts)
	err = unmarshalResponse(result.Body, &pp)
	if err != nil {
		return nil, err
	}
	return pp, nil
}

// GetProduct returns the information for a specific currency pair, given its ID.
func (c *Client) GetProduct(productID string) (Product, error) {
	rp := path.Join("/products", productID)
	result, err := c.Request("GET", rp, nil, nil)
	if err != nil {
		return Product{}, err
	}
	p := Product{}
	err = unmarshalResponse(result.Body, &p)
	if err != nil {
		return Product{}, err
	}
	return p, nil
}

// GetProductStats returns the 24 hour stats for a specific currency pair.
func (c *Client) GetProductStats(productID string) (Stats, error) {
	rp := path.Join("/products", productID, "stats")
	result, err := c.Request("GET", rp, nil, nil)
	if err != nil {
		return Stats{}, err
	}
	s := Stats{}
	err = unmarshalResponse(result.Body, &s)
	if err != nil {
		return Stats{}, err
	}
	return s, nil
}

func timestamp() string {
	now := time.Now()
	t := fmt.Sprintf("%d", now.UTC().Unix())
	return t
}
