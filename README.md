# certificate-manager

A utility that (ab)uses the 1password CLI to find outdated certificates.
The 1password CLI takes a long time (~1s) to retrieve the details of a secret so the secrets are cached in a sqlite database and refreshed when the `updated_at` changes.

You can add custom filtering based on categories and tags via customizing the code in `onepassword.go`.

See https://developer.1password.com/docs/cli
