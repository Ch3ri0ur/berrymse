# Helper Scripts
These are a set of convenience scripts, to register the autostart.

Run as sudo. e.g.  `sudo ./register.sh`

`berrymse.service` expects the executable to be at `/home/pi/berrymse/armv7l/berrymse` and that port 80 and `dev/video0` are to be used. If this is not the case modify the file accordingly before registering.

## register.sh

This is the only script necessary. It adds the service file to the services folder and restarts the daemon. 

This new service is then enabled (added to autostart).
## Unregister

Undoes Register

|Script|Function|
|-----|--|
|register.sh| Registers Service and adds autostart|
|unregister.sh| Undoes register|
|disable.sh| Removes autostart |
|stop.sh| Stops the currently running service (does not change autostart) |
|start.sh| Starts the service |
|restart.sh| Restarts the service |



