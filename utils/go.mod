module github.com/evsio0n/uinetd/utils

go 1.14

replace (
	github.com/evsio0n/uinetd/pkg/check => ../pkg/check
	github.com/evsio0n/uinetd/pkg/logger => ../pkg/logger
)

require (
	github.com/evsio0n/uinetd/pkg/check v0.0.0-00010101000000-000000000000
	github.com/evsio0n/uinetd/pkg/logger v0.0.0-00010101000000-000000000000
)
