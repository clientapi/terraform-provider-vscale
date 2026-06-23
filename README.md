# VScale (Selectel VDS) Terraform & OpenTofu Provider

This is a custom Terraform and OpenTofu provider for managing cloud resources on VScale (Selectel VDS) platform. Built using the modern **Terraform Plugin Framework**.

## Features

- **Resources**:
  - `vscale_scalet`: Manage virtual private servers (rebuild OS, update SSH keys, upgrade plans).
  - `vscale_ssh_key`: Manage SSH keys in VScale panel.
  - `vscale_domain`: Manage DNS domains.
  - `vscale_domain_record`: Manage resource DNS records (A, MX, TXT, SRV, CNAME, etc.).
  - `vscale_domain_tag`: Tag domains.
  - `vscale_ptr_record`: Manage PTR records for server IPs.
  - `vscale_backup`: Trigger server backups.
- **Data Sources**:
  - `vscale_account`: Fetch account profile and balance.
  - `vscale_images`: List available OS images.
  - `vscale_locations`: List datacenters.
  - `vscale_rplans`: List tariff plans.
  - `vscale_prices`: Read resources pricing details.
  - `vscale_backups`: List existing backups.

---

## Configuration

Configure the provider using your VScale API token:

```hcl
terraform {
  required_providers {
    vscale = {
      source  = "vscale/vscale"
      version = "1.0.0"
    }
  }
}

provider "vscale" {
  token = "your-api-token"
}
```

Alternatively, you can export the token as an environment variable:
```bash
export VSCALE_TOKEN="your-api-token"
```

---

## Local Development

### Prerequisites
- Go 1.21+
- Terraform 1.0+ or OpenTofu 1.6+

### Build Provider
To compile the provider binary locally:
```bash
go build -o terraform-provider-vscale
```

### Dev Overrides Setup
Create a CLI config file (e.g. `tofurc` or `terraformrc`):
```hcl
provider_installation {
  dev_overrides {
    "vscale/vscale" = "/absolute/path/to/provider/directory"
  }
  direct {}
}
```
Set the environment variable pointing to this config before running plan/apply:
```bash
# Unix
export TF_CLI_CONFIG_FILE="/path/to/tofurc"

# Windows PowerShell
$env:TF_CLI_CONFIG_FILE="C:\path\to\tofurc"
```

---

## Registry Publication

### 1. Terraform Registry
1. Push this repository to GitHub under the name `terraform-provider-vscale`.
2. Register in [Terraform Registry](https://registry.terraform.io) using your GitHub account.
3. Import the provider.
4. Push a new Git tag (e.g. `v1.0.0`) to trigger the release workflow. The binaries will be automatically compiled, signed, and published.

### 2. OpenTofu Registry
1. Publish your provider on GitHub and push the release tag.
2. Submit a Pull Request to [github.com/opentofu/registry](https://github.com/opentofu/registry) adding your metadata JSON config (e.g., `providers/v/vscale.json` with repository link).

---

## License

This project is licensed under the Mozilla Public License 2.0 - see the [LICENSE](LICENSE) file for details.
