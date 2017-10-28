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

// RemoveAdmin removes a user as a project admin
func (p *Project) RemoveAdmin(email string) error {
	return p.post("/contributors/remove", map[string]string{"email": email}, nil, nil)
}

// RemoveContributor removes a contributor from the language
func (l *Language) RemoveContributor(email string) error {
	return l.post("/contributors/remove", map[string]string{"email": email}, nil, nil)
}
