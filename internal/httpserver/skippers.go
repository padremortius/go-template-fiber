package httpserver

// func mySkipper(c *fiber.Ctx) bool {
// 	if (strings.HasPrefix(c.Path(), "/env")) ||
// 		(strings.HasPrefix(c.Path(), "/info")) ||
// 		(strings.HasPrefix(c.Path(), "/health")) ||
// 		(strings.HasPrefix(c.Path(), "/favicon.ico")) {
// 		return true
// 	}
// 	return false
// }

// func myJwtSkipper(c echo.Context) bool {
// 	if mySkipper(c) {
// 		return true
// 	}
// 	if strings.HasPrefix(c.Path(), "getToken") {
// 		return true
// 	}
// 	return false
// }
