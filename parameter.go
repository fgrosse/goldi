package goldi

func isParameter(p string) bool {
	if len(p) < 2 {
		return false
	}

	return p[0] == '%' && p[len(p)-1] == '%'
}
