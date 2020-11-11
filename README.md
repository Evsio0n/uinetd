## uinetd - A Udp forward tool just like rinetd

### Description
Redirects TCP&UDP connections from a local address & port to another address.
Configuration files is just in /etc/uintd.conf.(Like rinetd).
The process will open a socket to listen all the port from local bindings.

### Usage

local address|local bind port|remote address|remote port|`tcp` or `udp` or `all`
-|-|-|-|-
0.0.0.0|80|192.168.1.2|80|ALL

### Filter Connections
`allow *.*.*.*` or `deny *.*.*.*` (IPv6 available)

### Log who is using your redirection server (If you want.)

#### `loglevel 0-2`

loglevel|number|description
-|-|-
loglevel|0|Note  every connection  `TIme, From IP, Port` to `Destination IP,Port`
loglevel|1|Note only prohibited connection `Time, From IP, Port` to `Destination IP,Port`
loglevel|2|Note only prohibited connection `Time`