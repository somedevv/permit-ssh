# Permit

Another CLI tool written in Go.  \
**Permit** stores users and their SSH pub keys to help automate the task of autorising SSH connection in servers.

## Install

This tool depends, for now, in your ssh config to authenticate into servers. For AWS support you must have the ```aws cli``` installed and configured. \
Rename and place the executable on ```$HOME/.local/bin/permit``` and create the folder ```$HOME/.local/bin/.permit_data```. \
Then, add it to your path adding ```export PATH="$PATH:$HOME/.local/bin"``` to your ```.*rc``` file.

## Usage

``` text
permit - Your own SSH key manager and friend, made by somedevv.

  Usage:

    // WITH IP
    permit [add|remove] -user x -key "x" -ip xx.xx.xx.xx

    // WITH AWS
    permit [add|remove] -user x -key "x" aws --instance x --region x --profile x

  Subcommands:
    add             Subcommand 'add' is a command that adds a key to a server or a user.
    remove          Subcommand 'remove' is a command that deletes an existing user/key from a server or from saved users.
    help            Prints help
    list            Lists all saved users and keys
    interactive     Subcommand 'interactive' puts the tool in interactive mode.

  General Flags:
        --version       Displays the program version string.
    -h  --help          Displays help with available flag, subcommand, and positional value parameters.
    -ip --address       IP address of the server.
    -u  --user          Name of the user.
    -k  --key           Public key of the user.

  AWS Specific Flags:
    -i  --instance      Instance ID of the server.
    -r  --region        AWS region of the server.
    -p  --profile       AWS profile to use.
```

## Config (beta)

Create a file ```$HOME/.local/bin/.permit_data/config.json``` with the following content:

``` yaml
db type: local
```

## To Do

- [x] Configuration file
- [ ] Installer
- [ ] Full featured Interactive mode
- [x] AWS Integration for EC2 instances
- [ ] Support for external DB
