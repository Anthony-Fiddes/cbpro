package coinbase

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
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
		return "", errors.New("coinbase: arguments to sign cannot be empty")
	}
	prehash := timestamp + method + requestPath
	secretKey, err := base64.StdEncoding.DecodeString(c.Secret)
	if err != nil {
		return "", fmt.Errorf("coinbase: error decoding secret %w", err)
	}
	mac := hmac.New(sha256.New, secretKey)
	_, err = io.WriteString(mac, prehash)
	if err != nil {
		return "", fmt.Errorf("coinbase: error hashing signature %w", err)
	}
	sig := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return sig, nil
}

// Request makes a request to the Coinbase Pro API using the given client.
func (c *Client) Request(method, requestPath string, query url.Values, body io.Reader) (*http.Response, error) {
	if c.Doer == nil {
		return nil, errors.New("coinbase: Client.Doer is nil")
	}

	u, err := url.Parse(apiURL)
	if err != nil {
		return nil, fmt.Errorf("coinbase: error parsing url %w", err)
	}
	u.Path = path.Join(u.Path, requestPath)
	u.RawQuery = query.Encode()

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, fmt.Errorf("coinbase: error generating request %w", err)
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
		return nil, fmt.Errorf("coinbase: error performing request %w", err)
	}
	return result, nil
}

func timestamp() string {
	now := time.Now()
	t := fmt.Sprintf("%d", now.UTC().Unix())
	return t
}
