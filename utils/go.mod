module utils

go 1.14
require (
	github.com/evsio0n/uinetd/logger v0.0.0-master-incapable
	github.com/evsio0n/uinetd/check v0.0.0-master-incapable
)

replace (
 	github.com/evsio0n/uinetd/logger => ../pkg/logger
    github.com/evsio0n/uinetd/check => ../pkg/check
 )