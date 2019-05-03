# How to run this:

## `1.` Setup the haproxy with 2 servers:
```sh
cd ../testing
go build

# Open 2 extra terminal tab and run:
./testing -name serv-1 -addr :8080
./testing -name serv-2 -addr :8081
# Go back to the first terminal tab

./runDock.sh
``` 

## `2.` Run this example:
```
go build
./example
```
