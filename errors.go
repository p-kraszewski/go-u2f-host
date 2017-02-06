package u2fhost

// #cgo LDFLAGS: -lu2f-host
// #include <stdlib.h>
// #include <u2f-host/u2f-host.h>
import "C"

import "errors"

var (
	// ==== Common errors ==================================

	// ErrMemory Memory allocation error
	ErrMemory = errors.New("Memory error")

	// ErrJSON Bad JSON format
	ErrJSON = errors.New("Json error")

	// ErrBase64 Bad Base64 format
	ErrBase64 = errors.New("Base64 error")

	// ==== Server errors ==================================

	// //ErrCrypto Error in cryptography
	// ErrCrypto = errors.New("Cryptographic error")
	//
	// // ErrOrigin Origin does not match
	// ErrOrigin = errors.New("Origin mismatch")
	//
	// // ErrChallenge Challenge error
	// ErrChallenge = errors.New("Challenge error")
	//
	// // ErrSignature Signature mismatch
	// ErrSignature = errors.New("Signature mismatch")
	//
	// // ErrFormat Message format error
	// ErrFormat = errors.New("Message format error")

	// ==== Host errors ==================================

	// ErrTransport Transport error
	ErrTransport = errors.New("Transport error")
	// ErrNoU2FDevice No U2F device
	ErrNoU2FDevice = errors.New("No U2F device")
	// ErrAuthentication Autherntication error
	ErrAuthentication = errors.New("Autherntication error")
	// ErrTimeout Timeout error
	ErrTimeout = errors.New("Timeout error")
	// ErrSize Size error
	ErrSize = errors.New("Size error")

	// ==== Errors internal to Go binding ==================================

	// ErrInvalidPubKey Invalid PubKey format
	ErrInvalidPubKey = errors.New("Invalid PubKey format")

	// ErrDeviceNumber Invalid device number
	ErrDeviceNumber = errors.New("Invalid device number")

	//ErrOther Unknown error
	ErrOther = errors.New("Unknown error")
)

var errorList = []error{ErrMemory, ErrTransport, ErrJSON, ErrBase64, ErrNoU2FDevice, ErrAuthentication, ErrTimeout, ErrSize}

func iToErr(e C.u2fh_rc) error {
	if e >= 0 {
		return nil
	}

	if e < -8 {
		return ErrOther
	}

	return errorList[-(e - 1)]
}
