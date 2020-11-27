module github.com/evsio0n/uinetd/utils

go 1.14
require (
	github.com/evsio0n/uinetd/pkg/logger v0.0.0-master-incapable
	github.com/evsio0n/uinetd/pkg/check v0.0.0-master-incapable
)

replace (
 	github.com/evsio0n/uinetd/pkg/logger => ../pkg/logger
    github.com/evsio0n/uinetd/pkg/check => ../pkg/check
 )