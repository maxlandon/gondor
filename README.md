
Gondor - Go Maltego Transform Framework
=====

The Gondor frameworks aims to port the useful [Canari Framework](https://github.com/redcanari/canari3), 
written in Python 3, and integrate it with Paterva's own functionality set, like Transform servers 
usage/registration/management. Currently the state of Maltego Transforms libraries in Go is quite meager,
and this project aims to solve this. Another major goal is to streamline even more the production and use 
of Transforms, whether as independent executables or within larger projects & codebases.

## Main Features

This library aims to provide a few advantages over the Canari Framework, 
in part through a smart use of the Go toolchain and module management system.

- The framework goes with a command line tool, `gondor`, which is used to create,
  manage and build transforms, as well as to produce the various configurations and
  profiles that are required by Maltego when adding Transforms and Entities.

- The `gondor` command-line tool is partly a wrapper around the Go toolchain (which 
  is needed to build the transforms). **The tool should ship with its own Go toolchain**,
  and take advantage of the Go environement vars and its module system, so as **to work
  within a well-defined, hermetic build & setup environment** (see [Sliver](https://github
  .com/BishopFox/sliver) for a reference usage of an embedded Go toolchain). In the same
  vein, the **framework should allow a user to create, write and use Transforms/Entities
  from within arbitrary filesystem locations**, if used as a library in your program.

- The framework should allow to **produce standalone transform binaries, or be used 
  as a library**. For example, you might have a Go program producing entities among
  many other things, and you can't simply restart it each time you want to run a
  transform: this is where you want your program to be able to either send data to
  Maltego asynchronously, as well as make Maltego to send requests to it.

- To the best extent possible, the framework should **allow writers to create Transforms
  manipulating Go-native types**. This is achieved, in part, by the use of custom XML
  marshalling implementations (sounds scary, but don't be fooled: Go is a reliable friend).

- The framework should offer the capacity to create, manage and use more advanced Maltego
  functionality, such as [Machines](https://docs.maltego.com/support/solutions/articles/
  15000019249-machines-transform-macros-), or even the [Transform Distribution Servers 
  (iDTS)](https://docs.maltego.com/support/solutions/articles/15000020198-what-is-itds-#
  the-public-tds-0-0). From the command line tool, users can **start a Transform server, setup
  its various details** (API tokens, certs, etc), and **easily load arbitrary Transform Sets**.
  This process should be even more automated than the current Maltego Web UI setup process.


## Documentation

Currently there is no completion support for the command line tool, but all commands support
a `-h` or `--help` option, which will print their subcommands, options and various descriptions.
In the future this project will ship a complete documentation on how to write Transforms with it,
as well as some reference workflow examples, from a blank repo all the way to a running Transform server. 


## Support
 
- All contributions and ideas are welcome.
- All open issues will be solved as quickly as possible.
- Feature requests will be studied on a case-by-case basis.


## License - GPLv3

Gondor is licensed under [GPLv3](https://www.gnu.org/licenses/gpl-3.0.en.html).
