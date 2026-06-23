---
layout: ""
page_title: "vscale_ptr_record Resource - terraform-provider-vscale"
description: |-
  Manages a DNS PTR (Reverse) record in VScale.
---

# vscale_ptr_record (Resource)

Manages a DNS PTR (Reverse) record in VScale.

## Example Usage

```hcl
resource "vscale_ptr_record" "my_ptr" {
  ip      = "82.148.16.208"
  content = "example.com"
}
```

## Schema

### Required

- `content` (String) Domain name (value of the reverse record).
- `ip` (String) IP address for which the PTR record will be created.

### Computed

- `id` (String) Unique identifier of the PTR record.
