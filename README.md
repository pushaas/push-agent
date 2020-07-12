# push-agent

This is documented [here](https://github.com/pushaas/pushaas-docs#component-push-agent).

## running locally

Requires [push-redis](https://github.com/pushaas/push-redis) and [push-stream](https://github.com/pushaas/push-stream) to be running.

```shell
make run
```

## publishing images

```shell
make docker-push TAG=<tag>
```

---
