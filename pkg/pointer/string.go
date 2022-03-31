package pointer

func String(s string) *string {
	return &s
}

func ToString(p *string) string {
	if p == nil {
		return ""
	}

	return *p
}
