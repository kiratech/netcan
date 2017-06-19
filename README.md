# Network Container Analyzer (NetCAn)

![img/example.png](img/example.png)

## Run in a Docker container

```
docker run --privileged --net=host -v /:/rootfs:ro fntlnz/netcan -rootfs=/rootfs
```
