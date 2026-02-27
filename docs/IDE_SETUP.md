# IDE Setup Guide

This guide helps you configure your IDE for optimal development experience with nrfiber.

## Table of Contents

- [Visual Studio Code](#visual-studio-code)
- [GoLand / IntelliJ IDEA](#goland--intellij-idea)
- [Vim / Neovim](#vim--neovim)
- [Sublime Text](#sublime-text)
- [General Go Setup](#general-go-setup)
- [Troubleshooting](#troubleshooting)

---

## Visual Studio Code

### Recommended Extensions

1. **Go** (golang.go) - Official Go extension
   ```
   code --install-extension golang.go
   ```

2. **Go Test Explorer** (premparihar.gotestexplorer) - For running tests
   ```
   code --install-extension premparihar.gotestexplorer
   ```

3. **Error Lens** (usernamehw.errorlens) - Inline error display
   ```
   code --install-extension usernamehw.errorlens
   ```

### Settings Configuration

Create or update `.vscode/settings.json` in your project:

```json
{
  "go.useLanguageServer": true,
  "go.lintTool": "golangci-lint",
  "go.lintOnSave": "workspace",
  "go.formatTool": "gofmt",
  "go.toolsManagement.autoUpdate": true,
  "go.testFlags": ["-v"],
  "go.testTimeout": "30s",
  "go.coverOnSave": false,
  "go.coverOnSingleTest": false,
  "go.coverOnSingleTestFile": false,
  "go.coverageDecorator": {
    "type": "gutter"
  },
  "[go]": {
    "editor.formatOnSave": true,
    "editor.codeActionsOnSave": {
      "source.organizeImports": true
    }
  },
  "gopls": {
    "experimentalWorkspaceModule": true,
    "build.directoryFilters": [
      "-node_modules",
      "-vendor"
    ],
    "ui.semanticTokens": true,
    "ui.completion.usePlaceholders": true
  }
}
```

### Launch Configuration

Create `.vscode/launch.json` for debugging:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Launch Example (v3)",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/examples/fiber-v3-basic",
      "env": {
        "NEW_RELIC_LICENSE_KEY": "your-license-key-here"
      }
    },
    {
      "name": "Launch Example (v2)",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/examples/fiber-v2-basic",
      "env": {
        "NEW_RELIC_LICENSE_KEY": "your-license-key-here"
      }
    },
    {
      "name": "Test Current Package",
      "type": "go",
      "request": "launch",
      "mode": "test",
      "program": "${fileDirname}"
    },
    {
      "name": "Test Current File",
      "type": "go",
      "request": "launch",
      "mode": "test",
      "program": "${fileDirname}",
      "args": [
        "-test.run",
        "${selectedText}"
      ]
    }
  ]
}
```

### Tasks Configuration

Create `.vscode/tasks.json` for common tasks:

```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "Test v3",
      "type": "shell",
      "command": "cd v3 && go test -v ./...",
      "group": "test",
      "presentation": {
        "reveal": "always",
        "panel": "new"
      }
    },
    {
      "label": "Test v2",
      "type": "shell",
      "command": "cd v2 && go test -v ./...",
      "group": "test",
      "presentation": {
        "reveal": "always",
        "panel": "new"
      }
    },
    {
      "label": "Test All",
      "type": "shell",
      "command": "cd v3 && go test -v ./... && cd ../v2 && go test -v ./...",
      "group": {
        "kind": "test",
        "isDefault": true
      },
      "presentation": {
        "reveal": "always",
        "panel": "new"
      }
    },
    {
      "label": "Lint",
      "type": "shell",
      "command": "golangci-lint run ./...",
      "group": "test",
      "presentation": {
        "reveal": "always",
        "panel": "shared"
      }
    },
    {
      "label": "Run Example v3",
      "type": "shell",
      "command": "cd examples/fiber-v3-basic && go run main.go",
      "presentation": {
        "reveal": "always",
        "panel": "dedicated"
      }
    },
    {
      "label": "Run Example v2",
      "type": "shell",
      "command": "cd examples/fiber-v2-basic && go run main.go",
      "presentation": {
        "reveal": "always",
        "panel": "dedicated"
      }
    }
  ]
}
```

### Workspace Recommendations

Create `.vscode/extensions.json`:

```json
{
  "recommendations": [
    "golang.go",
    "premparihar.gotestexplorer",
    "usernamehw.errorlens",
    "eamodio.gitlens",
    "github.vscode-pull-request-github"
  ]
}
```

---

## GoLand / IntelliJ IDEA

### Initial Setup

1. **Install Go Plugin** (if not already installed)
   - Go to `Settings/Preferences → Plugins`
   - Search for "Go"
   - Install and restart IDE

2. **Configure GOPATH and GOROOT**
   - Go to `Settings/Preferences → Go → GOROOT`
   - Ensure your Go installation is detected
   - Verify GOPATH is set correctly

### Project Configuration

1. **Enable Go Modules**
   - Go to `Settings/Preferences → Go → Go Modules`
   - Check "Enable Go modules integration"
   - Set "Vendoring mode" to "Disabled" (unless you use vendoring)

2. **Configure Code Style**
   - Go to `Settings/Preferences → Editor → Code Style → Go`
   - Use default Go formatting (gofmt)
   - Enable "Optimize imports on the fly"

3. **Setup File Watchers**
   - Go to `Settings/Preferences → Tools → File Watchers`
   - Add watchers for:
     - `gofmt` (if not auto-enabled)
     - `golangci-lint` (optional)

### Run Configurations

#### For v3 Example:
1. Right-click `examples/fiber-v3-basic/main.go`
2. Select `Modify Run Configuration`
3. Add environment variables:
   ```
   NEW_RELIC_LICENSE_KEY=your-license-key
   NEW_RELIC_APP_NAME=nrfiber-v3-test
   ```

#### For v2 Example:
1. Right-click `examples/fiber-v2-basic/main.go`
2. Select `Modify Run Configuration`
3. Add environment variables:
   ```
   NEW_RELIC_LICENSE_KEY=your-license-key
   NEW_RELIC_APP_NAME=nrfiber-v2-test
   ```

#### For Tests:
1. Go to `Run → Edit Configurations`
2. Click `+` → `Go Test`
3. Configure:
   - **Name**: Test v3
   - **Directory**: `<project>/v3`
   - **Pattern**: `.*`
   - Click OK

### Useful Keyboard Shortcuts

| Action | macOS | Windows/Linux |
|--------|-------|---------------|
| Run Tests | `Ctrl + Shift + R` | `Ctrl + Shift + F10` |
| Debug Tests | `Ctrl + Shift + D` | `Ctrl + Shift + F9` |
| Go to Definition | `Cmd + B` | `Ctrl + B` |
| Find Usages | `Alt + F7` | `Alt + F7` |
| Refactor/Rename | `Shift + F6` | `Shift + F6` |
| Format Code | `Cmd + Alt + L` | `Ctrl + Alt + L` |
| Optimize Imports | `Ctrl + Alt + O` | `Ctrl + Alt + O` |

---

## Vim / Neovim

### Using vim-go

1. **Install vim-go**
   
   Add to your `.vimrc` or `init.vim`:
   ```vim
   Plug 'fatih/vim-go', { 'do': ':GoUpdateBinaries' }
   ```

2. **Configure vim-go**
   
   ```vim
   " Enable auto-formatting on save
   let g:go_fmt_autosave = 1
   let g:go_fmt_command = "gofmt"
   
   " Enable imports organization
   let g:go_imports_autosave = 1
   
   " Enable syntax highlighting
   let g:go_highlight_types = 1
   let g:go_highlight_fields = 1
   let g:go_highlight_functions = 1
   let g:go_highlight_function_calls = 1
   let g:go_highlight_operators = 1
   let g:go_highlight_extra_types = 1
   
   " Use gopls
   let g:go_def_mode='gopls'
   let g:go_info_mode='gopls'
   
   " Run tests
   autocmd FileType go nmap <leader>t <Plug>(go-test)
   autocmd FileType go nmap <leader>tf <Plug>(go-test-func)
   autocmd FileType go nmap <Leader>c <Plug>(go-coverage-toggle)
   
   " Navigate between functions
   autocmd FileType go nmap <Leader>n <Plug>(go-alternate-edit)
   ```

3. **Key Mappings for nrfiber Development**
   
   ```vim
   " Test v3
   autocmd FileType go nmap <leader>t3 :!cd v3 && go test -v ./...<CR>
   
   " Test v2
   autocmd FileType go nmap <leader>t2 :!cd v2 && go test -v ./...<CR>
   
   " Run example
   autocmd FileType go nmap <leader>rv3 :!cd examples/fiber-v3-basic && go run main.go<CR>
   autocmd FileType go nmap <leader>rv2 :!cd examples/fiber-v2-basic && go run main.go<CR>
   ```

### Using coc.nvim with gopls

1. **Install coc.nvim**
   ```vim
   Plug 'neoclide/coc.nvim', {'branch': 'release'}
   ```

2. **Install gopls**
   ```vim
   :CocInstall coc-go
   ```

3. **Configure coc-settings.json**
   ```json
   {
     "go.goplsOptions": {
       "experimentalWorkspaceModule": true
     },
     "languageserver": {
       "golang": {
         "command": "gopls",
         "rootPatterns": ["go.mod"],
         "filetypes": ["go"]
       }
     }
   }
   ```

---

## Sublime Text

### Install GoSublime

1. Install Package Control (if not installed)
2. Press `Ctrl/Cmd + Shift + P`
3. Type "Install Package" and press Enter
4. Search for "GoSublime" and install

### Configure GoSublime

Create or edit `Packages/User/GoSublime.sublime-settings`:

```json
{
  "env": {
    "GOPATH": "$HOME/go",
    "PATH": "$GOPATH/bin:$PATH"
  },
  "fmt_cmd": ["gofmt"],
  "on_save": [
    {
      "cmd": "gs_fmt",
      "args": {}
    }
  ],
  "autocomplete_tests": true,
  "autocomplete_closures": true
}
```

### Build System

Create a build system for testing: `Tools → Build System → New Build System`:

```json
{
  "shell_cmd": "go test -v",
  "file_regex": "^\\s*(.+\\.go):(\\d+):():(.*)$",
  "working_dir": "${file_path}",
  "selector": "source.go"
}
```

---

## General Go Setup

### Install Go Tools

Install essential Go development tools:

```bash
# gopls - Language Server
go install golang.org/x/tools/gopls@latest

