{
  description = "Patrick's Site Flake";

  inputs = {
    nixpkgs.url = "nixpkgs/nixpkgs-unstable";
    cedar.url = "github:ptdewey/cedar";
    # cedar.url = "git+file:///home/patrick/projects/cedar";
  };

  outputs = { self, nixpkgs, cedar }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
      cedarDrv = cedar.packages.${system}.default;
      buildSite = import ./nix/build-site.nix { inherit pkgs cedarDrv; };
      buildTailwind = import ./nix/build-tailwind.nix { inherit pkgs; };
    in {
      # `nix develop` directives
      devShells.${system}.default = pkgs.mkShell {
        packages = with pkgs; [ tailwindcss_4 ];
        buildInputs = [ cedarDrv ];
      };

      # `nix run` directives
      apps.${system} = {
        # Build SSG static HTML files
        cedar = cedar.apps.${system}.cedar;

        # Build stylesheets
        tailwind = {
          type = "app";
          program = "${buildTailwind}/bin/run";
        };

        # Build site (static HTML files and stylesheets)
        default = {
          type = "app";
          program = "${buildSite}/bin/run";
        };
      };
    };
}
