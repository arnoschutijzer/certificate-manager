# certificate-manager

```bash
$ go run ./cmd/cli/main.go
```

A utility that (ab)uses the 1password CLI to find outdated certificates.
The 1password CLI takes a long time (~1s) to retrieve the details of a secret so we use goroutines to fetch them in parallel.

You can add custom filtering based on categories and tags via customizing the code in `internal/onepassword/onepassword.go`.

See https://developer.1password.com/docs/cli
