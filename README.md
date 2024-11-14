# monitorConnections
Monitor's incoming and outgoing initiating connections for both TCP and UDP.

###Parameters
```
-device <device id e.g. eth0>
-exclude-public (hides connections from and to public IP addresses)
-exclude-udp (hides UDP traffic)
```
###Building
```
go build monitorConnections.go
```


###CLI Example
```
./monitorConnections -device en1 -exclude-public -exclude-udp
```
