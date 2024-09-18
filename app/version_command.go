package app

import "squirrel/types"

func VersionCommand(p types.Printer, appName string, appVersion string) Command {
	return func(args ...string) {
		p("{0} v{1}", appName, appVersion)
		p("{blue}http://github.com/ehsun7b/squirrel{/blue}\n")
	}
}
