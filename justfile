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
