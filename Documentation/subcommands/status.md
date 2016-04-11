# rkt status

Given a pod UUID, you can get the exit status of its apps.
Note that the apps are prefixed by `app-`.

```
$ rkt status 046b5bde
state=exited
created=2016-04-11 16:47:04.909 -0700 PDT
started=2016-04-11 16:47:05.129 -0700 PDT
finished=2016-04-11 16:48:35.746 -0700 PDT
pid=30380
exited=true
app-redis=0
app-etcd=0
```

If the pod is still running, you can wait for it to finish and then get the status with `rkt status --wait UUID`

## Options

| Flag | Default | Options | Description |
| --- | --- | --- | --- |
| `--wait` |  `false` | `true` or `false` | Toggle waiting for the pod to exit |

## Global options

See the table with [global options in general commands documentation](../commands.md#global-options).
