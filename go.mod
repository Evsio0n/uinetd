module github.com/evsio0n/uinetd

go 1.14

replace (
	github.com/evsio0n/uinetd/pkg/check => ./pkg/check
	github.com/evsio0n/uinetd/pkg/logger => ./pkg/logger
	github.com/evsio0n/uinetd/utils => ./utils
)

require github.com/evsio0n/uinetd/utils v0.0.0-00010101000000-000000000000
