package client

import (
	"encoding/asn1"
	"fmt"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/l10n"
	"net"
	"runtime"
	"syscall"

	"net/http"
	"reflect"
	"strings"

	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/lib"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client/metadata"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/security"
)

// GetRuntimeUserAgentHeader returns User-Agent header for go runtime
func GetRuntimeUserAgentHeader() string {
	return fmt.Sprintf("vAPI/%s Go/%s (%s; %s)", metadata.RuntimeVersion, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}

// CopyContextsToHeaders sets request headers using execution context properties
func CopyContextsToHeaders(ctx *core.ExecutionContext, req *http.Request) {
	appCtx := ctx.ApplicationContext()
	secCtx := ctx.SecurityContext()

	if appCtx != nil {
		for key, value := range appCtx.GetAllProperties() {
			keyLowerCase := strings.ToLower(key)
			switch keyLowerCase {
			case lib.HTTP_USER_AGENT_HEADER:
				// Prepend application user agent to runtime user agent
				vapiUserAgent := req.Header.Get(lib.HTTP_USER_AGENT_HEADER)
				userAgent := fmt.Sprintf("%s %s", *value, vapiUserAgent)
				req.Header.Set(lib.HTTP_USER_AGENT_HEADER, userAgent)
			case lib.HTTP_ACCEPT_LANGUAGE:
				req.Header.Set(lib.HTTP_ACCEPT_LANGUAGE, *value)
			default:
				req.Header.Set(lib.VAPI_HEADER_PREFIX+keyLowerCase, *value)
			}
		}
	}

	if secCtx != nil {
		if secCtx.Property(security.AUTHENTICATION_SCHEME_ID) == security.SESSION_SCHEME_ID {
			if sessionId, ok := secCtx.Property(security.SESSION_ID).(string); ok {
				req.Header.Set(lib.VAPI_SESSION_HEADER, sessionId)
			} else {
				log.Errorf("Invalid session ID in security context. Skipping setting request header. Expected string but was %s",
					reflect.TypeOf(secCtx.Property(security.SESSION_ID)))
			}
		}
	}
}

func getVAPIError(err error) *data.ErrorValue {
	log.Error(err)
	if netError, ok := err.(net.Error); ok && netError.Timeout() {
		err := l10n.NewRuntimeError("vapi.server.timedout", map[string]string{"errMsg": err.Error()})
		errVal := bindings.CreateErrorValueFromMessages(bindings.TIMEDOUT_ERROR_DEF, []error{err})
		return errVal
	}

	switch t := err.(type) {
	case *net.OpError:
		err := l10n.NewRuntimeError("vapi.server.unavailable", map[string]string{"errMsg": err.Error()})
		errVal := bindings.CreateErrorValueFromMessages(bindings.SERVICE_UNAVAILABLE_ERROR_DEF, []error{err})
		return errVal
	case syscall.Errno:
		if t == syscall.ECONNREFUSED {
			err := l10n.NewRuntimeError("vapi.server.unavailable", map[string]string{"errMsg": err.Error()})
			errVal := bindings.CreateErrorValueFromMessages(bindings.SERVICE_UNAVAILABLE_ERROR_DEF, []error{err})
			return errVal
		}
	case asn1.SyntaxError:
		err := l10n.NewRuntimeError("vapi.security.authentication.certificate.invalid", map[string]string{"errMsg": err.Error()})
		errVal := bindings.CreateErrorValueFromMessages(bindings.SERVICE_UNAVAILABLE_ERROR_DEF, []error{err})
		return errVal
	case asn1.StructuralError:
		err := l10n.NewRuntimeError("vapi.security.authentication.certificate.invalid", map[string]string{"errMsg": err.Error()})
		errVal := bindings.CreateErrorValueFromMessages(bindings.SERVICE_UNAVAILABLE_ERROR_DEF, []error{err})
		return errVal
	}

	err = l10n.NewRuntimeError("vapi.protocol.client.request.error", map[string]string{"errMsg": err.Error()})
	errVal := bindings.CreateErrorValueFromMessages(bindings.SERVICE_UNAVAILABLE_ERROR_DEF, []error{err})
	return errVal
}
