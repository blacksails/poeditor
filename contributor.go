package poeditor

// ListContributors lists all contributors registered in POEditor
func (poe *POEditor) ListContributors() ([]Contributor, error) {
	var res listContributorsResult
	err := poe.post("/contributors/list", nil, nil, &res)
	if err != nil {
		return []Contributor{}, err
	}
	return res.toContributors(), nil
}

// ListContributors lists all contributors registered under the project
func (p *Project) ListContributors() ([]Contributor, error) {
	var res listContributorsResult
	err := p.post("/contributors/list", nil, nil, &res)
	if err != nil {
		return []Contributor{}, err
	}
	return res.toContributors(), nil
}

// ListContributors lists all contributors registered under the language
func (l *Language) ListContributors() ([]Contributor, error) {
	var res listContributorsResult
	err := l.post("/contributors/list", nil, nil, &res)
	if err != nil {
		return []Contributor{}, err
	}
	return res.toContributors(), nil
}

// AddContributor adds a user as a project admin
func (p *Project) AddContributor(name, email string) error {
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

// RemoveContributor removes a user as a project admin
func (p *Project) RemoveContributor(email string) error {
	return p.post("/contributors/remove", map[string]string{"email": email}, nil, nil)
}

// RemoveContributor removes a contributor from the language
func (l *Language) RemoveContributor(email string) error {
	return l.post("/contributors/remove", map[string]string{"email": email}, nil, nil)
}

type Contributor struct {
	Name        string
	Email       string
	Permissions []Permission
}

type Permission struct {
	ProjectID   string
	ProjectName string
	Type        string
}

type listContributorsResult struct {
	Contributors []contributor `json:"contributors"`
}

func (r listContributorsResult) toContributors() []Contributor {
	cs := make([]Contributor, len(r.Contributors))
	for i, c := range r.Contributors {
		cs[i] = Contributor{
			Name:  c.Name,
			Email: c.Email,
		}
		ps := make([]Permission, len(c.Permissions))
		for ip, p := range c.Permissions {
			ps[ip] = Permission{
				ProjectID:   p.Project.ID,
				ProjectName: p.Project.Name,
				Type:        p.Type,
			}
		}
		cs[i].Permissions = ps
	}
	return cs
}

type contributor struct {
	Name        string                  `json:"name"`
	Email       string                  `json:"email"`
	Permissions []contributorPermission `json:"permissions"`
}

type contributorPermission struct {
	Project contributorProjectPermission `json:"project"`
	Type    string                       `json:"type"`
}

type contributorProjectPermission struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
