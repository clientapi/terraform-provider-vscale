package main

import (
	"context"
	"flag"
	"log"

	"terraform-provider-vscale/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Generate the Terraform provider
func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/vscale/vscale",
	}

	err := providerserver.Serve(context.Background(), provider.New("1.0.0"), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
