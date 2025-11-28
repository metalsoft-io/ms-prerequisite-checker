# MetalSoft Prerequisites Check

Stand-alone tool to check access required to install and operate MetalSoft global and site controllers.

## Checks

### Global Controller installation

This test is performed with command `global-install`

Arguments:

* `ms-repo` - URL of the repo with MetalSoft images (default to `http://repo.metalsoft.io`)
* `ms-repo-secure` - Secure URL of the repo with MetalSoft images (default to `https://repo.metalsoft.io`)
* `ms-registry` - URL of the MetalSoft registry (defaults to `https://registry.metalsoft.dev`)

Checks the following:

* HTTP on port 80 to `ms-repo`
* HTTPS on port 443 to `ms-repo`
* HTTPS on port 443 to `ms-registry`
* ICMP to `1.1.1.1`
* HTTP on port 80 to `1.1.1.1`
* HTTPS on port 443 to `1.1.1.1`
* HTTPS on port 443 to <https://downloads.dell.com/>
* HTTP on port 80 to <http://downloads.linux.hpe.com/>
* HTTPS on port 443 to <https://quay.io/>
* HTTPS on port 443 to <https://gcr.io/>
* HTTPS on port 443 to <https://cloud.google.com/>
* HTTPS on port 443 to <https://helm.traefik.io/>
* HTTPS on port 443 to <https://k8s.io/>
* TCP on port 587 to `smtp.office365.com`

### Global Controller operation

This test is performed with command `global-operate`

Checks the following:

* HTTPS on port 443 to <registry.metalsoft.dev>
* HTTP on port 80 to <repo.metalsoft.io>
* HTTPS on port 443 to <repo.metalsoft.io>

### Site Controller installation

This test is performed with command `site-install`

Checks the following:

* HTTPS on port 443 to <registry.metalsoft.dev>
* HTTP on port 80 to <repo.metalsoft.io>
* HTTPS on port 443 to <repo.metalsoft.io>

### Site Controller operation

This test is performed with command `site-operate`

Arguments:

* `global-controller-hostname` - IP address or hostname of the global controller.
* `nfs-server` (optional) - NFS server for use by the site controller.

Checks the following:

* HTTP on port 80 to `global-controller-hostname`
* HTTPS on port 443 to `global-controller-hostname`
* TCP on port 9003 to `global-controller-hostname`
* TCP on port 9009 to `global-controller-hostname`
* TCP on port 9010 to `global-controller-hostname`
* TCP on port 9011 to `global-controller-hostname`
* TCP on port 9090 to `global-controller-hostname`
* TCP on port 9091 to `global-controller-hostname`
* UDP on port 53 to `global-controller-hostname`
* TCP on port 111 to `nfs-server` - performed if the optional argument is provided
* UDP on port 111 to `nfs-server` - performed if the optional argument is provided
* TCP on port 2049 to `nfs-server` - performed if the optional argument is provided
* UDP on port 2049 to `nfs-server` - performed if the optional argument is provided

If the global controller is not installed and operational run the mock services on the node that will host it.
The mock service listens on the following ports and protocols:

* HTTP on port 80
* HTTPS on port 443
* TCP on port 9003
* TCP on port 9009
* TCP on port 9010
* TCP on port 9011
* TCP on port 9090
* TCP on port 9091
* UDP on port 53

Start the service with the following command:

```bash
ms-prerequisite-check global-service
```

Arguments:

* `listen-ip` (optional) - IP address to listen on. By default listens on all interfaces.

### Switch connectivity

This test is performed with command `site-manage-switch`

Arguments:

* `nos` - The switch NOS - one of (OS10, SONiC, JunOS, Cisco).
* `management-ip` - IP address of the switch management port.
* `username` - Username of the switch management admin user.
* `password` - Password of the switch management admin user.

Checks the following:

