## uinetd - Yet another udp & tcp forward tool just like rinetd

### Description
Redirects TCP & UDP connections from a local address & port to another address.
Configuration files is just in /etc/uintd.conf.(Like rinetd).
The process will open a socket to listen all the port from local bindings.

### Usage

local address|local bind port|remote address|remote port|`tcp` or `udp` or `all`
-|-|-|-|-
0.0.0.0|80|192.168.1.2|80|ALL
0.0.0.0|80|192.168.1.2|80|TCP

### Filter Connections 
`allow *.*.*.*` or `deny *.*.*.*` (IPv6 available)

### Log who is using your redirection server (If you want.)

#### `loglevel 0-4` (If you note 0 there will be no log recorded)

loglevel|number|description
-|-|-
loglevel|4|Note  every connection  `Time, From IP, Port` to `Destination IP,Port`, Errors
loglevel|3|Note only prohibited connection `Time, From IP, Port` to `Destination IP,Port` , Errors
loglevel|2|Note only prohibited connection `Time` , Errors
loglevel|1|Errors
