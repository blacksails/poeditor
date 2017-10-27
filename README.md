# poeditor [![GoDoc](https://godoc.org/github.com/blacksails/poeditor?status.svg)](https://godoc.org/github.com/blacksails/poeditor)

This is a delicious go API wrapper for POEditor. The project strives to convert
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
The following endpoints are wrapped/still need to get wrapped

- [x] Projects
  - [x] List
  - [x] View
  - [x] Add
  - [x] Update
  - [x] Delete
  - [x] Upload
  - [x] Sync
  - [x] Export
- [ ] Languages
  - [x] Available
  - [x] List
  - [x] Add
  - [ ] Update
  - [ ] Delete
- [ ] Terms
  - [ ] List
  - [ ] Add
  - [ ] Update
  - [ ] Delete
  - [ ] Add comment
- [ ] Contributors
  - [ ] List
  - [ ] Add
  - [ ] Remove
