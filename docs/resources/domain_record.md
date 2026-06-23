---
layout: ""
page_title: "vscale_domain_record Resource - terraform-provider-vscale"
description: |-
  Manages a DNS record for a Domain in VScale.
---

# vscale_domain_record (Resource)

Manages a DNS record for a Domain in VScale.

## Example Usage

```hcl
resource "vscale_domain_record" "www" {
  domain_id = vscale_domain.my_domain.id
  name      = "www"
  type      = "A"
  ttl       = 300
  content   = "1.2.3.4"
}
```

## Schema

### Required

- `domain_id` (Number) The ID of the domain to attach this record to.
- `name` (String) The name of the record (e.g., `www`, `mail`, or `@` for root).
- `type` (String) Record type (e.g., `SOA`, `NS`, `A`, `AAAA`, `CNAME`, `SRV`, `MX`, `TXT`, `SPF`).

### Optional

- `content` (String) Value of the record (e.g., an IP address). Omitted for SRV records.
- `port` (Number) Port for SRV records.
- `priority` (Number) Priority for MX and SRV records.
- `target` (String) Target hostname for SRV records.
- `ttl` (Number) Time to Live (TTL) in seconds. Min 60, Max 604800. Defaults to `300`.
- `weight` (Number) Weight for SRV records.

### Computed

- `id` (String) Unique identifier of the domain record.
