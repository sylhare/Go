package middleware

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/pem"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log"
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

	privateKeyPEM := `-----BEGIN PRIVATE KEY-----
MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQC7VJTUt9Us8cKj
MzEfYyjiWA4R4/M2bS1GB4t7NXp98C3SC6dVMvDuictGeurT8jNbvJZHtCSuYEvu
NMoSfm76oqFvAp8Gy0iz5sxjZmSnXyCdPEovGhLa0VzMaQ8s+CLOyS56YyCFGeJZ
qgtzJ6GR3eqoYSW9b9UMvkBpZODSctWSNGj3P7jRFDO5VoTwCQAWbFnOjDfH5Ulg
p2PKSQnSJP3AJLQNFNe7br1XbrhV//eO+t51mIpGSDCUv3E0DDFcWDTH9cXDTTlR
ZVEiR2BwpZOOkE/Z0/BVnhZYL71oZV34bKfWjQIt6V/isSMahdsAASACp4ZTGtwi
VuNd9tybAgMBAAECggEBAKTmjaS6tkK8BlPXClTQ2vpz/N6uxDeS35mXpqasqskV
laAidgg/sWqpjXDbXr93otIMLlWsM+X0CqMDgSXKejLS2jx4GDjI1ZTXg++0AMJ8
sJ74pWzVDOfmCEQ/7wXs3+cbnXhKriO8Z036q92Qc1+N87SI38nkGa0ABH9CN83H
mQqt4fB7UdHzuIRe/me2PGhIq5ZBzj6h3BpoPGzEP+x3l9YmK8t/1cN0pqI+dQwY
dgfGjackLu/2qH80MCF7IyQaseZUOJyKrCLtSD/Iixv/hzDEUPfOCjFDgTpzf3cw
ta8+oE4wHCo1iI1/4TlPkwmXx4qSXtmw4aQPz7IDQvECgYEA8KNThCO2gsC2I9PQ
DM/8Cw0O983WCDY+oi+7JPiNAJwv5DYBqEZB1QYdj06YD16XlC/HAZMsMku1na2T
N0driwenQQWzoev3g2S7gRDoS/FCJSI3jJ+kjgtaA7Qmzlgk1TxODN+G1H91HW7t
0l7VnL27IWyYo2qRRK3jzxqUiPUCgYEAx0oQs2reBQGMVZnApD1jeq7n4MvNLcPv
t8b/eU9iUv6Y4Mj0Suo/AU8lYZXm8ubbqAlwz2VSVunD2tOplHyMUrtCtObAfVDU
AhCndKaA9gApgfb3xw1IKbuQ1u4IF1FJl3VtumfQn//LiH1B3rXhcdyo3/vIttEk
48RakUKClU8CgYEAzV7W3COOlDDcQd935DdtKBFRAPRPAlspQUnzMi5eSHMD/ISL
DY5IiQHbIH83D4bvXq0X7qQoSBSNP7Dvv3HYuqMhf0DaegrlBuJllFVVq9qPVRnK
xt1Il2HgxOBvbhOT+9in1BzA+YJ99UzC85O0Qz06A+CmtHEy4aZ2kj5hHjECgYEA
mNS4+A8Fkss8Js1RieK2LniBxMgmYml3pfVLKGnzmng7H2+cwPLhPIzIuwytXywh
2bzbsYEfYx3EoEVgMEpPhoarQnYPukrJO4gwE2o5Te6T5mJSZGlQJQj9q4ZB2Dfz
et6INsK0oG8XVGXSpQvQh3RUYekCZQkBBFcpqWpbIEsCgYAnM3DQf3FJoSnXaMhr
VBIovic5l0xFkEHskAjFTevO86Fsz1C2aSeRKSqGFoOQ0tmJzBEs1R6KqnHInicD
TQrKhArgLXX4v3CddjfTRJkFWDbE/CkvKZNOrcf1nhaGCPspRJj2KUkj1Fhl9Cnc
dn/RsYEONbwQSjIfMPkvxF+8HQ==
-----END PRIVATE KEY-----`

	block, _ = pem.Decode([]byte(privateKeyPEM))
	if block == nil || block.Type != "PRIVATE KEY" {
		log.Fatal("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("failed to parse private key: %v", err)
	}

	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		log.Fatal("not an RSA private key")
	}

	// Create a new token object
	token = jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub":   "987654321",
		"name":  "Joe Dohn",
		"admin": false,
	})

	// Sign the token with the private key
	log.Printf("Signing token with private key: %v", rsaPrivateKey)
	tokenString, err := token.SignedString(rsaPrivateKey)
	if err != nil {
		log.Fatalf("failed to sign token: %v", err)
	}

	fmt.Println("Signed Token:", tokenString)

	token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
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
