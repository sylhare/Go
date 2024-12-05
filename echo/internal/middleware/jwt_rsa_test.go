package middleware

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"math/big"
	"strings"
	"testing"
)

func TestRSA(t *testing.T) {
	pemData := `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAu1SU1LfVLPHCozMxH2Mo
4lgOEePzNm0tRgeLezV6ffAt0gunVTLw7onLRnrq0/IzW7yWR7QkrmBL7jTKEn5u
+qKhbwKfBstIs+bMY2Zkp18gnTxKLxoS2tFczGkPLPgizskuemMghRniWaoLcyeh
kd3qqGElvW/VDL5AaWTg0nLVkjRo9z+40RQzuVaE8AkAFmxZzow3x+VJYKdjykkJ
0iT9wCS0DRTXu269V264Vf/3jvredZiKRkgwlL9xNAwxXFg0x/XFw005UWVRIkdg
cKWTjpBP2dPwVZ4WWC+9aGVd+Gyn1o0CLelf4rEjGoXbAAEgAqeGUxrcIlbjXfbc
mwIDAQAB
-----END PUBLIC KEY-----`

	block, _ := pem.Decode([]byte(pemData))
	if block == nil || block.Type != "PUBLIC KEY" {
		fmt.Println("failed to decode PEM block containing public key")
		return
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		fmt.Println("failed to parse DER encoded public key:", err)
		return
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		fmt.Println("not an RSA public key")
		return
	}

	n := rsaPub.N
	e := rsaPub.E

	fmt.Printf("Modulus (n): %s\n", n.String())
	fmt.Printf("Exponent (e): %d\n", e)

	// Encode the modulus as Base64 URL without padding
	nBytes := n.Bytes()
	nBase64 := base64.URLEncoding.EncodeToString(nBytes)
	nBase64 = strings.TrimRight(nBase64, "=")
	fmt.Println("Encoded n:", nBase64)
	fmt.Println("Encoded e:", intToBase64(e))

	// Parse the token
	var jwtToken = "eyJhbGciOiJSUzI1NiIsImtpZCI6InRlc3Qta2V5LWlkIiwidHlwIjoiSldUIn0.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.Vm4Dg5bskBWMUGTYPO6Tgge2nLJ1Oa6rkN3B-s9PsFQ9zvQQ_2V1K73X70X3zW5JiXElPRgXoJu-op0UwVt34uPdfMrMeQ5O1Ja-H6TvO2JDYMMgYz1yPp36-UY73Y7t-i2RTEXMrc9_piGtOimL7lpE5N58iQFxG4GmjyGgiLZvvczYYS0EdpJrVx4brT5pFMQ-ltPxLByw6z8jqpwzqNGssHlmzObtsKYysHOaYYfHlTDff2PeGgu6Fb5ZkNRhkQaEjCYXs1eoVVYu2w8v6FBe8sgzaQkkfiOhuiSxu0vGRv3breSe0J2xUbM1RJUjibZSnLuuFruV6wKsJWivDA"
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		return rsaPub, nil
	})
	if err != nil {
		fmt.Println("Error parsing token:", err)
	}

	// Print the token claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Println("Token Claims:")
		for key, value := range claims {
			fmt.Printf("%s: %v\n", key, value)
		}
	} else {
		fmt.Println("Invalid token claims")
	}
	fmt.Println("Token is valid:", token.Valid)

	// Encoded values
	nBase64 = "u1SU1LfVLPHCozMxH2Mo4lgOEePzNm0tRgeLezV6ffAt0gunVTLw7onLRnrq0_IzW7yWR7QkrmBL7jTKEn5u-qKhbwKfBstIs-bMY2Zkp18gnTxKLxoS2tFczGkPLPgizskuemMghRniWaoLcyehkd3qqGElvW_VDL5AaWTg0nLVkjRo9z-40RQzuVaE8AkAFmxZzow3x-VJYKdjykkJ0iT9wCS0DRTXu269V264Vf_3jvredZiKRkgwlL9xNAwxXFg0x_XFw005UWVRIkdgcKWTjpBP2dPwVZ4WWC-9aGVd-Gyn1o0CLelf4rEjGoXbAAEgAqeGUxrcIlbjXfbcmw"
	eBase64 := "AQAB"

	// Decode n
	nBytes, err = base64.RawURLEncoding.DecodeString(nBase64)
	if err != nil {
		fmt.Println("Error decoding n:", err)
		return
	}
	n = new(big.Int).SetBytes(nBytes)

	// Decode e
	eBytes, err := base64.RawURLEncoding.DecodeString(eBase64)
	if err != nil {
		fmt.Println("Error decoding e:", err)
		return
	}
	e = int(new(big.Int).SetBytes(eBytes).Uint64())

	// Construct the RSA public key
	rsaPub = &rsa.PublicKey{
		N: n,
		E: e,
	}

	// Print the RSA public key
	fmt.Printf("RSA Public Key: %+v\n", rsaPub)

	// Optionally, encode the public key to PEM format
	pubASN1, err := x509.MarshalPKIXPublicKey(rsaPub)
	if err != nil {
		fmt.Println("Error marshaling public key:", err)
		return
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	})
	fmt.Printf("PEM Encoded Public Key:\n%s\n", pubPEM)
}

func intToBase64(n int) string {
	// Convert the integer to a byte slice
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, uint32(n))

	// Remove leading zero bytes
	var start int
	for start = 0; start < len(bytes); start++ {
		if bytes[start] != 0 {
			break
		}
	}
	bytes = bytes[start:]

	// Encode the byte slice to base64
	encoded := base64.StdEncoding.EncodeToString(bytes)
	return encoded
}
