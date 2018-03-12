package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/sl1pm4t/terraform-provider-tfstate/tfstate"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: tfstate.Provider})
}
