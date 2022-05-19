# Permit

Another CLI tool written in Go.  \
**Permit** stores users and their SSH pub keys to help automate the task of autorising SSH connection in servers.

## Usage

``` text
permit - Your own SSH key manager and friend, made by somedevv.

  Usage:
    permit [add|remove] -user x -key "x" -ip xx.xx.xx.xx

  Subcommands:
    add             Subcommand 'add' is a command that adds a key to a server or a user.
    remove          Subcommand 'remove' is a command that deletes an existing user/key from a server or from saved users.
    help            Prints help
    list            Lists all saved users and keys
    interactive     Subcommand 'interactive' puts the tool in interactive mode.

  Flags:
       --version        Displays the program version string.
    -h --help           Displays help with available flag, subcommand, and positional value parameters.
```

This tool depends, for now, in your ssh config to authenticate into servers. \
You need to place the executable on ```$HOME/.local/bin``` and create a folder ```$HOME/.local/bin/.permit_data```

## To Do

- [ ] Configuration file
- [ ] Installer
- [ ] Full featured Interactive mode
- [ ] AWS Integration for EC2 instances
- [ ] Suport for external DB
