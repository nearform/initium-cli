package k8s

import (

	_"k8s.io/client-go/util/cert"
	"fmt"
	"testing"
	_"crypto/rsa"
	
)

var k8sInfo = map[string]string{
	"caCrt":  `
	-----BEGIN CERTIFICATE-----
MIIBkjCCATegAwIBAgIIRvqLB4BCntMwCgYIKoZIzj0EAwIwIzEhMB8GA1UEAwwY
azNzLWNsaWVudC1jYUAxNjgzMjA2OTEyMB4XDTIzMDUwNDEzMjgzMloXDTI0MDUw
MzEzMjgzMlowMDEXMBUGA1UEChMOc3lzdGVtOm1hc3RlcnMxFTATBgNVBAMTDHN5
c3RlbTphZG1pbjBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABIK6xKvVGY+MUdH0
MahZ+B4TPZC6EhboO6lPL+7eNl6W+Ar60Hje+QKATxciNX0p4r0CVs9HmWYTgQPQ
uyRieQijSDBGMA4GA1UdDwEB/wQEAwIFoDATBgNVHSUEDDAKBggrBgEFBQcDAjAf
BgNVHSMEGDAWgBT/yKboRYR/66YB21hnmZK9CPZfFTAKBggqhkjOPQQDAgNJADBG
AiEAhoLrnKOq/6Cu9ZSf9GtseSEPekMYUBy0wnkVOJ/+nSICIQCXcHjqVaWxmO8U
DZvlOseXlHljNP54oFOuXJ4iFpKK+Q==
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIBeDCCAR2gAwIBAgIBADAKBggqhkjOPQQDAjAjMSEwHwYDVQQDDBhrM3MtY2xp
ZW50LWNhQDE2ODMyMDY5MTIwHhcNMjMwNTA0MTMyODMyWhcNMzMwNTAxMTMyODMy
WjAjMSEwHwYDVQQDDBhrM3MtY2xpZW50LWNhQDE2ODMyMDY5MTIwWTATBgcqhkjO
PQIBBggqhkjOPQMBBwNCAAR8MKRsvffcDVlIv2ffOFuZbmoQdfpJVyWS72lrFHHZ
XmmIZPFMJbhgJ8/B20nV1zMSwYFk33AEIIWI+elS27Xno0IwQDAOBgNVHQ8BAf8E
BAMCAqQwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQU/8im6EWEf+umAdtYZ5mS
vQj2XxUwCgYIKoZIzj0EAwIDSQAwRgIhANmzIE9boUJ/t0/7d6QxbRBS0hmXkwap
OHK+JHIkIJUWAiEAw3vv7B3fDqeh5USkDctTf6fNTVj3o9LQ8A3o+DzUzN0=
-----END CERTIFICATE-----` ,
	"endpoint":  "https://127.0.0.1:6443",
	"token": `-----BEGIN EC PRIVATE KEY-----
	MHcCAQEEIGppDcvonjlKLWtpKQAsTzd5tWNcxnq1nGaRPt2n+PXToAoGCCqGSM49
	AwEHoUQDQgAEgrrEq9UZj4xR0fQxqFn4HhM9kLoSFug7qU8v7t42Xpb4CvrQeN75
	AoBPFyI1fSnivQJWz0eZZhOBA9C7JGJ5CA==
	-----END EC PRIVATE KEY-----`,
}

func TestConfig(t *testing.T) {

	// key, err := rsa.GenerateKey(cryptorand.Reader, 2048)
	// if err != nil {
	// 	t.Fatalf("rsa key failed to generate: %v", err)
	// }
	
	_, err := Config(k8sInfo["endpoint"], k8sInfo["token"], []byte(k8sInfo["caCrt"]))

	if err != nil {
		t.Fatalf(fmt.Sprintf("Error: %s", err))
	}

	
}


