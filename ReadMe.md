Setup
=====

Install go -> https://go.dev/doc/install
  download installer, run it.
    setup path in shell config (~/zprofile)
      export GOPATH=~/go/bin
        test: go -v

To run with auto changes, download and config wgo package -> https://github.com/bokwoon95/wgo
  $ go install github.com/bokwoon95/wgo@latest
    setup alias in shell settings in .zshrc
      https://github.com/bokwoon95/wgo

We will use a routing package, the chi router
  https://github.com/go-chi/chi
  % go get -u github.com/go-chi/chi/v5

We need a postgres driver to talk to the database
  https://github.com/jackc/pgx/
  % go get github.com/jackc/pgx/v4/stdlib

We need a databased migrations package
  https://github.com/pressly/goose
  % go install github.com/pressly/goose/v3/cmd/goose@latest
	% go get github.com/pressly/goose/v3/cmd/goose@latest

We need a package to help with encryption
  % go get golang.org/x/crypto/bcrypt


To run project
==============
% wgo run main.go {-port 8000}

To run, watching more file types than just .go files:
% wgo run -file .gohtml -file .sql -file .css -file .sql -file .js main.go -port 8000

