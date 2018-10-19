# Ghosts

![alt text](https://raw.githubusercontent.com/lobre/ghosts/master/img/logo.png)

Helper web interface for hosts declared in Docker containers

## Behavior

In order for a entry to appear, we need a correct host.
That means a protocol, a host (or ip), a port and a path.

Override rules are as follows.

### protocol

 1. label ghosts.protocol
 2. traefik entrypoint
 3. http default value

### host

 1. label ghosts.host or traefik frontend rule (according to traefik_mode)
 2. if direct_mode: container ip

## Parameters

    proxy_ip
    traefik_mode
    direct_mode (using container ip)
    hosts_file

## Container parameters

 - logo
 - category
 - name
 - description
 - auth
 - proto
 - host
 - path
 - hide
