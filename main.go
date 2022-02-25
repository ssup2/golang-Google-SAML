package main

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"net/url"

	"github.com/crewjam/saml/samlsp"
)

func samlRequestPrinter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Header : %+v\n", r.Header)
		fmt.Printf("Body : %+v\n", r.Body)
		next.ServeHTTP(w, r)
	})
}

func echoSession(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%v\n", samlsp.SessionFromContext(r.Context()))
}

func main() {
	// Load certificate keypair
	keyPair, err := tls.LoadX509KeyPair("myservice.cert", "myservice.key")
	if err != nil {
		panic(err)
	}
	keyPair.Leaf, err = x509.ParseCertificate(keyPair.Certificate[0])
	if err != nil {
		panic(err)
	}

	// Get identity provider info from identity provider meta URL
	idpMetadataURL, err := url.Parse("http://localhost:8080/realms/ssup2/protocol/saml/descriptor")
	if err != nil {
		panic(err)
	}
	idpMetadata, err := samlsp.FetchMetadata(context.Background(), http.DefaultClient,
		*idpMetadataURL)
	if err != nil {
		panic(err)
	}

	// Get SAML service provider middleware
	rootURL, err := url.Parse("http://localhost:8000")
	if err != nil {
		panic(err)
	}
	samlSP, _ := samlsp.New(samlsp.Options{
		URL:         *rootURL,
		Key:         keyPair.PrivateKey.(*rsa.PrivateKey),
		Certificate: keyPair.Leaf,
		IDPMetadata: idpMetadata,
	})

	// Set SAML's metadata and ACS (Assertion Consumer Service) endpoint with SAML request printer
	http.Handle("/saml/", samlRequestPrinter(samlSP))

	// Set session handler to print session info
	app := http.HandlerFunc(echoSession)
	http.Handle("/session", samlSP.RequireAccount(app))

	// Serve HTTP
	http.ListenAndServe(":8000", nil)
}
