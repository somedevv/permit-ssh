# Permit

Another CLI tool written in Go.  \
**Permit** stores users and their SSH pub keys to help automate the task of autorising SSH connection in servers.

## Usage

``` bash
permit [options]
```

Usage examples:

``` bash
# Only with a key
permit -k RSA... -ip root@XX.XX.XX.XX

# Only with a saved user
permit -u example -ip root@XX.XX.XX.XX

# Both
permit -k RSA... -u example -ip root@XX.XX.XX.XX
```

This tool depends, for now, in your ssh config to authenticate into servers.

## Flags

``` text
  --version           Displays the program version string.

  -h --help           Displays help with available flag, subcommand, and
                      positional valueparameters.

  -del --delete       Delete a user or key. If IP is set, the user will be
                      deleted from the server, otherwise, the user will be
                      deleted from the database

  -u --user           The user to add or delete

  -k --key            The key to add or delete

  -ip --address       The IP of the server to add or delete the user

  -l --list           List all the users in the database

  -i --interactive    Interactive mode

   // TODO //

  -s, --search        Search for user.
 
   // TODO //
```
