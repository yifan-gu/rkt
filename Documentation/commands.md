# rkt Commands

Work in progress.
Please contribute if you see an area that needs more detail.

## Downloading Images (ACIs)

[aci-images]: https://github.com/appc/spec/blob/master/spec/aci.md#app-container-image
[appc-discovery]: https://github.com/appc/spec/blob/master/spec/discovery.md#app-container-image-discovery

rkt runs applications packaged according to the open-source [App Container Image][aci-images] specification.
ACIs consist of the root filesystem of the application container, a manifest, and an optional signature.

ACIs are named with a URL-like structure.
This naming scheme allows for a decentralized discovery of ACIs, related signatures and public keys.
rkt uses these hints to execute [meta discovery][appc-discovery].

* [trust](subcommands/trust.md)
* [fetch](subcommands/fetch.md)

## Running Pods

[metadata-spec]: https://github.com/appc/spec/blob/master/spec/ace.md#app-container-metadata-service
[rkt-mds]: subcommands/metadata-service.md

rkt can execute ACIs identified by name, hash, local file path, or URL.
If an ACI hasn't been cached on disk, rkt will attempt to find and download it.
To use rkt's [metadata service][metadata-spec], enable registration with the `--mds-register` flag when [invoking it][rkt-mds].

* [run](subcommands/run.md)
* [stop](subcommands/stop.md)
* [enter](subcommands/enter.md)
* [prepare](subcommands/prepare.md)
* [run-prepared](subcommands/run-prepared.md)

## Pod inspection and management

rkt provides subcommands to list, get status, and clean its pods.

* [list](subcommands/list.md)
* [status](subcommands/status.md)
* [gc](subcommands/gc.md)
* [rm](subcommands/rm.md)
* [cat-manifest](subcommands/cat-manifest.md)

## Interacting with the local image store

rkt provides subcommands to list, inspect and export images in its local store.

* [image](subcommands/image.md)

## Metadata Service

The metadata service helps running apps introspect their execution environment and assert their pod identity.

* [metadata-service](subcommands/metadata-service.md)

## API Service

The API service allows clients to list and inspect pods and images running under rkt.

* [api-service](subcommands/api-service.md)

## Misc

* [version](subcommands/version.md)
* [config](subcommands/config.md)

## Global Options

In addition to the flags used by individual `rkt` commands, `rkt` has a set of global options that are applicable to all commands.

| Flag | Default | Options | Description |
| --- | --- | --- | --- |
| `--cpuprofile (hidden flag)` | `` | A file path | Write CPU profile to the file |
| `--debug` |  `false` | `true` or `false` | Prints out more debug information to `stderr` |
| `--dir` | `/var/lib/rkt` | A directory path | Path to the `rkt` data directory |
| `--insecure-options` |  none | **none**: All security features are enabled<br/>**http**: Allow HTTP connections. Be warned that this will send any credentials as clear text.<br/>**image**: Disables verifying image signatures<br/>**tls**: Accept any certificate from the server and any host name in that certificate<br/>**ondisk**: Disables verifying the integrity of the on-disk, rendered image before running. This significantly speeds up start time.<br/>**pubkey**: Allow fetching pubkeys via insecure connections (via HTTP connections or from servers with unverified certificates). This slightly extends the meaning of the `--trust-keys-from-https` flag.<br/>**all**: Disables all security checks | Comma-separated list of security features to disable |
| `--local-config` |  `/etc/rkt` | A directory path | Path to the local configuration directory |
| `--memprofile (hidden flag)` | `` | A file path | Write memory profile to the file |
| `--system-config` |  `/usr/lib/rkt` | A directory path | Path to the system configuration directory |
| `--trust-keys-from-https` |  `false` | `true` or `false` | Automatically trust gpg keys fetched from HTTPS (or HTTP if the insecure `pubkey` option is also specified) |
| `--user-config` |  `` | A directory path | Path to the user configuration directory |

## Logging

By default, rkt will send logs directly to stdout/stderr, allowing them to be captured by the invoking process.
On host systems running systemd, rkt will attempt to integrate with journald on the host.
In this case, the logs can be accessed directly via journalctl.

#### Accessing logs via journalctl

To read the logs of a running pod, get the pod's machine name from `machinectl`:

```
$ machinectl
MACHINE                                  CLASS     SERVICE
rkt-bc3c1451-2e81-45c6-aeb0-807db44e31b4 container rkt

1 machines listed.
```

or `rkt list --full`

```
$ rkt list --full
UUID                                  APP    IMAGE NAME                              IMAGE ID             STATE    CREATED                             STARTED                             NETWORKS
bc3c1451-2e81-45c6-aeb0-807db44e31b4  etcd   coreos.com/etcd:v2.3.4                  sha512-7f05a10f6d2c  running  2016-05-18 10:07:35.312 +0200 CEST  2016-05-18 10:07:35.405 +0200 CEST  default:ip4=172.16.28.83
                                      redis  registry-1.docker.io/library/redis:3.2  sha512-6eaaf936bc76
```

The pod's machine name will be the pod's UUID prefixed with `rkt-`.
Given this machine name, logs can be retrieved by `journalctl`:

```
$ journalctl -M rkt-bc3c1451-2e81-45c6-aeb0-807db44e31b4
[...]
```

To get logs from one app in the pod:

```
$ journalctl -M rkt-bc3c1451-2e81-45c6-aeb0-807db44e31b4 -t etcd
[...]
$ journalctl -M rkt-bc3c1451-2e81-45c6-aeb0-807db44e31b4 -t redis
[...]
```

Additionaly, logs can be programmatically accessed via the [sd-journal API](https://www.freedesktop.org/software/systemd/man/sd-journal.html).

##### Stopped pod

To read the logs of a stopped pod, use:

```
journalctl -m _MACHINE_ID=132f9d560e3f4d1eba8668efd488bb62

[...]
```

On some distributions such as Ubuntu, persistent journal storage is not enabled by default. In this case, it is not possible to get the logs of a stopped pod. Persistent journal storage can be enabled with `sudo mkdir /var/log/journal` before starting the pods.
