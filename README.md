# Ghosts

![logo](https://raw.githubusercontent.com/lobre/ghosts/master/static/logo.png | width=50)

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
