---
layout: ""
page_title: "vscale_ssh_key Resource - terraform-provider-vscale"
description: |-
  Manages an SSH key in VScale.
---

# vscale_ssh_key (Resource)

Manages an SSH key in VScale.

## Example Usage

```hcl
resource "vscale_ssh_key" "my_key" {
  name = "production-key"
  key  = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQ..."
}
```

## Schema

### Required

- `key` (String) The public SSH key.
- `name` (String) Name of the SSH key.

### Computed

- `id` (String) Unique identifier of the SSH key.
