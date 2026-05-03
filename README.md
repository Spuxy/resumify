## Resumify
Generator tool in Go to build your awesome CV :page_facing_up:

## Why
Reason is to generate simple and clean site based on most popular configuration extension YML, where people can easily add new experiences and do it via VCS.

## Usage

Resumify supports two modes: a one-shot **build** (great for GitHub Pages / CI) and a long-running **server** (great for local preview).

### 1. Write your CV
Make your own repository on GitHub, e.g. `https://github.com/{username}/cv`, and add a `cv.yml`.

### 2a. Build mode (GitHub Pages / CI)
Generate `index.html` once and exit — perfect for committing to a `gh-pages` branch from a GitHub Action.

```sh
# from a local file
./resumify --build --src cv.yml

# from a remote raw URL
./resumify --build --src https://raw.githubusercontent.com/{username}/cv/{branch}/cv.yml

# custom output path
./resumify --build --src cv.yml --out dist/index.html
```

### 2b. Server mode (local preview)
Copy `.env.example` to `src/.env` and edit it to point at your `cv.yml`, then run resumify:

```sh
cp .env.example src/.env
# edit src/.env: set src=https://raw.githubusercontent.com/<you>/cv/master/cv.yml
./resumify
```

- `/preview` — render the CV in the browser
- `/generate` — render and write `index.html` to disk

`--src` also overrides the configured source in server mode. `src/.env` is gitignored so it never gets committed.

### 2c. GitHub Action (zero-install)
Use the bundled composite action from your CV repo — no Go toolchain or local build needed. It checks out resumify, builds it, and drops `index.html` in your workspace.

Create `.github/workflows/deploy.yml` in your CV repo:

```yaml
name: Deploy CV
on:
  push:
    branches: [main]
  workflow_dispatch:

permissions:
  contents: read
  pages: write
  id-token: write

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Generate CV
        uses: Spuxy/resumify@master
        with:
          src: cv.yml
          out: dist/index.html

      - uses: actions/upload-pages-artifact@v3
        with:
          path: dist

  deploy:
    needs: build
    runs-on: ubuntu-latest
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    steps:
      - id: deployment
        uses: actions/deploy-pages@v4
```

**Action inputs:**

| Input | Description | Default |
|---|---|---|
| `src` | Path to `cv.yml` (relative to your repo) or a raw URL | `cv.yml` |
| `out` | Output path for the generated HTML | `index.html` |
| `ref` | Resumify ref (tag, branch, or SHA) to pin | `master` |
| `go-version` | Go version used to build resumify | `1.21` |
| `target-repo` | `owner/repo` to push the generated HTML to (e.g. `username/username.github.io`). Empty = skip push | `""` |
| `target-branch` | Branch in `target-repo` to push to | `main` |
| `target-path` | Path inside `target-repo` for the HTML | `index.html` |
| `token` | PAT / fine-grained token with `contents:write` on `target-repo`. Required when `target-repo` is set | `""` |
| `commit-message` | Commit message used when pushing | `chore: update CV via resumify` |
| `commit-user-name` | Git `user.name` for the push commit | `github-actions[bot]` |
| `commit-user-email` | Git `user.email` for the push commit | `41898282+github-actions[bot]@users.noreply.github.com` |

### 2d. Push to a separate `username.github.io` repo
If you keep `cv.yml` in one repo and serve from `username/username.github.io`, give the action a token and let it push directly — no Pages workflow needed in either repo.

1. Create a fine-grained PAT scoped to your `<username>.github.io` repo with **Contents: Read and write**.
2. Add it as a secret named `CV_DEPLOY_TOKEN` in your CV repo.
3. Drop this workflow into your CV repo at `.github/workflows/deploy.yml` (also available as a ready-to-copy file at [`examples/cv-repo-workflow.yml`](examples/cv-repo-workflow.yml)):

```yaml
name: Deploy CV
on:
  push:
    branches: [main]
  workflow_dispatch:

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: Spuxy/resumify@master
        with:
          src: cv.yml
          token: ${{ secrets.CV_DEPLOY_TOKEN }}
```

That's it — the workflow has no hardcoded usernames. When `target-repo` is empty, the action defaults to `<owner>/<owner>.github.io` (lowercased) using the workflow's repo owner, so anyone forking the CV repo just adds their own `CV_DEPLOY_TOKEN` secret and ships.

Override only if your setup differs:

```yaml
      - uses: Spuxy/resumify@master
        with:
          src: cv.yml
          token: ${{ secrets.CV_DEPLOY_TOKEN }}
          target-repo: my-org/portfolio
          target-branch: gh-pages
          target-path: cv/index.html
```

The action skips the commit when the rendered HTML is unchanged, so re-runs are safe.

## Flags
| Flag | Description | Default |
|---|---|---|
| `--build` | Generate `index.html` and exit | `false` |
| `--src` | YAML source (URL or local path); overrides `.env` | `""` |
| `--out` | Output path for `--build` mode | `index.html` |

## CV YAML fields
Beyond the obvious ones, two top-level fields control where the template pulls media from — set them in your `cv.yml`:

| Field | Description | Example |
|---|---|---|
| `photo` | Full URL to your profile photo | `https://myuser.github.io/images/me.jpg` |
| `assets` | Base URL prepended to gallery `picture_path` entries | `https://raw.githubusercontent.com/myuser/cv/master` |

If you leave them empty the template still renders, but profile/gallery images will be broken.

## TODO
Create CI/CD <br />
Create unit test <br />
Implement more colorschemes (only gruvbox is implemented) <br />
