all:
    @nix run

build:
    @nix run .#cedar

style:
    @nix run .#tailwind

shell:
    @nix develop --command "zsh"

clean:
    @rm -rf public

serve:
    @pushd public || exit 0 && python -m http.server && popd || exit 0

test:
    @go test ./... -cover -coverprofile=cover.out
