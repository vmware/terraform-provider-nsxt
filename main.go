// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/vmware/terraform-provider-nsxt/nsxt"
)

// func main() {
// 	debugMode := flag.Bool("debug", false, "Enable debug mode for debuggers like Delve (dlv).")
// 	flag.Parse()

// 	opts := &plugin.ServeOpts{
// 		ProviderFunc: func() *schema.Provider {
// 			return nsxt.Provider()
// 		},
// 	}

// 	if *debugMode {
// 		opts.Debug = true
// 		opts.ProviderAddr = "registry.terraform.io/vmware/nsxt"
// 	}

// 	plugin.Serve(opts)
// }

func main() {
	ctx := context.Background()

	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	upgradedSdkServer, err := tf5to6server.UpgradeServer(
		ctx,
		nsxt.Provider().GRPCProvider, // Old terraform-plugin-sdk provider
	)

	if err != nil {
		log.Fatal(err)
	}

	providers := []func() tfprotov6.ProviderServer{
		providerserver.NewProtocol6(nsxt.NewFrameworkProvider()),
		func() tfprotov6.ProviderServer {
			return upgradedSdkServer
		},
	}

	muxServer, err := tf6muxserver.NewMuxServer(ctx, providers...)

	if err != nil {
		log.Fatal(err)
	}

	var serveOpts []tf6server.ServeOpt

	if debug {
		serveOpts = append(serveOpts, tf6server.WithManagedDebug())
	}

	err = tf6server.Serve(
		"registry.terraform.io/<namespace>/<provider_name>",
		muxServer.ProviderServer,
		serveOpts...,
	)

	if err != nil {
		log.Fatal(err)
	}
}
