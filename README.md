<p align="center">
    <a href="https://youtu.be/PifMSY8-PO4" target="_blank">
      <img src="/assets/gopher2.png" width="180" />
    </a>
    <h3 align="center">Poodle</h3>
    <p align="center">A fast and beautiful command line tool to build API requests</p>
    <p align="center">
        <a href="https://github.com/Clivern/Poodle/actions"><img src="https://github.com/Clivern/Poodle/workflows/Build/badge.svg"></a>
        <a href="https://github.com/Clivern/Poodle/actions"><img src="https://github.com/Clivern/Poodle/workflows/Release/badge.svg"></a>
        <a href="https://github.com/Clivern/Poodle/releases"><img src="https://img.shields.io/badge/Version-2.2.3-red.svg"></a>
        <a href="https://goreportcard.com/report/github.com/Clivern/Poodle"><img src="https://goreportcard.com/badge/github.com/Clivern/Poodle?v=2.2.3"></a>
        <a href="https://github.com/Clivern/Poodle/blob/master/LICENSE"><img src="https://img.shields.io/badge/LICENSE-MIT-orange.svg"></a>
    </p>
</p>

<p align="center">
    <img src="/assets/poodle_showcase.gif" width="85%" />
</p>
<h4 align="center">
     <a href="https://youtu.be/PifMSY8-PO4" target="_blank">:unicorn: Check out the full Demo!</a>
</h4>
<br/>

Poodle is an interactive command line tool to build and test web APIs based on a pre-built definitions.

`Poodle` has the following features:

- Register your web services and endpoints easily.
- Use variables in endpoints definitions.
- Painless debugging and interaction with APIs.
- Search web services and endpoints interactively.
- Edit services and endpoints easily ([config is just a TOML file](/misc/service_definition.toml)).
- Sync services via Gist automatically.

## Documentation

Download [the latest poodle binary](https://github.com/Clivern/Poodle/releases). Also install [fzf](https://github.com/junegunn/fzf) for better searching otherwise poodle will use a built-in one. Make it executable from everywhere.

```zsh
$ curl -sL https://github.com/Clivern/Poodle/releases/download/x.x.x/poodle_x.x.x_OS.tar.gz | tar xz
```

To list all commands and options

```zsh
$ poodle help

A fast and beautiful command line tool to build API requests

If you have any suggestions, bug reports, or annoyances please report
them to our issue tracker at <https://github.com/clivern/poodle/issues>

Usage:
  poodle [command]

Available Commands:
  call        Interact with one of the configured services
  configure   Configure Poodle
  delete      Delete a service definition file
  edit        Edit service definition file
  help        Help about any command
  license     Print the license
  new         Creates a new service definition file
  sync        Sync services definitions
  version     Print the version number

Flags:
  -c, --config string   config file (default "/Users/Clivern/poodle/config.toml")
  -h, --help            help for poodle
  -v, --verbose         verbose output

Use "poodle [command] --help" for more information about a command.
```

To configure poodle, You will need to provide your github username and oauth token with a `gist` scope if you need the backup/sync feature

```zsh
$ poodle configure
```

To sync definitions with backend. for now only github gists supported

```zsh
$ poodle sync
```

To create a new service.

```zsh
$ poodle new
```

by default we use `https://httpbin.org` as service API for testing so change with your web service API.

To edit a previously created service file:

```zsh
$ poodle edit
```

To start calling your services endpoints:

```zsh
$ poodle call
```

To start calling a custom service endpoints without storing it globally:

```zsh
# Assuming that .poodle.toml exists on current dir
$ poodle call -f ./.poodle.toml
```

To delete a service definition file:

```zsh
$ poodle delete
```

For more info how to create an API service definition, [please check this file.](/misc/service_definition.toml)


## Upgrade Notes

To upgrade from version `1.x.x` to version `2.x.x`, It is recommended to delete the remote github gist and poodle will sync changes in a new proper way.

## Versioning

For transparency into our release cycle and in striving to maintain backward compatibility, Poodle is maintained under the [Semantic Versioning guidelines](https://semver.org/) and release process is predictable and business-friendly.

See the [Releases section of our GitHub project](https://github.com/clivern/poodle/releases) for changelogs for each release version of Poodle. It contains summaries of the most noteworthy changes made in each release.


## Bug tracker

If you have any suggestions, bug reports, or annoyances please report them to our issue tracker at https://github.com/clivern/poodle/issues


## Security Issues

If you discover a security vulnerability within Poodle, please send an email to [hello@clivern.com](mailto:hello@clivern.com)


## Contributing

We are an open source, community-driven project so please feel free to join us. see the [contributing guidelines](CONTRIBUTING.md) for more details.


## License

© 2020, Clivern. Released under [MIT License](https://opensource.org/licenses/mit-license.php).

**Poodle** is authored and maintained by [@clivern](http://github.com/clivern).
