---
layout: ""
page_title: "vscale_prices Data Source - terraform-provider-vscale"
description: |-
  Provides price list for VScale resources.
---

# vscale_prices (Data Source)

Provides price list for VScale resources.

## Example Usage

```hcl
data "vscale_prices" "current" {}
```

## Schema

### Computed

- `backup_price` (Number) Backup storage cost per GB.
- `huge_hour` (Number) Cost per hour for `huge` configuration.
- `huge_month` (Number) Cost per month for `huge` configuration.
- `large_hour` (Number) Cost per hour for `large` configuration.
- `large_month` (Number) Cost per month for `large` configuration.
- `medium_hour` (Number) Cost per hour for `medium` configuration.
- `medium_month` (Number) Cost per month for `medium` configuration.
- `monster_hour` (Number) Cost per hour for `monster` configuration.
- `monster_month` (Number) Cost per month for `monster` configuration.
- `period` (String) The pricing period start date.
- `small_hour` (Number) Cost per hour for `small` configuration.
- `small_month` (Number) Cost per month for `small` configuration.
