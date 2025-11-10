# Sensors project

based on

Need to return the chart element and the chart script separaltey

https://github.com/MugTree/www-template

## www template

Basic web template for creating go sites. Has login, contacts, auth, compression, logging basic routing. Decent templates. Embedded fs so a single binary.

## Dependencies

- templ - go templating language - see the make file and https://templ.guide/
- air - live reload for go apps https://github.com/air-verse/air
- goose - db migrations - https://github.com/pressly/goose
- sqlite - good starting point - see scripts folder for rescue and restore scripts
- datastar - reactivity events etc - https://data-star.dev/
- hurl - https://hurl.dev/ - bit like curl but easier to use

## To get up and running - clone folder and create a repo

```bash
    cp -R /Users/me/Developer/go-projects/www-template /Users/me/Developer/go-projects/some-site
    cd /Users/me/Developer/go-projects/some-site
```

Remove and .git directories. Create a new one...

Install any comand line requirements

```bash
    go install github.com/pressly/goose/v3/cmd/goose@latest
    go install github.com/a-h/templ/cmd/templ@latest
    go install github.com/air-verse/air@latest
```

## To set up the .env

required to get the goose variables set in env for migrations

```bash
    ./load_env.sh
```

## To run the website

```bash
    make start-dev
```

- This starts up the site on http://localhost:8080
- Plus starts up a live reload on http://localhost:7331 for ease of dev when adding web features

### Tests

Tests are minimal at this point but I've added a test for the code in the www package that starts and stops the server

```bash
cd www
go test
```

If you need to check the return values of the JSON functions of the website Hurl is a good option. See below for a one liner that calls the
http://localhost:8080/api/get-sites using the JSON api token

```bash
hurl --variable api_key=some-api-key-q3we get_sites_from_api.hurl
```
