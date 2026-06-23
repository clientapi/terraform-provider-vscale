---
layout: ""
page_title: "vscale_domain Resource - terraform-provider-vscale"
description: |-
  Manages a DNS domain in VScale.
---

# vscale_domain (Resource)

Manages a DNS domain in VScale.

## Example Usage

```hcl
resource "vscale_domain" "my_domain" {
  name      = "example.com"
  bind_zone = "$ORIGIN example.com.\n..."
}
```

## Schema

### Required

- `name` (String) The domain name (e.g., `example.com`).

### Optional

- `bind_zone` (String) Optional zone file content in BIND format.

### Computed

- `id` (String) Unique identifier of the Domain.
