resource "tls_private_key" "rsa_4096" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "digitalocean_ssh_key" "default" {
  name       = var.name
  public_key = tls_private_key.rsa_4096.public_key_openssh
}

resource "local_file" "id_rsa" {
  content  = tls_private_key.rsa_4096.private_key_pem
  filename = "../id_rsa"
}
