/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: BSD-2-Clause */

package server

import (
	"fmt"
	"net"
	"net/http"

	"github.com/vmware/vsphere-automation-sdk-go/runtime/l10n"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/log"
	"golang.org/x/net/context"
)

// Server Wraps http server and provides basic functionality
// such as server start and stop
type Server struct {
	srv *http.Server
}

// NewServer returns Server object
func NewServer(portnumber string, handler http.Handler) *Server {
	return &Server{srv: &http.Server{Addr: portnumber, Handler: handler}}
}

// Start Starts http server
func (s *Server) Start() chan error {
	errorChannel := make(chan error)
	go func() {
		// catch when server panics
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				var errorMsg string
				if ok {
					errorMsg = err.Error()
				} else {
					errorMsg = fmt.Sprintf("Server error occured: %#v", r)
				}

				errorChannel <- l10n.NewRuntimeError(
					"vapi.protocol.server.error",
					map[string]string{"err": errorMsg})
			}
		}()
		if err := s.srv.ListenAndServe(); err != nil {
			errorChannel <- l10n.NewRuntimeError(
				"vapi.protocol.server.error",
				map[string]string{"err": err.Error()})
		}
	}()

	log.Infof("HTTP Server starting on %s", s.srv.Addr)
	return errorChannel
}

// WaitForRunningPort waits until a server port gets used by the OS.
// Should be used after call to Start method. Returned error channel from Start
// method is expected as input parameter.
// Closes server in case of error.
func (s *Server) WaitForRunningPort(errChannel chan error, seconds int) (bool, error) {
	check, err := WaitForFunc(seconds, func() (bool, error) {
		select {
		case err := <-errChannel:
			return false, err
		default:
			_, err := net.Dial("tcp", s.Address())
			if err == nil {
				return true, nil
			}
			return false, nil
		}
	})
	if err != nil || !check {
		s.Stop()
		return false, err
	}

	log.Info("Server started successfully on ", s.Address())
	return true, nil
}

// Stop Stops http server gracefully
func (s *Server) Stop() {
	log.Info("Stopping Http Server")
	if err := s.srv.Shutdown(context.TODO()); err != nil {
		panic(err) // failure/timeout shutting down the server gracefully
	}
	log.Infof("Server at %s stopped.", s.srv.Addr)
}

// Address Returns server address
func (s *Server) Address() string {
	return s.srv.Addr
}