* HTTP connection to `management-ip` on port 80
* HTTPS connection to `management-ip` on port 443
* SSH connection to `management-ip` on port 22 using the provided `username` and `password`
* NETCONF - SSH connection to `management-ip` on port 830 using the provided `username` and `password` - performed when the `nos` is "JunOS"

### Server connectivity

This test is performed with command `site-manage-server`

Arguments:

* `vendor` - The server vendor - one of (Dell, HP, Lenovo).
* `bmc-ip` - IP address of the server BMC interface.
* `username` - Username of the server BMC admin user.
* `password` - Password of the server BMC admin user.
* `iso-link` (optional) - Link to an ISO to test mounting virtual media.

Checks the following:

* Redfish - HTTPS connection to `bmc-ip` on port 443
* SSH - SSH connection to `bmc-ip` on port 22 using the provided `username` and `password`
* IPMI - UDP connection to `bmc-ip` on port 623
* VNC - HTTP connection to `bmc-ip` on port 5901 - performed when the `vendor` is "Dell" and the `vnc-password` is provided

### Site Controller inbound connections

To test reachability from the servers and switches to the site controller use the `site-service` command.
In this mode the tool will listen on the specified IP address (or all interfaces if omitted) for the following inbound requests:

* DHCPv4 on port 53 - will print the summary of the received packet without responding
  * NOTE: This function is implemented for Linux systems only and requires elevated permissions!

## Building

### TLS Certificate

The TLS certificate used for the HTTPS service are embedded in the tool.
To generate new certificates use the following command:

```bash
cd ./certs/certs
go run "/usr/local/go/src/crypto/tls/generate_cert.go" --rsa-bits=2048 --host=localhost
```

```bash
cd .\certs\certs
go run "C:\Program Files\Go\src\crypto\tls\generate_cert.go" --rsa-bits=2048 --host=localhost
```

### Linux

```bash
GOOS='linux' GOARCH='amd64' CGO_ENABLED='0' go build -ldflags "-s -X main.version=7.0.0" -o ./bin/linux_amd64/ms-prerequisite-check ./cmd/cli
```

### Windows PowerShell

```bash
$Env:GOOS='windows' ; $Env:GOARCH='amd64' ; go build -ldflags "-s -X main.version=7.0.0" -o ./bin/windows_amd64/ms-prerequisite-check.exe ./cmd/cli
```

## Running

### Startup parameters

| Parameter               | Type   | Default value          | Description                      |
| ----------------------- | ------ | ---------------------- | -------------------------------- |
| -log-level              | string | info                   | Log level: debug,info,warn,error |

## Examples

### Prerequisites for installing the global controller

```bash
ms-prerequisite-check -log-level=debug global-install
```

### Prerequisites for running the global controller

```bash
ms-prerequisite-check -log-level=debug global-operate
```

### Prerequisites for installing the site controller

```bash
ms-prerequisite-check -log-level=debug site-install
```

### Prerequisites for running the site controller

Run the mock services on the global controller node.

```bash
ms-prerequisite-check -log-level=debug global-service
```

Optional arguments:

* `listen-ip` - IP address on which to listen for incoming requests

Test the connectivity by running the tool on the site controller node.
The `global-controller-hostname` argument points to the global controller node.

```bash
ms-prerequisite-check -log-level=debug site-operate global-controller-hostname=metal.acme.com
```

Optional arguments:

* `nfs-server` - points to the NFS server for the site controller storage

### Test connectivity to managed switch

```bash
ms-prerequisite-check -log-level=debug site-manage-switch nos=SONiC management-ip=1.2.3.4 username=admin password=secret
```

### Test connectivity to managed server

```bash
ms-prerequisite-check -log-level=debug site-manage-server vendor=Dell bmc-ip=1.1.1.1 username=root password=calvin
```

Optional arguments:

* `iso-link` - location of an ISO image to test mounting virtual media

### Site controller mock service

Run the mock services on the site controller node.

```bash
ms-prerequisite-check -log-level=debug site-service
```

Optional arguments:

* `listen-ip` - IP address on which to listen for incoming requests
