---
layout: ""
page_title: "vscale_scalet Resource - terraform-provider-vscale"
description: |-
  Manages a Scalet (Virtual Private Server) in VScale.
---

# vscale_scalet (Resource)

Manages a Scalet (Virtual Private Server) in VScale.

## Example Usage

```hcl
resource "vscale_scalet" "my_server" {
  name      = "production-server"
  make_from = "ubuntu_20.04_64_001_master"
  rplan     = "medium"
  location  = "spb0"
  do_start  = true
  password  = "SecurePassword123!"

  keys = [
    12345
  ]
}
```

## Schema

### Required

- `location` (String) Datacenter location ID (e.g. `spb0`, `msk0`).
- `make_from` (String) The ID of the OS image or backup from which to build the Scalet. Changing this triggers rebuild.
- `name` (String) Name of the Scalet.
- `rplan` (String) Tariff plan ID (e.g. `small`, `medium`, `large`). Changing this triggers plan upgrade.

### Optional

- `do_start` (Boolean) Whether to start the Scalet immediately after creation. Defaults to `false`.
- `keys` (List of Number) List of SSH key IDs to load onto the Scalet. Can be updated dynamically.
- `password` (String, Sensitive) Root password for the Scalet (if keys are not used). Changing this triggers rebuild.

### Computed

- `id` (String) Unique identifier (CTID) of the Scalet.
- `private_address` (Object) Private network interface details.
- `public_address` (Object) Public network interface details.
- `status` (String) Status of the Scalet (e.g., `started`, `stopped`, `defined`).
