# serial-gateway

A serial gateway that runs in the cluster. It will read a serial port mounted from the node and process the messages received

## Virtrual Serial

Taken from this [stackoverflow post](https://stackoverflow.com/questions/22568878/emulate-serial-port)

To test the side car, in one terminal run the following to create a virtual serial device.

```sh
socat PTY,link=./virtual-tty,raw,echo=0 -
```

Then you can run this in a different terminal to read from the device.  Now you can type input into the first terminal, hit enter to make it available in the side car program.

```sh
SERIAL_PORT=./virtual-tty go run cmd/serial/main.go
```
