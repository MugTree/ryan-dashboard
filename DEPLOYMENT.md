# Deployment

Idea here was to create something relatively simple taking advantage of Go's facility of building binaries for
other platforms. Another nice feature is that with smaller sites all the assests can be baked in using go's embed directive.

Basically build the binary for the remote host - AMD linux - and use shell scripts to manage it. Start / stop etc. It just then sits behind a reverse proxy. The database is sqlite, possibly would have been better to use postgres as could use a visual client like dataGrip to admin it over the network - definitely simpler if you are working with others.

## Local setup

Uses a variety of scripts to get stuff in place

```cmd
    ./scripts/upload_site.sh
    ./scripts/www_stop.sh
    ./scripts/www_stop.sh
    ./scripts/rescue_db.sh
    ./scripts/restore_db.sh
```

Those scripts require that the sudoers file is setup correctly.

### Setting the user so that we don't need sudo

When we run any GitHub actions or shell sessions we don't want to have to deal with a prompt.

**!! on the current box the user is portable simply use `sudo visudo /etc/sudoers.d/portable` and add the new services to the end of the existing line**

Otherwise...

**Create a create a user file in**
`/etc/sudoers.d/youruser`

**And then run something like this command as root**

```
    cat > /etc/sudoers.d/youruser << SUDO
    youruser ALL=(ALL)  NOPASSWD: /bin/systemctl start someservice.service
    SUDO
```

Verify the permissions

```
    sudo chmod 440 /etc/sudoers.d/youruser
```

And you can then run

```
    sudo systemctl start youruser.service
```

From a local shell script or whatever....

### Using sqlite remotely

Have to use the command line on this but, trade off is not having to run a DB server.

If I want to update the database - I use a script to pull or rescue the database (rescue_db.sh) which pulls a copy down
I then run the changes against that and then use the restore script to replace the modified file.

The rescue script creates a backup of the database first so not all is lost if you make a mistake.

## Server setup

Assumes ubuntu linux or some such ...

Idea here is to run a binary using systemd to manage it.

We have a simple unit file that on the ubuntu distro lives here

```
    /etc/systemd/system
```

Here's an example file 'www.notzero.co.uk.service'

```
    [Unit]
    Description=www.notzero.co.uk
    After=network.target

    [Service]
    Environment="IS_PRODUCTION=true"
    WorkingDirectory=/home/portable/notzero.co.uk/www
    ExecStart=/home/portable/notzero.co.uk/www/bin/notzero.www.amd64
    User=portable
    Restart=always
```

All fairly self explanatory

To work with this file

```
     systemctl start www.notzero.co.uk.service
     systemctl stop www.notzero.co.uk.service
```

Apart from IS_PRODUCTION the ENV variables will live in a .env file next to the exe.

To read any errors on start up

```
    journalctl -u www.notzero.co.uk.service
```

## Caddy

Use Caddy as a reverse proxy and to generate the ssl certs for free.
https://caddyserver.com/

Here's a typical entry

```
    somesite.co.uk {
        redir https://www.somesite.co.uk{uri}
    }

    www.somesite.co.uk {
        reverse_proxy localhost:10000

        file_server

        encode gzip
        log {
                output file /var/log/caddy/www.somesite.log
                format json
        }
    }
```

The service file lives at

```
    /lib/systemd/system/somesite.service
```

Useful commands

```
    find / -name Caddyfile
```

To get the location of the Caddyfile

```
    systemctl reload caddy`
```

Any time you make a change to the config

### General point about services

Good idea to set it to start at boot...

```
    sudo systemctl enable application.service
```

#### Then the basics...

```
    sudo systemctl start somesite.service
    sudo systemctl stop somesite.service
    sudo systemctl restart somesite.service
    sudo systemctl status somesite.service
```

Need ssh, http https (caddy), rysnc (copy files) and postgres

## Firewall - Using UFW

Below are the obvious rules that are needed

### ssh

To login and maintain the box

`sudo ufw allow OpenSSH

### Caddy

The front end webserver that I've chosen to use

`sudo ufw allow proto tcp from any to any port 80,443

### Rysnc

Need this for efficiently copying files

`sudo ufw allow 873/tcp

### Postgres

This as with other applications can be tunnelled through ssl
