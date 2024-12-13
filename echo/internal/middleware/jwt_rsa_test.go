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
	"strings"
	"testing"
)

var (
	publicKeyPEM = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAu1SU1LfVLPHCozMxH2Mo
4lgOEePzNm0tRgeLezV6ffAt0gunVTLw7onLRnrq0/IzW7yWR7QkrmBL7jTKEn5u
+qKhbwKfBstIs+bMY2Zkp18gnTxKLxoS2tFczGkPLPgizskuemMghRniWaoLcyeh
kd3qqGElvW/VDL5AaWTg0nLVkjRo9z+40RQzuVaE8AkAFmxZzow3x+VJYKdjykkJ
0iT9wCS0DRTXu269V264Vf/3jvredZiKRkgwlL9xNAwxXFg0x/XFw005UWVRIkdg
cKWTjpBP2dPwVZ4WWC+9aGVd+Gyn1o0CLelf4rEjGoXbAAEgAqeGUxrcIlbjXfbc
mwIDAQAB
-----END PUBLIC KEY-----`

	privateKeyPEM = `-----BEGIN PRIVATE KEY-----
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

	jwtTokenString = "eyJhbGciOiJSUzI1NiIsImtpZCI6InRlc3Qta2V5LWlkIiwidHlwIjoiSldUIn0.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWUsImlhdCI6MTUxNjIzOTAyMn0.Vm4Dg5bskBWMUGTYPO6Tgge2nLJ1Oa6rkN3B-s9PsFQ9zvQQ_2V1K73X70X3zW5JiXElPRgXoJu-op0UwVt34uPdfMrMeQ5O1Ja-H6TvO2JDYMMgYz1yPp36-UY73Y7t-i2RTEXMrc9_piGtOimL7lpE5N58iQFxG4GmjyGgiLZvvczYYS0EdpJrVx4brT5pFMQ-ltPxLByw6z8jqpwzqNGssHlmzObtsKYysHOaYYfHlTDff2PeGgu6Fb5ZkNRhkQaEjCYXs1eoVVYu2w8v6FBe8sgzaQkkfiOhuiSxu0vGRv3breSe0J2xUbM1RJUjibZSnLuuFruV6wKsJWivDA"
)

func decodePEMBlock(pemData string, blockType string) ([]byte, error) {
	block, _ := pem.Decode([]byte(pemData))
	if block == nil || block.Type != blockType {
		return nil, fmt.Errorf("failed to decode PEM block containing %s", blockType)
	}
	return block.Bytes, nil
}

func parseRSAPublicKey(pemData string) (*rsa.PublicKey, error) {
	bytes, err := decodePEMBlock(pemData, "PUBLIC KEY")
	if err != nil {
		return nil, err
	}
	pub, err := x509.ParsePKIXPublicKey(bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DER encoded public key: %v", err)
	}
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}
	return rsaPub, nil
}

func parseRSAPrivateKey(pemData string) (*rsa.PrivateKey, error) {
	bytes, err := decodePEMBlock(pemData, "PRIVATE KEY")
	if err != nil {
		return nil, err
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}
	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA private key")
	}
	return rsaPrivateKey, nil
}

func intToBase64(n int) string {
	bytes := make([]byte, 4)
	binary.BigEndian.PutUint32(bytes, uint32(n))
	var start int
	for start = 0; start < len(bytes); start++ {
		if bytes[start] != 0 {
			break
		}
	}
	bytes = bytes[start:]
	return base64.StdEncoding.EncodeToString(bytes)
}

func publicKeyToPEM(pub *rsa.PublicKey) string {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		log.Fatalf("failed to marshal public key: %v", err)
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	})
	return string(pubPEM)
}

func TestRSA(t *testing.T) {
	var rsaPub *rsa.PublicKey

	t.Run("ParseRSAPublicKey", func(t *testing.T) {
		var err error
		rsaPub, err = parseRSAPublicKey(publicKeyPEM)
		if err != nil {
			t.Fatalf("failed to parse RSA public key: %v", err)
		}
	})

	t.Run("FromRSAtoJWKS", func(t *testing.T) {
		n := rsaPub.N
		e := rsaPub.E
		fmt.Printf("Modulus (n): %s\n", n.String())
		fmt.Printf("Exponent (e): %d\n", e)

		nBytes := n.Bytes()
		nBase64 := base64.URLEncoding.EncodeToString(nBytes)
		nBase64 = strings.TrimRight(nBase64, "=")
		fmt.Println("Encoded n:", nBase64)
		fmt.Println("Encoded e:", intToBase64(e))
	})

	t.Run("SignToken", func(t *testing.T) {
		rsaPrivateKey, err := parseRSAPrivateKey(privateKeyPEM)
		if err != nil {
			t.Fatalf("failed to parse RSA private key: %v", err)
		}

		token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
			"sub":   "987654321",
			"name":  "Joe Dohn",
			"admin": false,
		})

		tokenString, err := token.SignedString(rsaPrivateKey)
		if err != nil {
			t.Fatalf("failed to sign token: %v", err)
		}
		fmt.Println("Signed Token:", tokenString)

		token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return rsaPub, nil
		})
		if err != nil {
			t.Fatalf("error parsing token: %v", err)
		}
		fmt.Println("Token is valid:", token.Valid)
	})

	t.Run("ParseJWTToken", func(t *testing.T) {
		token, err := jwt.Parse(jwtTokenString, func(token *jwt.Token) (interface{}, error) {
			return rsaPub, nil
		})
		if err != nil {
			t.Fatalf("error parsing token: %v", err)
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			fmt.Println("Token Claims:")
			for key, value := range claims {
				fmt.Printf("%s: %v\n", key, value)
			}
		} else {
			t.Fatalf("invalid token claims")
		}
		fmt.Println("Token is valid:", token.Valid)
	})

	t.Run("PublicKeyToPEM", func(t *testing.T) {
		fmt.Println("Public Key PEM:")
		fmt.Println(publicKeyToPEM(rsaPub))
	})
}
