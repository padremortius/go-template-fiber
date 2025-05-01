package httpserver

func mySkipper() []string {
	return []string{
		"/env", "/health", "/info", "/favicon.ico", "/prometheus",
	}
}
