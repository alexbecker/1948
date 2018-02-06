# Nineteen-Fourty-Eight (1948)

#### A minimalist web framework

1948 is designed to help you make a simple and performant website without unnecessary cruft.
It consists of:
 * a static templating system,
 * an extensible, lightweight webserver,
 * a set of management scripts,
 * and a set of hooks for plugin integration.

## Prerequisits

Currently, 1948 only supports Unix-like systems. Using 1948 requires:
 * git
 * Python 3
 * pip
 * Go
 * make

## Setup

 - Clone the repository by running `git clone https://github.com/alexbecker/1948`.
 - *(Recommended)* Create a Python 3 [virtualenv](https://virtualenv.pypa.io/en/stable/) in the repository.
 - Install Jinja2 by running `pip install jinja2`.
 - Run `./init_project.sh` to create a new project.

## Customization

The `init_project.sh` script creates a skeleton project for you in the `local/` directory.
The resulting files can be customized to create your site:
 * The `local/env.sh` script is used to export configuration variables for various scripts and for the webserver.
 * The static files of your site are built from `local/templates/` using `local/template_conf.py` for any configuration.
 * The server exposes hooks for extension in `local/go/src/local/`, and can be modified to add dynamic behavior to your site.
 * The Makefile at `local/Makefile` exposes a `local` target that is made whenever the project is made.
Since this is made after the static files, the `STATIC_DEPS` variable is available and targets can be appended to it to make before static files are made.
 * The `local/build.sh` script can be extended to add additional steps to the build process.
 * The `local/deploy.sh` script can be extended to add additional steps to the deploy process.
 * The `local/server_side/` directory will be copied to the server when the project is deployed.
 * The `local/server_side/install.sh` script can be extended to add additional steps to the server-side installation.

## Building and Deploying

To build the entire site, run `./build`. Alternatively, individual targets found in `Makefile` can be built directly.

To deploy the site, run `source local/env.sh; ./deploy.sh` (or provide other environment variables if desired).

To install server-side dependencies, on the server run `source env.sh; ./install.sh`
(note that this is a different install script, located in `server_side/`).

## Plugins

There is not yet a plugin repository or any type of plugin management system.
However, there are hooks built into the webserver, build system and various scripts to allow creating plugins.

By convention, a plugin is installed in the `local/plugins/` directory,
and contains an `install.sh` script that should be run on first setup.
If the plugin requires installation on the 
Plugins may require additional configuration, or server-side installation,
conventially via a script at `plugins/<plugin_dir>/server_side/install.sh`.

## Tips

 * If you want to maintain multiple projects, make `local/` a symlink.
 * If you want to deploy to multiple environments, e.g. a dev or staging environment, create additonal `local/env_*.sh` files.

## TODOs

 * Plugin management infrastructure.
 * Story for handling upgrades, since install scripts are intended for fresh installations.
