resource "digitalocean_vpc" "default" {
  name     = "${var.name}-vpc"
  region   = var.region
  ip_range = var.cidr_block
}

resource "digitalocean_firewall" "node" {
  name = "${var.name}-node"
  droplet_ids = concat(
    digitalocean_droplet.manager[*].id,
    digitalocean_droplet.worker[*].id,
  )

  inbound_rule {
    protocol         = "tcp"
    port_range       = "4789"
    source_addresses = [var.cidr_block]
  }

  inbound_rule {
    protocol         = "tcp"
    port_range       = "7946"
    source_addresses = [var.cidr_block]
  }

  inbound_rule {
    protocol         = "udp"
    port_range       = "7946"
    source_addresses = [var.cidr_block]
  }

  inbound_rule {
    protocol         = "tcp"
    port_range       = "9323"
    source_addresses = [var.cidr_block]
  }

  outbound_rule {
    protocol              = "tcp"
    port_range            = "1-65535"
    destination_addresses = ["0.0.0.0/0", "::/0"]
  }

  outbound_rule {
    protocol              = "udp"
    port_range            = "1-65535"
    destination_addresses = ["0.0.0.0/0", "::/0"]
  }

  outbound_rule {
    protocol              = "icmp"
    destination_addresses = ["0.0.0.0/0", "::/0"]
  }
}

resource "digitalocean_firewall" "manager" {
  name        = "${var.name}-manager"
  droplet_ids = digitalocean_droplet.manager[*].id

  inbound_rule {
    protocol         = "tcp"
    port_range       = "22"
    source_addresses = ["0.0.0.0/0", "::/0"]
  }

  inbound_rule {
    protocol         = "tcp"
    port_range       = "80"
    source_addresses = ["0.0.0.0/0", "::/0"]
  }

  inbound_rule {
    protocol         = "tcp"
    port_range       = "443"
    source_addresses = ["0.0.0.0/0", "::/0"]
  }

  inbound_rule {
    protocol         = "udp"
    port_range       = "443"
    source_addresses = ["0.0.0.0/0", "::/0"]
  }

  inbound_rule {
    protocol         = "tcp"
    port_range       = "2377"
    source_addresses = [var.cidr_block]
  }

  inbound_rule {
    protocol         = "tcp"
    port_range       = "5000"
    source_addresses = [var.cidr_block]
  }
}


resource "digitalocean_firewall" "worker" {
  name        = "${var.name}-worker"
  droplet_ids = digitalocean_droplet.worker[*].id

  inbound_rule {
    protocol         = "tcp"
    port_range       = "22"
    source_addresses = [var.cidr_block]
  }
}
