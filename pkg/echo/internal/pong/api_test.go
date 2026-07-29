package pong

import (
	"echo/internal/middleware"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPing(t *testing.T) {
	// Create a new Echo instance
	e := PongEchoServer()

	// Create a new HTTP request
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()

	// Serve the HTTP request
	e.ServeHTTP(rec, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"ping":"pong"}`, rec.Body.String())
}

func startMockJWKS() *httptest.Server {
	jwks := `{
		"keys": [
			{
				"kty": "RSA",
				"kid": "test-key-id",
				"use": "sig",
				"alg": "RS256",
				"n": "u1SU1LfVLPHCozMxH2Mo4lgOEePzNm0tRgeLezV6ffAt0gunVTLw7onLRnrq0_IzW7yWR7QkrmBL7jTKEn5u-qKhbwKfBstIs-bMY2Zkp18gnTxKLxoS2tFczGkPLPgizskuemMghRniWaoLcyehkd3qqGElvW_VDL5AaWTg0nLVkjRo9z-40RQzuVaE8AkAFmxZzow3x-VJYKdjykkJ0iT9wCS0DRTXu269V264Vf_3jvredZiKRkgwlL9xNAwxXFg0x_XFw005UWVRIkdgcKWTjpBP2dPwVZ4WWC-9aGVd-Gyn1o0CLelf4rEjGoXbAAEgAqeGUxrcIlbjXfbcmw",
				"e": "AQAB"
			}
		]
	}`

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(jwks))
	}))
}

func TestGetRestricted(t *testing.T) {
	mockJWKS := startMockJWKS()
	defer mockJWKS.Close()

	middleware.JwksURL = mockJWKS.URL

	e := PongEchoServer()

	req := httptest.NewRequest(http.MethodGet, "/restricted", nil)
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6InRlc3Qta2V5LWlkIiwidHlwIjoiSldUIn0.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.Vm4Dg5bskBWMUGTYPO6Tgge2nLJ1Oa6rkN3B-s9PsFQ9zvQQ_2V1K73X70X3zW5JiXElPRgXoJu-op0UwVt34uPdfMrMeQ5O1Ja-H6TvO2JDYMMgYz1yPp36-UY73Y7t-i2RTEXMrc9_piGtOimL7lpE5N58iQFxG4GmjyGgiLZvvczYYS0EdpJrVx4brT5pFMQ-ltPxLByw6z8jqpwzqNGssHlmzObtsKYysHOaYYfHlTDff2PeGgu6Fb5ZkNRhkQaEjCYXs1eoVVYu2w8v6FBe8sgzaQkkfiOhuiSxu0vGRv3breSe0J2xUbM1RJUjibZSnLuuFruV6wKsJWivDA")
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"ping":"Welcome John Doe!"}`, rec.Body.String())
}
