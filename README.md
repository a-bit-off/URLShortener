# URLShortener

## Description
This is a rest API service that shortens links, you can run it on your server and use.

## Download
Write in your terminal `git clone https://github.com/a-bit-off/URLShortener`.

## Build
To start a project you need to be in main dir `"URLShortener"`, and write `make` in your terminal.

## Usage
You started the service locally with the `make` command. <br />
In the `local.yml` configuration file, you can see on which host you are running it (`8082`). The URL looks like this: `"http://localhost:8082"`. <br />
You have the ability to add, delete and get a link by alias. <br />

<br />

Add: <br />
In your terminal: <br />
`curl -X POST -H "Content-Type: application/json" -d '{
     "url": "your url",
     "alias": "your alias"
}' http://localhost:8082/url` <br />

<br />

Get: <br />
`"http://localhost:8082/" + "your alias"` <br />

<br />
 
Delete URL: <br />
`"http://localhost:8082/url/" + "your alias"` <br />

<br />
Or you may use Postman. <br />

## Stack
config - `cleanenv`. <br />
logger - `slog`. <br />
storage - `sqlite`. <br />
router - `chi`, `middleware`. <br />
build - `Makefile`. <br />
test - `httptest`, `testing`, `mockery`. <br />