# golangci-lint - Linter
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# delve - Debugger
go install github.com/go-delve/delve/cmd/dlv@latest

# go test coverage
go install golang.org/x/tools/cmd/cover@latest

# godoc - Documentation
go install golang.org/x/tools/cmd/godoc@latest
```

### Project-Specific Configuration

#### Create `.golangci.yml`

```yaml
run:
  timeout: 5m
  tests: true
  skip-dirs:
    - vendor
    - examples

linters:
  enable:
    - gofmt
    - golint
    - govet
    - errcheck
    - staticcheck
    - unused
    - gosimple
    - structcheck
    - varcheck
    - ineffassign
    - deadcode
    - typecheck

linters-settings:
  gofmt:
    simplify: true
  golint:
    min-confidence: 0.8
  govet:
    check-shadowing: true

issues:
  exclude-rules:
    - path: _test\.go
      linters:
        - errcheck
```

#### Create `.editorconfig`

```ini
root = true

[*]
charset = utf-8
end_of_line = lf
insert_final_newline = true
trim_trailing_whitespace = true

[*.go]
indent_style = tab
indent_size = 4

[*.{yml,yaml}]
indent_style = space
indent_size = 2

[*.{json,md}]
indent_style = space
indent_size = 2
```

### Environment Variables

Add to your shell profile (`.bashrc`, `.zshrc`, etc.):

```bash
# Go environment
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
export GO111MODULE=on

