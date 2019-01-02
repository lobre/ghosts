# Ghosts

<img src="https://raw.githubusercontent.com/lobre/ghosts/master/static/logo.png" width="250">

> Host entries generator, automatic network connection and web interface for Docker containers.

This Go program will listen for Docker events and fill the gaps of:

 - Generating `/etc/hosts` entries according to specific labels declared on Docker containers.
 - Generating a web interface to list web exposed containers in a nice grid.
 - Auto connect network of web exposed containers to an external proxy.

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

![screenshot](https://raw.githubusercontent.com/lobre/ghosts/master/img/screenshot.png)

## Quickstart with Docker

To start the server, use the following command.

    docker run --rm --name ghosts -v /var/run/docker.sock:/var/run/docker.sock -v /etc/hosts:/app/hosts -p 8080:8080 lobre/ghosts

### Windows

To let the container edit the `C:\Windows\System32\drivers\etc\hosts`, we need to update the permissions of the file to let the current user edit without Administrator rights.

![screenshot](https://raw.githubusercontent.com/lobre/ghosts/master/img/windows_permissions.png)

Then, we also need to add the `hostsforcewindowsstyle` parameter of ghosts to be sure we use Windows style end-of-line characters.

    docker run --rm --name ghosts -v /var/run/docker.sock:/var/run/docker.sock -v /c/etc/hosts:/app/hosts -p 8080:8080 lobre/ghosts -hostsforcewindowsstyle

#### Docker machine Virtualbox

By default, only `C:\Users` is shared to the VM. So the hosts file won't be available to Docker containers by default. We need to add a specific shared mount.

![screenshot](https://raw.githubusercontent.com/lobre/ghosts/master/img/vbox_shared.png)

## Modes

Ghosts has two different modes.

 - **Direct mode (default)**: containers are accessed using their direct IP.
 - **Proxy mode**: containers are accessed through the IP of a predefined proxy (the proxy IP can be defined using the binary parameter `-proxyIP`).

## Binary parameters

    Usage of ./ghosts:
      -addr string
            Web app address and port (default ":8080")
      -help string
            Change the Web help link (default "https://github.com/lobre/ghosts/blob/master/README.md")
      -hosts string
            Custom location for hosts file
      -hostsforcewindowsstyle
            Force CRLF end of lines and one entry per line when generating hosts entries
      -nohelp
            Disable help on web interface
      -nohosts
            Don't generate hosts file
      -noweb
            Don't start web server
      -proxycontainername string
            Name of proxy container
      -proxyip string
            Specific proxy IP for hosts entries (default "127.0.0.1")
      -proxymode
            Enable proxy
      -proxynetautoconnect
            Enable automatic network connection between proxy and containers
      -webnavbgcolor
            Color of navbar on the web interface (default "#f1f1fc")
      -webnavtextcolor
            Color of the navbar text on the web interface (default "#50596c")

## Container parameters as labels

 - `ghosts.host`: Comma separated list of hosts for container (e.g. mycontainer.local.com).
 - `ghosts.path`: Comma separated list of paths for container (e.g. /my-path).
 - `ghosts.port`: Override internal exposed port.
 - `ghosts.proto`: "http" or "https" (default "http").
 - `ghosts.name`: Define web name. Otherwise taken from the container name.
 - `ghosts.auth`: Define if auth protected entry. If true, a lock will be displayed on the web interface.
 - `ghosts.category`: Define a web category. Defaults to "Apps". Supports multiple values (comma separated list).
 - `ghosts.logo`: Define a web logo. Defaults to a generated avatar with the initials of the entry name.
 - `ghosts.description`: Add a web description that will appear as a tooltip.
 - `ghosts.noweb`: Don't show on the web.
 - `ghosts.nohosts`: Don't generate entry in hosts file.
 - `ghosts.nonetautoconnect`: Don't connect to proxy network.
 - `ghosts.direct`: Use direct container IP in hosts file even if in proxy mode.
 - `ghosts.webdirect`: Use direct container IP directly in web view even if in proxy mode.

### Segments

You can define multiple sets of urls/port using segments. They can be defined using the following labels structure.

 - `ghosts.<my_segment_name>.host`
 - `ghosts.<my_segment_name>.path`
 - `ghosts.<my_segment_name>.port`
 - `ghosts.<my_segment_name>.proto`

The name of the segment will be shown on the web interface. This feature can be useful if your container has multiple vhosts (e.g. frontend and backend).
