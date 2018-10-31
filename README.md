# Ghosts

<img src="https://raw.githubusercontent.com/lobre/ghosts/master/static/logo.png" width="250">

> Host entries generator and web interface for Docker containers.

This Go program will listen for Docker events and generate `/etc/hosts` entries according to specific labels declared on Docker containers.

On top of that, a friendly web interface will list web exposed containers on to a nice grid.

See the web interface after having created the following containers.

    docker run -d --name test1 \
        --label ghosts.host="test1.local" \
        nginx

    docker run -d --name test2 \
        --label ghosts.category="System" \
        --label ghosts.host="test2.local" \
        nginx

    docker run -d --name test3 \
        --label ghosts.host="test3.local" \
        --label ghosts.name="Friendly app" \
        nginx

    docker run -d --name test4 \
        --label ghosts.host="test4.local" \
        --label ghosts.name="Jenkins" \
        --label ghosts.logo="https://wiki.jenkins.io/download/attachments/2916393/logo.png" \
        nginx

![screenshot](https://raw.githubusercontent.com/lobre/ghosts/master/screenshot.png)

## Quickstart with Docker

To start the server, use the following command.

    docker run --rm --name ghosts -v /var/run/docker.sock:/var/run/docker.sock -v /etc/hosts:/app/hosts -p 8080:8080 lobre/ghosts

## Binary parameters

    Usage of ./ghosts:
      -addr string
            Web app address and port (default ":8080")
      -help string
            Change the Web help link (default "https://github.com/lobre/ghosts/blob/master/README.md")
      -hosts string
            Custom location for hosts file
      -nohelp
            Disable help on web interface
      -nohosts
            Don't generate hosts file
      -noweb
            Don't start web server
      -proxyip string
            Specific proxy IP for hosts entries (default "127.0.0.1")
      -proxymode
            Enable proxy
      -traefikmode
            Enable integration with Traefik proxy

## Container parameters

 - `ghosts.host`: Host of container (e.g. mycontainer.local.com). If in traefik mode, it can be taken from `traefik.frontend.rule`.
 - `ghosts.port`: Override port. Otherwise taken from exposed ports or traefik port.
 - `ghosts.name`: Define web name. Otherwise taken from the container name.
 - `ghosts.proto`: Define web protocol. If in traefik mode, it can be taken from `traefik.frontend.entryPoints`.
 - `ghosts.auth`: Define if auth protected entry. If in traefik mode, it can be taken from `traefik.frontend.auth.basic`.
 - `ghosts.category`: Define a web category. Defaults to "Apps".
 - `ghosts.logo`: Define a web logo. Defaults to a generated avatar with the initials of the entry name.
 - `ghosts.description`: Add a web description that will appear as a tooltip.
 - `ghosts.noweb`: Don't show on the web.
 - `ghosts.nohosts`: Don't generate entry in hosts file.
 - `ghosts.direct`: Use direct container IP in hosts file even if in proxy mode.
 - `ghosts.webdirect`: Use direct container IP directly in web view even if in proxy mode.
