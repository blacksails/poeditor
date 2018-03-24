# poeditor [![GoDoc](https://godoc.org/github.com/blacksails/poeditor?status.svg)](https://godoc.org/github.com/blacksails/poeditor)

This is a delicious go API wrapper for [POEditor](https://poeditor.com). The project strives to convert
the POEditor REST API to an idiomatic go API.

## Getting started

Install the library with `go get` or `dep`

```
go get github.com/blacksails/poeditor
dep ensure -add github.com/blacksails/poeditor
```

Provide an API token and go nuts

```go
poe := poeditor.New("YOUR API TOKEN")

// Get projects
ps, _ := poe.ListProjects()

// Export all languages in project folders
wd, _ := os.Getwd()
for _, p := range ps {
    pDir := filepath.Join(wd, "translations", strconv.Itoa(p.ID))
    os.MkdirAll(pDir)
    ls, _ := p.ListLanguages()
    for _, l := range ls {
        f, _ := os.Create(pDir, fmt.Sprintf("%s.po", l.Code))
        l.Export(poeditor.FileFormatPO, []string{}, []string{}, f)
    }
}
```

## Wrapper completeness
All of the API endpoints have been implemented. A few of them lack proper
testing. Personally I dont use all of them, so I am only using a few in
production. If you find something that doesn't work please file an issue and
I will try to make a fix asap. Pull requests are also very welcome.
