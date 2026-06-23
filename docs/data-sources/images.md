---
layout: ""
page_title: "vscale_images Data Source - terraform-provider-vscale"
description: |-
  Provides list of available OS images in VScale.
---

# vscale_images (Data Source)

Provides list of available OS images in VScale.

## Example Usage

```hcl
data "vscale_images" "all" {}
```

## Schema

### Computed

- `images` (List of Object) List of available OS images. (see [below for nested schema](#nestedatt--images))

<a id="nestedatt--images"></a>
### Nested Schema for `images`

Computed:

- `active` (Boolean) Whether the image is active.
- `description` (String) Human-readable description.
- `id` (String) Identifier of the image (OS name).
- `locations` (List of String) Locations where this image is available.
- `rplans` (List of String) Tariff plans compatible with this image.
- `size` (Number) Size of the image in MB.
