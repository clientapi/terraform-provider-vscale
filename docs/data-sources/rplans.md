---
layout: ""
page_title: "vscale_rplans Data Source - terraform-provider-vscale"
description: |-
  Provides list of available pricing/tariff plans (rplans) in VScale.
---

# vscale_rplans (Data Source)

Provides list of available pricing/tariff plans (rplans) in VScale.

## Example Usage

```hcl
data "vscale_rplans" "all" {}
```

## Schema

### Computed

- `rplans` (List of Object) List of available plans. (see [below for nested schema](#nestedatt--rplans))

<a id="nestedatt--rplans"></a>
### Nested Schema for `rplans`

Computed:

- `addresses` (Number) Number of IP addresses included.
- `cpus` (Number) CPU core count.
- `disk` (Number) Disk size.
- `id` (String) Identifier of the plan (e.g. `small`, `medium`).
- `locations` (List of String) Locations where this plan is available.
- `memory` (Number) RAM size in MB.
- `templates` (List of String) Templates compatible with this plan.
