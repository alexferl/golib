package server

import (
	"crypto/tls"

	"github.com/klauspost/cpuid/v2"
)

// defaultCipherSuites is the ordered list of all the cipher
// suites we want to support by default, assuming AES-NI
// (hardware acceleration for AES).
var defaultCipherSuitesWithAESNI = []uint16{
	tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
	tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
}

// defaultCipherSuites is the ordered list of all the cipher
// suites we want to support by default, assuming lack of
// AES-NI (NO hardware acceleration for AES).
var defaultCipherSuitesWithoutAESNI = []uint16{
	tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
	tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
	tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
}

// getOptimalDefaultCipherSuites returns an appropriate cipher
// suite to use depending on the hardware support for AES.
func getOptimalDefaultCipherSuites() []uint16 {
	if cpuid.CPU.Supports(cpuid.AESNI) {
		return defaultCipherSuitesWithAESNI
	}
	return defaultCipherSuitesWithoutAESNI
}

// defaultCurves is the list of only the curves or key exchange
// mechanisms we want to use by default. The order is irrelevant.
var defaultCurves = []tls.CurveID{
	tls.X25519MLKEM768,
	tls.X25519,
	tls.CurveP256,
}
