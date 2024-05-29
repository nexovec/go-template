go install github.com/g4s8/envdoc@latest # generates .env documentation
go install golang.org/x/tools/gopls@latest # go language server
go install github.com/kisielk/godepgraph@latest # renders go dependency graph svg, see ./scripts/draw_deps.sh
sudo apt install -y graphviz # as a godepgraph dependency

# useful CLI programs
go install github.com/jesseduffield/lazydocker@latest
go install github.com/jesseduffield/lazygit@latest
go install github.com/xxxserxxx/gotop/v4/cmd/gotop@latest
go install github.com/hhatto/gocloc/cmd/gocloc@latest
go install github.com/dundee/gdu/v5/cmd/gdu@latest

# code generation and linters
go install github.com/a-h/templ/cmd/templ@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# migrations, db access
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

