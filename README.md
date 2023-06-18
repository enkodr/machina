# machina

This project is a command line application written in Go that allows users to 
create and configure virtual machines using YAML files. 
The application supports VM creation using either `libvirt` or `qemu`.

## Installation

```bash
curl -sL https://raw.githubusercontent.com/enkodr/machina/main/bin/install | sh -
```

## How to Use

Once you've installed the application, you can use it to create, manage, 
and delete virtual machines.

### Available Commands

* `copy` - Copies files from the host to the VM and vice versa.
* `create` - Creates a new virtual machine based on a YAML configuration file.
* `delete` - Deletes an existing virtual machine.
* `health` - Shows if all the dependencies are installed.
* `list` - Lists all existing virtual machines.
* `shell` - Enters a VM shell.
* `start` - Starts an existing virtual machine.
* `stop` - Stops a running virtual machine.
* `template` - Lists the available templates or download one if a name is specified.

### Examples

**Creating a new virtual machine from the ubuntu template:**

```bash
machina create ubuntu
```

**Creating a new virtual machine from an existing file:**

```bash
machina create -f template.yaml
```

**Starting an existing virtual machine:**

```bash
machina start my_vm
```

**Sopping a running virtual machine:**

```bash
machina stop my_vm
```

**Copying files from the host to the VM:**

```bash
machina copy /path/to/host my_vm:/path/to/guest
```

**Copying files from the VM to the host:**

```bash
machina copy my_vm:/path/to/guest /path/to/host 
```

**Deleting an existing virtual machine:**

```bash
machina delete my_vm
```

**Checking the health of the system:**

```bash
machina health
```

**Listing all existing virtual machines:**

```bash
machina shell my_vm
```

**Listing available templates:**

```bash
machina template
```

**Downloading a template:**

```bash
machina template template_name
```

## Configuration

When the first machine is created, a file with the configuration is created
on the `~/.config/machina/config.yaml`

## Working with Templates

The tool provides the capability to use pre-configured templates for creating 
virtual machines. Templates make it easy to create VMs with standardized settings.

### Machine Name (name)
The `name` key is used to specify the name of the virtual machine. 
This name should be unique within the system.

### OS Variant (variant)
The `variant` key is used to set the operating system variant to be used 
for the virtual machine. This is typically required for the libvirt hypervisor. 
You can obtain a list of available variants for your system by running the command 
`virt-install --os-variant list`.

### Image Configuration (image)
The `image` key is used to specify the image to be used for provisioning the virtual machine. 
It includes the URL of the image file and its checksum for verification.

### User Credentials (credentials)
The `credentials` key is used to define the user credentials for the default user of the virtual machine. 
It includes the username, password, and optional group memberships.

### Hardware Resources (resources)
The `resources` key is used to specify the hardware resources allocated to the virtual machine. 
It includes the number of CPU cores (**cpus**), amount of memory (**memory**), and disk size (**disk**).

### Scripts (scripts)
The `scripts` key is used to define  the scripts that will be executed inside the virtual machine. 
It includes an `install` script, which is executed during machine installation, 
and an `init` script, which is invoked by the .bashrc shell.

### Mount Point (mount)
The `mount` key is used to define mount points from the host machine into the virtual machine. 
Each mount point includes a name, `hostPath` (path on the host), and `guestPath` (path inside the virtual machine). 
This configuration uses the `virtio-9p` driver for mounting.

#### Example

Templates are defined in YAML files similar to VM configurations. 
Below is an example of a template configuration:

```yaml
# The name to use for the machine.
# This needs to be a unique name in the system.
name: ubuntu

# This value sets the OS variant o use. This is only needed for `libvirt` hypervisor
# To grab a list of the ones available for your system you can just run 
# `virt-install --os-variant list`
variant: "ubuntu22.04"

# The image to be used to provision the machine
image : 
  url: "https://cloud-images.ubuntu.com/jammy/20230428/jammy-server-cloudimg-amd64.img"
  checksum: "sha256:3e1898e9a0cc39df7d9c6af518584d53a647dabfbba6595d1a09dd82cabe8a87"
  
# The user credentials to be set for the default user.
# This user will have root access without asking for password.
credentials:
  username: machina
  password: machina
  groups:
  - "users"
  - "admin"
  
# This option specifies the hardware you want to set for the machine
resources:
  # Sets the number of cores to set for the machine
  cpus: 2
  # Sets the ammount of RAM to set for the machine. You must use the 
  # standard G or M units, for Gigabyte and Megabyte respectivly.
  memory: "2G"
  # Sets the ammount of space to define for the VM virtual disk.
  disk: "50G"

# The scripts are executed inside of the virtual machine.
# For compatibility with different distro's the scripts are not inherited
scripts:
  # The `install` script is execute during the machine installation
  install: |
    #!/bin/bash

  # The `init` script is invoced by the `.bashrc` shell
  init: |
    #!/bin/bash
    echo "Welcome to your virtual machine"

# Mounts defines a set of mount points from the host into the VM, where the
# `hostPath` is the path in the host and the `guestPath` sets the path inside the VM 
# where the mount will be defined. 
# This will use `virtio-9p` driver.
mount:
  # name: share
  # hostPath: "/path/to/host/folder"
  # guestPath: "/path/to/guest/folder"
```

To create a new template, save your configuration in a file with the .yaml 
extension in the templates directory. The name of the file will be the name 
of the template.

### Override a template

You can also create a template that uses a base template.
Any value set on this template will override the values from the extended template.

```yaml
name: my_vm

extends: ubuntu

resources: 
  cpus: 4
  memory: "8G"
  disk: "100G"
```

## Cluster
The `Cluster` type in represents a logical grouping of machines or nodes that work together
to achieve a common goal. It serves as a fundamental building block for managing and 
orchestrating distributed systems. 

The Cluster configuration allows you to define various aspects of the cluster, 
including its name, parameters, output, and the list of machines or instances it comprises.

### Params (params)
The `Params` field allows you to define a set of parameters specific to the cluster. 
These parameters can be used to customize the behavior or settings of the cluster components. 
Each parameter consists of a name-value pair, allowing you to configure the cluster to suit your requirements.

### Outputs (outputs)
The `Output` field enables you to specify the desired output or results associated with the cluster. 
This can include any information or artifacts generated by the cluster during its operation. 
The output field provides a structured way to capture and store relevant data produced by the 
cluster for further analysis or consumption.

### Instances (instances)
Furthermore, the Instances field allows you to define a list of machines or nodes that form part of the cluster. 
Each machine in the list represents a separate entity `Machine` with its own set of properties, as defined above. 
This allows you to tailor the behavior and characteristics of each machine within the cluster to achieve the desired system behavior.

### Example

```yaml
# The kind Cluster allows the creation of a cluster of Machines
kind: Cluster 

# The name for the cluster. This needs to be unique in the system.
name: kubernetes

# Params allows you to set parametes to be used in the scripts.
# These parameters will be available for all the machines within the workspace.
params:
  hello: "World"

# The outputs allows you to define the outputs from the cluster.
# These outputs will be available for all the machines within the workspace.
outputs:
- message

# The machines uses a list of kind Machine 
machines:
- kind: Machine
  name: server
  
  extends: ubuntu
  # Number of machines to create
  replicas: 1

  scripts:
    install: |
      #!/bin/bash
      echo "Hello {{ .Params.hello }}"
      echo "Hello from server" > $(outputs.message)


- kind: Machine
  name: client
  extends: ubuntu
  replicas: 1

  scripts:
    install: |
      #!/bin/bash  
      echo < $(inputs.message)
  
```