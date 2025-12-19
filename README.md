# gwt

gwt is a personal TUI helper for managing `git worktree`.

It provides a simple interactive interface to:
- list existing worktrees
- switch between them
- create or delete worktrees safely

`gwt` is designed to be used via shell integration.

---

## Requirements

- git (with worktree support)
- zsh or bash

---

## Setup

Clone this repository and source the shell integration script.

### zsh

```sh
source /path/to/gwt/shell/gwt.zsh
````

### bash

```sh
source /path/to/gwt/shell/gwt.bash
```

Make sure `gwt-bin` is available in your PATH.

---

## Usage

```sh
gwt
```

This launches a TUI to select or manage git worktrees.

After selecting a worktree, your shell will `cd` into it automatically.

---

## Notes

* `gwt-bin` is an internal command used by the `gwt` shell function.
  You normally do not need to run it directly.
* Errors are printed to stderr. The selected path is printed to stdout.

You can check internal usage with:

```sh
gwt-bin --help
```

---

## Design Philosophy

* Shell-first (no `cd` inside binaries)
* Minimal output
* Safe defaults
* Personal workflow focused

---

## Status

This is a personal tool.
The interface and behavior may change without notice.
