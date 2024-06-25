# wakwaw

### Prerequisites

- [OpenTofu](https://opentofu.org/docs/intro/install/)
- [Ansible](https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html)
- [DigitalOcean Personal Access Token](https://docs.digitalocean.com/reference/api/create-personal-access-token/)
  - Permissions: `create`, `read`, `update`, and `delete` for droplet, firewall, ssh_key, and vpc.
- [Cloudflare Global API Key](https://developers.cloudflare.com/fundamentals/api/get-started/keys/#view-your-global-api-key)

### Setup

1. Create a `terraform.tfvars` file in the `tofu` directory.

   - Required: `do_token` (DigitalOcean Personal Access Token) & `region` (choose the region closest to you).

2. Copy the example environment file and update it:

   ```sh
   cp swarm/.env.example swarm/.env
   ```

   - Required: `DOMAIN` (your domain), `ACME_EMAIL` (your email), `CF_API_EMAIL` (Cloudflare Global API email), and `CF_API_KEY` (Cloudflare Global API key).
   - To generate an encrypted password, run:

     ```sh
     htpasswd -nb -B USERNAME PASSWORD
     ```

     - The output will be in the format: `USERNAME:ENCRYPTED_PASSWORD`.

3. Apply the Terraform configuration:

```sh
tofu -chdir=tofu apply
```

4. Set the correct permissions for your SSH key:

```sh
chmod 600 id_rsa
```

5. Run the Ansible playbook:

   ```sh
   ANSIBLE_CONFIG=ansible/ansible.cfg ansible-playbook ansible/main.yml -i hosts --private-key id_rsa
   ```

   - If you encounter the error `Error connecting: Error while fetching server API version: Not supported URL scheme http+docker`, run:

   ```sh
   ansible-galaxy collection install community.docker --upgrade
   ```

### Update

1. Apply the Terraform configuration:

```sh
tofu -chdir=tofu apply
```

2. Run the Ansible playbook:

   ```sh
   ANSIBLE_CONFIG=ansible/ansible.cfg ansible-playbook ansible/main.yml -i hosts --private-key id_rsa
   ```

#### If you only want to deploy

```sh
ANSIBLE_CONFIG=ansible/ansible.cfg ansible-playbook ansible/deploy.yml -i hosts --private-key id_rsa -e "stack_name=NAME"
```

### Destroy

1. Destroy the Terraform-managed infrastructure:

   ```sh
   tofu -chdir=tofu destroy
   ```

### Available services

| Name       | URL                                   | UI  | Auth                                                                 |
| ---------- | ------------------------------------- | --- | -------------------------------------------------------------------- |
| traefik    | https://traefik.example.com           | yes | `username: TRAEFIK_USERNAME`<br />`password: TRAEFIK_PASSWORD`       |
| portainer  | https://wakwaw.example.com/portainer  | yes | `username: admin`<br />`password: PORTAINER_PASSWORD`                |
| minio      | https://wakwaw.example.com/minio/     | yes | `username: MINIO_USERNAME`<br />`password: MINIO_PASSWORD`           |
| grafana    | https://wakwaw.example.com/grafana    | yes | `username: GRAFANA_USERNAME`<br />`password: GRAFANA_PASSWORD`       |
| prometheus | https://wakwaw.example.com/prometheus | yes | `username: PROMETHEUS_USERNAME`<br />`password: PROMETHEUS_PASSWORD` |
