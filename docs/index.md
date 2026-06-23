---
layout: ""
page_title: "Provider: VScale"
description: |-
  The VScale provider is used to interact with the resources supported by VScale (Selectel VDS) API.
---

# VScale Provider

The VScale provider is used to interact with VScale (Selectel VDS) cloud resources.
The provider must be configured with an API token before it can be used.

## Example Usage

```hcl
terraform {
  required_providers {
    vscale = {
      source  = "vscale/vscale"
      version = "1.0.0"
    }
  }
}

provider "vscale" {
  # token = var.vscale_token
}

resource "vscale_scalet" "my_server" {
  name      = "production-server"
  make_from = "ubuntu_20.04_64_001_master"
  rplan     = "medium"
  location  = "spb0"
  do_start  = true
}
```

## Schema

### Optional

- `token` (String, Sensitive) The VScale API token. Can also be set via the `VSCALE_TOKEN` environment variable.
