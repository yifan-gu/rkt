# rkt list

You can list all rkt pods.

```
$ rkt list
UUID		APP	    IMAGE NAME					                STATE	CREATED		    STARTED		    FINISHED	 NETWORKS
046b5bde	redis	registry-1.docker.io/library/redis:latest	running	11 seconds ago	10 seconds ago			     default:ip4=172.16.28.104
    		etcd	coreos.com/etcd:v2.0.9											
ae3981ac	nginx	registry-1.docker.io/library/nginx:latest	exited	1 minute ago	1 minute ago	1 minute ago	
```

You can view the full UUID as well as the image's ID by using the `--full` flag.

```
$ rkt list --full
UUID		                			APP	    IMAGE NAME					                IMAGE ID		    STATE	CREATED			                    STARTED					            FINISHED				           NETWORKS
046b5bde-6068-4c88-b541-5e8d2f2ad3ba	redis	registry-1.docker.io/library/redis:latest	sha512-9ae1fe1e18b2	running	2016-04-11 16:47:04.909 -0700 PDT	2016-04-11 16:47:05.129 -0700 PDT						               default:ip4=172.16.28.104
					                    etcd	coreos.com/etcd:v2.0.9				        sha512-91e98d7f1679					
ae3981ac-1b86-4d92-a1dd-1426e97d1726	nginx	registry-1.docker.io/library/nginx:latest	sha512-4e932fa575d6	exited	2016-04-11 16:45:19.665 -0700 PDT	2016-04-11 16:45:19.809 -0700 PDT	2016-04-11 16:45:19.809 -0700 PDT
```

## Options

| Flag | Default | Options | Description |
| --- | --- | --- | --- |
| `--full` |  `false` | `true` or `false` | Use long output format |
| `--no-legend` |  `false` | `true` or `false` | Suppress a legend with the list |

## Global options

See the table with [global options in general commands documentation](../commands.md#global-options).
