terraform {
  required_providers {
    vscale = {
      source  = "vscale/vscale"
      version = "1.0.0"
    }
  }
}

provider "vscale" {
  # token = "your-api-token"
}

# --- Data Sources Example ---
data "vscale_account" "current" {}

data "vscale_images" "all" {}

data "vscale_locations" "all" {}

data "vscale_rplans" "all" {}

data "vscale_prices" "current" {}

data "vscale_backups" "all" {}

# Output account email and balance
output "account_email" {
  value = data.vscale_account.current.email
}

output "account_balance" {
  value = data.vscale_account.current.balance
}

# --- Resources Example ---

# SSH Key
resource "vscale_ssh_key" "my_key" {
  name = "terraform-key"
  key  = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQ..."
}

# DNS Domain
resource "vscale_domain" "my_domain" {
  name = "example-vscale.com"
}

# DNS record
resource "vscale_domain_record" "www" {
  domain_id = vscale_domain.my_domain.id
  name      = "www"
  type      = "A"
  ttl       = 300
  content   = "1.2.3.4"
}

# Domain Tag
resource "vscale_domain_tag" "tag" {
  name    = "production"
  domains = [vscale_domain.my_domain.name]
}

# Scalet Server
resource "vscale_scalet" "my_server" {
  name      = "my-vscale-server"
  make_from = "ubuntu_20.04_64_001_master"
  rplan     = "medium"
  location  = "spb0"
  do_start  = true

  keys = [
    vscale_ssh_key.my_key.id
  ]
}

# Reverse DNS record (PTR)
resource "vscale_ptr_record" "reverse" {
  ip      = vscale_scalet.my_server.public_address.address
  content = "example-vscale.com"
}

# Server Backup
resource "vscale_backup" "my_backup" {
  scalet_id = vscale_scalet.my_server.id
  name      = "manual-backup-before-deploy"
}
