---
layout: ""
page_title: "vscale_backup Resource - terraform-provider-vscale"
description: |-
  Manages a server backup in VScale.
---

# vscale_backup (Resource)

Manages a server backup in VScale. Creating a backup triggers a backup of the source Scalet.

## Example Usage

```hcl
resource "vscale_backup" "my_backup" {
  scalet_id = 81661558
  name      = "manual-backup-before-deploy"
}
```

## Schema

### Required

- `name` (String) Name of the backup.
- `scalet_id` (Number) ID of the source Scalet.

### Computed

- `created` (String) Creation timestamp.
- `id` (String) Unique identifier of the backup.
- `location` (String) Datacenter location where backup is stored.
- `size` (Number) Backup size in GB.
- `status` (String) Status of the backup.
- `template` (String) OS template from which server was built.
