package u2fhost

// #cgo LDFLAGS: -lu2f-host
// #include <stdlib.h>
// #include <u2f-host/u2f-host.h>
import "C"
import "unsafe"

// Mode Library initialization mode
type Mode C.u2fh_initflags

const (
	// Production mode (no debuginfo printed)
	Production Mode = 0

	// Debug mode (debuginfo printed to console)
	Debug = 1
)

// Host U2F marker struct
type Host struct{}

// Context object
type Context struct {
	devs     *C.u2fh_devs
	maxIndex C.unsigned
}

// Start Must be called successfully before using any other functions.
func Start(mode Mode) (Host, error) {
	rc := C.u2fh_global_init(C.u2fh_initflags(mode))
	return Host{}, iToErr(rc)
}

// Stop Must be called to clean-up system.
func (host Host) Stop() {
	C.u2fh_global_done()
}

// Open Create context before registration/authentication calls.
func (host Host) Open() (*Context, error) {
	ctx := &Context{}

	rc := C.u2fh_devs_init(&ctx.devs)
	if rc != 0 {
		return nil, iToErr(rc)
	}

	rc = C.u2fh_devs_discover(ctx.devs, &ctx.maxIndex)
	if rc != 0 {
		ctx.Close()
		return nil, iToErr(rc)
	}

	return ctx, nil

}

// Close Destroy context
func (ctx *Context) Close() {
	C.u2fh_devs_done(ctx.devs)
}

// Register device in server
func (ctx *Context) Register(challenge string, origin string, present bool) (string, error) {
	var output *C.char
	var pres C.u2fh_cmdflags
	orig := C.CString(origin)
	defer C.free(unsafe.Pointer(orig))
	chal := C.CString(challenge)
	defer C.free(unsafe.Pointer(chal))

	if present {
		pres = 1
	}

	rc := C.u2fh_register(ctx.devs, chal, orig, &output, pres)
	if rc != 0 {
		return "", iToErr(rc)
	}
	return C.GoString(output), nil
}

// Authenticate device against server
func (ctx *Context) Authenticate(challenge string, origin string, present bool) (string, error) {
	var output *C.char
	var pres C.u2fh_cmdflags
	orig := C.CString(origin)
	defer C.free(unsafe.Pointer(orig))
	chal := C.CString(challenge)
	defer C.free(unsafe.Pointer(chal))

	if present {
		pres = 1
	}

	rc := C.u2fh_authenticate(ctx.devs, chal, orig, &output, pres)
	if rc != 0 {
		return "", iToErr(rc)
	}
	return C.GoString(output), nil
}

// IsAlive check if detected device is still active
func (ctx *Context) IsAlive(num uint8) error {
	if C.unsigned(num) > ctx.maxIndex {
		return ErrDeviceNumber
	}
	rc := C.u2fh_is_alive(ctx.devs, C.unsigned(num))
	if rc != 0 {
		return nil
	}
	return ErrNoU2FDevice
}

// GetDescription get textual description of device
func (ctx *Context) GetDescription(num uint8) (string, error) {
	if C.unsigned(num) > ctx.maxIndex {
		return "", ErrDeviceNumber
	}

	var buffer [1024]byte

	var buflen C.size_t = C.size_t(len(buffer))
	bufptr := (*C.char)(unsafe.Pointer(&buffer[0]))

	rc := C.u2fh_get_device_description(ctx.devs, C.unsigned(num), bufptr, &buflen)
	if rc != 0 {
		return "", iToErr(rc)
	}
	return string(buffer[:buflen]), nil
}
