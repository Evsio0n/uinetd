module uinetd

go 1.14

require (
	github.com/evsio0n/log v0.0.0-20200721012947-f56dac26ab64
	github.com/evsio0n/uinetd/pkg/check v0.0.0-master-incapable // indirect
	github.com/evsio0n/uinetd/pkg/logger v0.0.0-master-incapable
	github.com/evsio0n/uinetd/utils v0.0.0-master-incapable
)

replace (
	github.com/evsio0n/uinetd/pkg/check => ./pkg/check
	github.com/evsio0n/uinetd/pkg/logger => ./pkg/logger
	github.com/evsio0n/uinetd/utils => ./utils
)
