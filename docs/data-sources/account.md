---
layout: ""
page_title: "vscale_account Data Source - terraform-provider-vscale"
description: |-
  Fetches information about the current VScale account and billing balance.
---

# vscale_account (Data Source)

Fetches information about the current VScale account and billing balance.

## Example Usage

```hcl
data "vscale_account" "current" {}

output "balance" {
  value = data.vscale_account.current.balance
}
```

## Schema

### Computed

- `actdate` (String) Account activation date.
- `balance` (Number) Current balance amount.
- `bonus` (Number) Current bonus amount.
- `country` (String) User's country.
- `email` (String) Registered email address.
- `face_id` (String) Internal type identifier for physical/legal entity.
- `id` (String) Account ID.
- `middlename` (String) User's middle name.
- `mobile` (String) Registered mobile number.
- `name` (String) User's first name.
- `state` (String) Account state (1 for active, 0 for inactive).
- `surname` (String) User's last name.
