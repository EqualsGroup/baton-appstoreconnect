package client

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

const (
	BaseURL = "https://api.appstoreconnect.apple.com"

	// JWT tokens are valid for up to 20 minutes. We refresh at 15 minutes.
	tokenLifetime  = 20 * time.Minute
	tokenRefreshAt = 15 * time.Minute
)

// Client is an HTTP client for the App Store Connect API.
type Client struct {
	httpClient *http.Client
	issuerID   string
	keyID      string
	privateKey *ecdsa.PrivateKey

	mu          sync.Mutex
	token       string
	tokenExpiry time.Time
}

// New creates a new App Store Connect API client.
func New(_ context.Context, issuerID, keyID, privateKeyPath string) (*Client, error) {
	keyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file %s: %w", privateKeyPath, err)
	}

	privateKey, err := parseP8PrivateKey(keyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		issuerID:   issuerID,
		keyID:      keyID,
		privateKey: privateKey,
	}, nil
}

// parseP8PrivateKey parses an ECDSA private key from a .p8 (PEM-encoded PKCS#8) file.
func parseP8PrivateKey(data []byte) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block from private key")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PKCS8 private key: %w", err)
	}

	ecKey, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("private key is not an ECDSA key")
	}

	return ecKey, nil
}

// generateToken creates a new JWT for the App Store Connect API.
func (c *Client) generateToken() (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Issuer:    c.issuerID,
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(tokenLifetime)),
		Audience:  jwt.ClaimStrings{"appstoreconnect-v1"},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	token.Header["kid"] = c.keyID

	signedToken, err := token.SignedString(c.privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}

	return signedToken, nil
}

// getToken returns a valid JWT token, refreshing if necessary.
func (c *Client) getToken() (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.token != "" && time.Now().Before(c.tokenExpiry) {
		return c.token, nil
	}

	token, err := c.generateToken()
	if err != nil {
		return "", err
	}

	c.token = token
	c.tokenExpiry = time.Now().Add(tokenRefreshAt)

	return c.token, nil
}

// Do executes an HTTP request with JWT authentication and handles rate limiting.
func (c *Client) Do(ctx context.Context, req *http.Request) (*http.Response, error) {
	l := ctxzap.Extract(ctx)

	token, err := c.getToken()
	if err != nil {
		return nil, fmt.Errorf("failed to get auth token: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Handle rate limiting (429 Too Many Requests).
	if resp.StatusCode == http.StatusTooManyRequests {
		resp.Body.Close()

		retryAfter := resp.Header.Get("Retry-After")
		waitDuration := 60 * time.Second // default wait
		if retryAfter != "" {
			if seconds, parseErr := strconv.Atoi(retryAfter); parseErr == nil && seconds > 0 {
				waitDuration = time.Duration(seconds) * time.Second
			}
		}

		l.Warn("rate limited by App Store Connect API, retrying",
			zap.Duration("retry_after", waitDuration),
		)

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(waitDuration):
		}

		// Retry the request.
		return c.Do(ctx, req)
	}

	return resp, nil
}
