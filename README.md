# MetalSoft Prerequisites Check

Stand-alone tool to check access required to install and operate MetalSoft global and site controllers.

Requirements: <https://metalsoft.atlassian.net/browse/MS-4809>

## Building

### TLS Certificate

The TLS certificate used for the HTTPS service are embedded in the tool.
To generate new certificates use the following command:

```bash
cd ./certs/certs
go run /usr/local/go/src/crypto/tls/generate_cert.go" --rsa-bits=2048 --host=localhost
```

```bash
cd .\certs\certs
go run "C:\Program Files\Go\src\crypto\tls\generate_cert.go" --rsa-bits=2048 --host=localhost
```

### Linux

```bash
GOOS='linux' GOARCH='amd64' go build -ldflags "-s -X main.version=6.3.0" -o ./bin/linux_amd64/ms-prerequisite-check ./cmd/cli
```

### Windows PowerShell

```bash
$Env:GOOS='windows' ; $Env:GOARCH='amd64' ; go build -ldflags "-s -X main.version=6.3.0" -o ./bin/windows_amd64/ms-prerequisite-check.exe ./cmd/cli
```

## Running

### Startup parameters

| Parameter               | Type   | Environment variable     | Default value          | Description                                  |
| ----------------------- | ------ | ------------------------ | ---------------------- | -------------------------------------------- |
| -log-level              | string | LOG_LEVEL                | info                   | Log level: trace,debug,info,warn,error,fatal |

## Examples

### Prerequisites for installing the global controller

```bash
ms-prerequisite-check global-install k8s-repo=https://repo.metalsoft.io
```

Optional arguments:

* `ms-repo` - URL of the repo with MetalSoft images (default to `http://repo.metalsoft.io`)
* `ms-repo-secure` - Secure URL of the repo with MetalSoft images (default to `https://repo.metalsoft.io`)
* `ms-registry` - URL of the MetalSoft registry (defaults to `https://registry.metalsoft.dev`)

### Prerequisites for running the global controller

```bash
ms-prerequisite-check global-operate
```

### Prerequisites for installing the site controller

```bash
ms-prerequisite-check site-install
```

### Prerequisites for running the site controller

Run the mock services on the global controller node.

```bash
ms-prerequisite-check global-service
```

Optional arguments:

* `listen-ip` - IP address on which to listen for incoming requests

Test the connectivity by running the tool on the site controller node.
The `global-controller-hostname` argument points to the global controller node.

```bash
ms-prerequisite-check site-operate global-controller-hostname=metal.acme.com
```

Optional arguments:

* `nfs-server` - points to the NFS server for the site controller storage

### Test connectivity to managed switch

```bash
ms-prerequisite-check site-manage-switch nos=SONiC management-ip=1.2.3.4 username=admin password=secret
```

### Test connectivity to managed server

```bash
ms-prerequisite-check site-manage-server vendor=Dell bmc-ip=1.2.3.4 username=root password=calvin
```

Optional arguments:

* `iso-link` - location of an ISO image to test mounting virtual media
