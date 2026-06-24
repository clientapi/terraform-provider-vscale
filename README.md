# Terraform & OpenTofu Provider for VScale (Selectel VDS)

[![Terraform Registry](https://img.shields.io/badge/Terraform%20Registry-clientapi%2Fvscale-blueviolet)](https://registry.terraform.io/providers/clientapi/vscale/latest)
[![OpenTofu Compatible](https://img.shields.io/badge/OpenTofu-Compatible-yellow)](https://opentofu.org/)
[![License](https://img.shields.io/badge/License-MPL%202.0-blue.svg)](LICENSE)

A Terraform and OpenTofu provider for managing cloud resources on the **VScale (Selectel VDS)** platform. Built using the modern **Terraform Plugin Framework** in Go.

---

## Features

- **Resources**:
  - `vscale_scalet`: Manage virtual private servers (OS rebuilds, root password resets, dynamic SSH key assignment, sizing plan upgrades).
  - `vscale_ssh_key`: Manage public SSH keys in the VScale panel.
  - `vscale_domain`: Manage DNS domains (supports zone file imports in BIND format).
  - `vscale_domain_record`: Full control over DNS records (A, AAAA, MX, CNAME, TXT, SRV, etc.) with in-place updates.
  - `vscale_domain_tag`: Label and group DNS domains.
  - `vscale_ptr_record`: Manage reverse DNS (PTR) records for server IP addresses.
  - `vscale_backup`: Trigger manual server backups.
- **Data Sources**:
  - `vscale_account`: Fetch account profile details, email, and current balance/bonuses.
  - `vscale_images`: List available operating system images.
  - `vscale_locations`: List datacenters and regions.
  - `vscale_rplans`: List tariff plans.
  - `vscale_prices`: Read resource pricing metrics.
  - `vscale_backups`: List existing backups.

---

## Configuration

Add the provider configuration to your `.tf` file:

```hcl
terraform {
  required_providers {
    vscale = {
      source  = "clientapi/vscale"
      version = "~> 1.0.0"
    }
  }
}

provider "vscale" {
  # token = var.vscale_token
}
```

### Authentication
The provider requires a VScale API token. You can pass it explicitly in the `provider` block as shown above, or define the `VSCALE_TOKEN` environment variable (recommended):

```bash
export VSCALE_TOKEN="your-api-token"
```

---

## Local Development & Testing

If you want to contribute or build the provider from source:

### Prerequisites
- Go 1.21+
- Terraform 1.0+ or OpenTofu 1.6+

### Building the Binary
Compile the provider binary locally:
```bash
go build -o terraform-provider-vscale.exe
```

### Local Dev Overrides
To test local builds without downloading from a registry, configure a developer override in your CLI configuration file (`tofu.rc` or `terraform.rc` / `.terraformrc`):

```hcl
provider_installation {
  dev_overrides {
    "clientapi/vscale" = "C:/Users/den67rus/GolandProjects/terraform-provider-vscale"
  }
  direct {}
}
```

Set the environment variable pointing to this config before running plan/apply:
```powershell
# Windows PowerShell
$env:TF_CLI_CONFIG_FILE="C:\Users\den67rus\GolandProjects\terraform-provider-vscale\tofurc"

# Unix/macOS
export TF_CLI_CONFIG_FILE="/path/to/tofurc"
```

---

## License

This project is licensed under the Mozilla Public License 2.0 - see the [LICENSE](LICENSE) file for details.
