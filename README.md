# PERMIT-SSH

Another CLI tool written in Go.  \
**Permit** stores users and their SSH pub keys to help automate the task of autorising SSH connection in servers.

## Usage

``` bash
permit [options]
```

Usage examples:

``` bash
# Only with a key
permit -key RSA... -ip root@XX.XX.XX.XX

# Only with a saved user
permit -user example -ip root@XX.XX.XX.XX

# Both
permit -key RSA... -user example -ip root@XX.XX.XX.XX
```

## Options

``` text
 -help, --help
  Prints help information.

 -user, --user
  User to add.

 -key, --key
  SSH key to add.
 
 -ip, --ip
  IP address of the machine.

 -list, --list
  List stored users.

 -del, --del
  Delete saved user in DB.

   // TODO //

 -s, --search
  Search for user.
 
   // TODO //
```
