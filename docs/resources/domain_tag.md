---
layout: ""
page_title: "vscale_domain_tag Resource - terraform-provider-vscale"
description: |-
  Manages a DNS domain tag in VScale.
---

# vscale_domain_tag (Resource)

Manages a DNS domain tag in VScale.

## Example Usage

```hcl
resource "vscale_domain_tag" "production" {
  name    = "production"
  domains = ["example.com"]
}
```

## Schema

### Required

- `name` (String) Name of the tag.

### Optional

- `domains` (List of String) Domains attached to this tag.

### Computed

- `id` (String) Unique identifier of the domain tag.
