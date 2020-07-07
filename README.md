## UINETD - A Udp forward tool just like

### Description
Redirects TCP&UDP connections from a local address & port to another address.
Configure files is just in /etc/uintd.conf.(Like rinetd).
The process will open a socket to listen all the port from local bindings.

### Usage

local address|local bind port|remote address|remote port|`tcp` or `udp` or `all`
-|-|-|-|-

`0.0.0.0 80 192.168.1.2 80`

