package app

func VersionCommand(p Printer, appName string, appVersion string) Command {
	return func(args ...string) {
		p("{0} v{1}", appName, appVersion)
	}
}
