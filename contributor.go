package poeditor

// AddAdmin adds a user as a project admin
func (p *Project) AddAdmin(name, email string) error {
	return p.post("/contributors/add", map[string]string{
		"name":  name,
		"email": email,
		"admin": "1",
	}, nil, nil)
}

// AddContributor adds a user as a language contributor
func (l *Language) AddContributor(name, email string) error {
	return l.post("/contributors/add", map[string]string{
		"name":  name,
		"email": email,
	}, nil, nil)
}
