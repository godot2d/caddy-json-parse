[![Go](https://github.com/godot2d/caddy-json-parse/workflows/Go/badge.svg)](https://github.com/godot2d/caddy-json-parse/actions)

# caddy-json-parse
Caddy v2 module for parsing json request body.

## Installation

```
xcaddy build v2.0.0 \
    --with github.com/godot2d/caddy-json-parse
```


## Usage

`json_parse` parses the request body as json for reference as [placeholders](https://caddyserver.com/docs/caddyfile/concepts#placeholders).

### Caddyfile

Simply use the directive anywhere in a route. If set, `strict` responds with bad request if the request body is an invalid json.
```
json_parse [<strict>]
```

And reference variables via `{json.*}` placeholders. Where `*` can get as deep as possible. e.g. `{json.items.0.label}`

#### Special Feature: ServerName Extraction from Attach Field

This plugin includes special handling for extracting `ServerName` from nested JSON strings in the `attach` field. When your JSON request body contains an `attach` field with a JSON string that includes a `ServerName` property, you can access it directly using:

```
{json.attach.ServerName}
```

For example, with this request body:
```json
{
  "game": "com.arrow.defense3d",
  "orderId": "order-id",
  "uid": "test",
  "amount": "100",
  "tradeState": "SUCCESS",
  "timestamp": 123456789,
  "platform": "ios",
  "sku": "whatever",
  "attach": "{ \"roleId\": 100016, \"gameServerId\": 1, \"shopId\": 1, \"merchandiseId\": 1001 , \"ServerName\" : \"dev\"}"
}
```

You can access `"dev"` using `{json.attach.ServerName}`.


#### Examples

**Example 1:** Run a [command](https://github.com/abiosoft/caddy-exec) only if the github webhook is a push on master branch.
```
@webhook {
    expression {json.ref}.endsWith('/master')
}
route {
    json_parse # enable json parser
    exec @webhook git pull origin master
}
```

**Example 2:** Extract ServerName from attach field and respond with it.
```
:8080 {
    json_parse
    respond "ServerName: {json.attach.ServerName}, Game: {json.game}, OrderId: {json.orderId}"
}
```

### JSON

`json_parse` can be part of any route as an handler

```jsonc
{
  ...
  "routes": [
    {
      "handle": [
        {
          "handler": "json_parse",

          // if set to true, returns bad request for invalid json
          "strict": false 
        },
        ...
      ]
    },
  ...
  ]
}
```

## License

Apache 2
