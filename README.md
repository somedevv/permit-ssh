# PERMIT-SSH

Another CLI tool written in Go.  \
**Permit** stores users and their SSH pub keys to help automate the task of autorising SSH connection in servers.

## Usage

``` bash
permit [options]
```

Usage example:

``` bash
permit -key RSA... -user example -ip root@XX.XX.XX.XX
```

## Options

``` text
 -h, --help
  Prints help information.

 -u, --user
  User to add.

 -k, --key
  SSH key to add.

 -l, --list
  List stored users.

   // TODO //

 -d, --delete
  Delete user.

 -s, --search
  Search for user.
 
 -ip, --ip
  IP address of the machine.

   // TODO //
```
