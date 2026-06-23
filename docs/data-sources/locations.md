---
layout: ""
page_title: "vscale_locations Data Source - terraform-provider-vscale"
description: |-
  Provides list of available datacenters/locations in VScale.
---

# vscale_locations (Data Source)

Provides list of available datacenters/locations in VScale.

## Example Usage

```hcl
data "vscale_locations" "all" {}
```

## Schema

### Computed

- `locations` (List of Object) List of available locations. (see [below for nested schema](#nestedatt--locations))

<a id="nestedatt--locations"></a>
### Nested Schema for `locations`

Computed:

- `active` (Boolean) Whether the location is active.
- `description` (String) Human-readable description.
- `id` (String) Identifier of the location (e.g. `spb0`, `msk0`).
- `private_networking` (Boolean) Whether private networking is supported in this location.
- `rplans` (List of String) Tariff plans available in this location.
- `templates` (List of String) Templates available in this location.
