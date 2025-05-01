package main

import (
	"github.com/padremortius/go-template-fiber/internal/app"
)

var (
	aBuildNumber    = ""
	aBuildTimeStamp = ""
	aGitBranch      = ""
	aGitHash        = ""
)

// @Title go-template-fiber
// @Version 1.0
// @Description This is a template of api-server with fiber router.
// @ContactName padremortius
// @ContactUrl http://misko.su/support
// @ContactEmail support@misko.su
// @LicenseName MIT
// @LicenseURL https://en.wikipedia.org/wiki/MIT_License

func main() {
	app.Run(aBuildNumber, aBuildTimeStamp, aGitBranch, aGitHash)
}
