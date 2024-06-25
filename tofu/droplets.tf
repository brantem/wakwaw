resource "digitalocean_droplet" "manager" {
  count = var.managers

  image             = "docker-20-04"
  name              = "${var.name}-manager-${count.index + 1}"
  region            = var.region
  size              = var.manager_size
  monitoring        = true
  vpc_uuid          = digitalocean_vpc.default.id
  ssh_keys          = [digitalocean_ssh_key.default.id]
  graceful_shutdown = true
}

resource "digitalocean_droplet" "worker" {
  count = var.workers

  image             = "docker-20-04"
  name              = "${var.name}-worker-${count.index + 1}"
  region            = var.region
  size              = var.worker_size
  monitoring        = true
  vpc_uuid          = digitalocean_vpc.default.id
  ssh_keys          = [digitalocean_ssh_key.default.id]
  graceful_shutdown = true
}

resource "local_file" "hosts" {
  filename = "../hosts"
  content = templatefile("${path.module}/templates/hosts.tmpl",
    {
      manager = digitalocean_droplet.manager
      worker  = digitalocean_droplet.worker
    }
  )
}
