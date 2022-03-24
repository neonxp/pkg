# pkg

Selfhosted golang packages proxy.

Live version: https://go.neonxp.dev/

Package page: https://go.neonxp.dev/pkg/ (browser will redirects to documentation, but `go get go.neonxp.dev/pkg` will install package)

## Configuration

`config.json` file:

```json
{
    "title": "WEBSITE TITLE",
    "host": "go.neonxp.dev",
    "packages": {
        "pkg": {
            "pkg": "pkg",
            "vcs": "git",
            "repo": "https://github.com/neonxp/pkg",
            "desc": "Package description"
        },
        "jsonrpc2": {
            "pkg": "jsonrpc2",
            "vcs": "git",
            "repo": "https://github.com/neonxp/jsonrpc2",
            "desc": "Имплементация сервера JSON-RPC 2.0 с генериками."
        },
        "test": {
            "pkg": "test",
            "vcs": "git",
            "repo": "https://github.com/example/test",
            "desc": "Test description."
        },
        ...
    }
}

```

## Author

Alexander Kiryukhin <i@neonxp.dev>

## License

![GPL v3](https://www.gnu.org/graphics/gplv3-with-text-136x68.png)
