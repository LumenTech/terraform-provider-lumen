## Lumen API key

Generate access token for authentication with Lumen API endpoint. Users should have `morph-api` access. Execute below CLI with your Lumen _$username_ and _$password_ to generate access token for Lumen API.

```shell
$ curl -k -X POST 'https://api.lumen.com/oauth/token?grant_type=password&scope=write&client_id=morph-api' -d 'username=$username' -d 'password=$password'
```

It should generate a response like this:
```shell
{"access_token":"0000-0000-0000-0000","token_type":"bearer","refresh_token":"0000000000000","expires_in":000000000,"scope":"write"}
```
