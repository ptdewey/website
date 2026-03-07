{ pkgs, lib }:

pkgs.buildGoModule {
  src = ../.;
  name = "cedar";
  version = "0.1.0";
  vendorHash = "sha256-5h3jUNJk4hWPNYWXjFkLHuDJnU7xyhCwG9lsuP7jNcE=";

  meta = with lib; {
    description = "Lightweight static site generator.";
    homepage = "https://github.com/ptdewey/website";
    license = licenses.mit;
    maintainers = with maintainers; [ ptdewey ];
  };
}
