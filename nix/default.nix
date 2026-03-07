{ pkgs ? import <nixpkgs> { }, lib ? pkgs.lib }:

pkgs.buildGoModule {
  src = ../.;
  name = "cedar";
  version = "0.1.0";
  vendorHash = "sha256-vk6/9GiZ/meZ621AA0ClDqQTrcwlmQbY6ZnDJP9bOHo=";

  meta = with lib; {
    description = "Lightweight static site generator.";
    homepage = "https://github.com/ptdewey/website";
    license = licenses.mit;
    maintainers = with maintainers; [ ptdewey ];
  };
}
