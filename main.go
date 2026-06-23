// Copyright Aviatrix Systems, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/AviatrixSystems/terraform-provider-avxcloud/internal/provider"
)

// version is overridden at build time.
var version = "dev"

// Generate documentation.
//go:generate go tool tfplugindocs generate -provider-name avxcloud

func main() {
	debug := true
	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/aviatrixsystems/avxcloud",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)
	if err != nil {
		log.Fatal(err.Error())
	}
}
