---
layout: ""
page_title: "vscale_backups Data Source - terraform-provider-vscale"
description: |-
  Provides list of all backups in VScale.
---

# vscale_backups (Data Source)

Provides list of all backups in VScale.

## Example Usage

```hcl
data "vscale_backups" "all" {}
```

## Schema

### Computed

- `backups` (List of Object) List of all backups. (see [below for nested schema](#nestedatt--backups))

<a id="nestedatt--backups"></a>
### Nested Schema for `backups`

Computed:

- `active` (Boolean) Whether the backup is active.
- `created` (String) Creation timestamp.
- `id` (String) Identifier of the backup.
- `location` (String) Storage datacenter location.
- `locked` (Boolean) Whether the backup is locked.
- `name` (String) Name of the backup.
- `scalet_id` (Number) ID of the source Scalet.
- `size` (Number) Backup size in GB.
- `status` (String) Current status of the backup.
- `template` (String) OS template from which server was built.