# New Relic (for testing)
export NEW_RELIC_LICENSE_KEY="your-license-key-here"
export NEW_RELIC_APP_NAME="nrfiber-dev"

# Helpful aliases
alias gotest='go test -v ./...'
alias gotest3='cd v3 && go test -v ./... && cd ..'
alias gotest2='cd v2 && go test -v ./... && cd ..'
alias golint='golangci-lint run ./...'
```

---

## Troubleshooting

### gopls Not Working

**Problem**: Code completion or navigation not working.

**Solution**:
```bash
# Reinstall gopls
go install golang.org/x/tools/gopls@latest

# Clear gopls cache
rm -rf ~/.cache/gopls

# Restart your IDE/editor
```

### Module Download Issues

**Problem**: Cannot download modules or dependencies.

**Solution**:
```bash
# Clean module cache
go clean -modcache

# Update dependencies
go get -u ./...
go mod tidy

# Verify modules
go mod verify
```

### Import Path Not Recognized

**Problem**: IDE doesn't recognize `nrfiber/v3` or `nrfiber/v2` imports.

**Solution**:
1. Ensure you're in the correct directory
2. Run `go mod download`
3. Reload your IDE workspace
4. For VSCode: Run "Go: Restart Language Server"
5. For GoLand: File → Invalidate Caches / Restart

### Tests Not Running

**Problem**: Cannot run tests from IDE.

**Solution**:
1. Check you're in the correct directory (`v2/` or `v3/`)
2. Ensure `go.mod` exists in that directory
3. Run `go test ./...` from command line first
4. Check IDE test configuration includes correct working directory

### Debugger Issues

**Problem**: Breakpoints not hitting or debugger not starting.

**Solution**:
```bash
# Reinstall delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Check dlv works
dlv version

# For IDE: Check debug configuration includes correct program path
```

### Format on Save Not Working

**Problem**: Code not formatting automatically.

**Solution**:
- **VSCode**: Check `"editor.formatOnSave": true` in settings
- **GoLand**: Enable "Reformat code" in Settings → Tools → Actions on Save
- **Vim**: Ensure `let g:go_fmt_autosave = 1` is set
- Verify `gofmt` is in your PATH: `which gofmt`

---

## Additional Resources

- [Official Go Editor Setup](https://golang.org/doc/editors.html)
- [gopls Documentation](https://github.com/golang/tools/tree/master/gopls)
- [vim-go Tutorial](https://github.com/fatih/vim-go/wiki)
- [VSCode Go Extension](https://github.com/golang/vscode-go)
- [GoLand Documentation](https://www.jetbrains.com/help/go/)

---

## Contributing IDE Configurations

If you have IDE configurations that work well for nrfiber development, please contribute them!

1. Create a configuration file in `.vscode/`, `.idea/`, etc.
2. Add it to version control (if appropriate)
3. Document it in this guide
4. Submit a pull request

---

## Quick Start Checklist

- [ ] IDE installed with Go support
- [ ] Go 1.25.0+ installed
- [ ] `gopls` installed and configured
- [ ] `golangci-lint` installed
- [ ] Editor configured for format-on-save
- [ ] Run configurations set up for examples
- [ ] Test configurations working
- [ ] Environment variables configured
- [ ] Can run `go test ./...` successfully
- [ ] Can run example applications

Once all items are checked, you're ready to develop with nrfiber!