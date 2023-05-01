# machina

`machina` is a tool to create Virtual Machines based on yaml templates

## Installation

```bash
curl -sL https://raw.githubusercontent.com/enkodr/machina/main/bin/install | sh -
```

## Configuration

### libvirt

With `libvirt`, mounting is only possible on connection to `qemu:///system` and the following options needs to be added to the `/etc/libvirt/qemu.conf` file:

```toml
user = "my_username"
group = "my_group"
security_driver = "none"
```

## Templates

You can use an already defined [template](https://https://github.com/enkodr/machina/tree/main/templates) or create your own.

```yaml
name: ubuntu

variant: "ubuntu22.04"

image : 
  url: "https://cloud-images.ubuntu.com/jammy/20230428/jammy-server-cloudimg-amd64.img"
  checksum: "sha256:3e1898e9a0cc39df7d9c6af518584d53a647dabfbba6595d1a09dd82cabe8a87"
  
credentials:
  username: machina
  password: machina
  groups:
  - "users"
  - "admin"
  
specs:
  cpus: 2
  memory: "2G"
  disk: "50G"

scripts:
  install: |
    #!/bin/bash
    sudo apt update && sudo apt -y upgrade
  
  init: |
    #!/bin/bash
    echo "Welcome to your Virtual Machine!"
mount:
  name: workspace
  hostPath: "/home/user/workspace"
  guestPath: "/home/machina/workspace"
```

### Override a template

You can also create a template that uses a base template.
Any value set on this template will override the values from the extended template.

```yaml
name: my-vm

extends: ubuntu

specs: 
  cpus: 4
  memory: "8G"
  disk: "100G"
```

**Note:** Although a template can only extend a single template, is possible to create a
hierarchy of templates by extending a template that already extends another template.
