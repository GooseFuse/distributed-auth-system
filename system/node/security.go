package node

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// SecurityManager manages security aspects like TLS and signatures
type SecurityManager struct {
	tlsConfig  *tls.Config
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
	certPool   *x509.CertPool
	mutex      sync.RWMutex
}

// NewSecurityManager creates a new SecurityManager
func NewSecurityManager(nodeID string, certDir string) (*SecurityManager, error) {
	// Create certificate directory if it doesn't exist
	err := os.MkdirAll(certDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate directory: %w", err)
	}

	// Check if certificates already exist
	certPath := filepath.Join(certDir, "cert.pem")
	keyPath := filepath.Join(certDir, "key.pem")
	certExists := fileExists(certPath) && fileExists(keyPath)

	var privateKey *ecdsa.PrivateKey
	var cert tls.Certificate

	if !certExists {
		// Generate new key pair and certificate
		privateKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			return nil, fmt.Errorf("failed to generate private key: %w", err)
		}

		// Create self-signed certificate
		template := x509.Certificate{
			SerialNumber: big.NewInt(1),
			Subject: pkix.Name{
				Organization: []string{"Distributed Auth System"},
				CommonName:   nodeID,
			},
			NotBefore:             time.Now(),
			NotAfter:              time.Now().Add(365 * 24 * time.Hour), // 1 year
			KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
			ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
			BasicConstraintsValid: true,
			IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
			DNSNames:              []string{"localhost"},
		}

		derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
		if err != nil {
			return nil, fmt.Errorf("failed to create certificate: %w", err)
		}

		// Save certificate and key to files
		certOut, err := os.Create(certPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open cert.pem for writing: %w", err)
		}
		err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
		certOut.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to write cert.pem: %w", err)
		}

		keyOut, err := os.Create(keyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open key.pem for writing: %w", err)
		}
		keyBytes, err := x509.MarshalECPrivateKey(privateKey)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal private key: %w", err)
		}
		err = pem.Encode(keyOut, &pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes})
		keyOut.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to write key.pem: %w", err)
		}

		// Load the newly created certificate
		cert, err = tls.LoadX509KeyPair(certPath, keyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load certificate: %w", err)
		}
	} else {
		// Load existing certificate
		cert, err = tls.LoadX509KeyPair(certPath, keyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load certificate: %w", err)
		}

		// Extract private key from certificate
		key, ok := cert.PrivateKey.(*ecdsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("private key is not an ECDSA key")
		}
		privateKey = key
	}

	// Create certificate pool
	certPool := x509.NewCertPool()
	x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}
	certPool.AddCert(x509Cert)

	// Create TLS config
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
		RootCAs:      certPool,
		MinVersion:   tls.VersionTLS12,
	}

	return &SecurityManager{
		tlsConfig:  tlsConfig,
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
		certPool:   certPool,
	}, nil
}

// GetTLSConfig returns the TLS configuration
func (sm *SecurityManager) GetTLSConfig() *tls.Config {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()
	return sm.tlsConfig
}

// SignData signs data with the private key
func (sm *SecurityManager) SignData(data []byte) ([]byte, error) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	hash := sha256.Sum256(data)
	return ecdsa.SignASN1(rand.Reader, sm.privateKey, hash[:])
}

// VerifySignature verifies a signature with a public key
func (sm *SecurityManager) VerifySignature(data []byte, signature []byte, publicKey *ecdsa.PublicKey) bool {
	hash := sha256.Sum256(data)
	return ecdsa.VerifyASN1(publicKey, hash[:], signature)
}

// GetPublicKeyBytes returns the public key as bytes
func (sm *SecurityManager) GetPublicKeyBytes() ([]byte, error) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	return x509.MarshalPKIXPublicKey(sm.publicKey)
}

// GetPublicKeyHex returns the public key as a hex string
func (sm *SecurityManager) GetPublicKeyHex() (string, error) {
	keyBytes, err := sm.GetPublicKeyBytes()
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(keyBytes), nil
}

// ParsePublicKey parses a public key from bytes
func ParsePublicKey(keyBytes []byte) (*ecdsa.PublicKey, error) {
	key, err := x509.ParsePKIXPublicKey(keyBytes)
	if err != nil {
		return nil, err
	}
	publicKey, ok := key.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an ECDSA public key")
	}
	return publicKey, nil
}

// ParsePublicKeyHex parses a public key from a hex string
func ParsePublicKeyHex(keyHex string) (*ecdsa.PublicKey, error) {
	keyBytes, err := hex.DecodeString(keyHex)
	if err != nil {
		return nil, err
	}
	return ParsePublicKey(keyBytes)
}

// Helper function to check if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
