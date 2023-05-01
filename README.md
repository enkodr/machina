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