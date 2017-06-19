# Network Container Analyzer (NetCAN)

![img/example.png](img/example.png)

## Run in a Docker container

```
docker run --privileged --net=host -v /:/rootfs:ro kiratech/netcan -rootfs=/rootfs
```

## Run standalone

```
sudo netcan
```

You will find netcan using your browser at: http://<machine-ip>:8000

## Project status

NetCAN is a newborn project with the purpose of helping in analyzing networking interfaces and links happening trough local or distributed bridges.
The main goal of NetCAN are:
- providing a standalone web interface that shows the current network configuration inside the host or the cluster (as indipendent as possible respect to the cluster manager or container engine)
- providing a set of tools to analyze network traffic happening inside containers by applying filters and capturing packets
- providing a way to export netcan datas to external tools like Elastic

### Checklist
- [x] Simple visualizer (first version, not the definitive one)
- [ ] Support for multiple nodes (cluster mode)
- [ ] Discussion on what to analyze to write the analysis tools (TBD)
- [ ] Definitive version of the visualizer


