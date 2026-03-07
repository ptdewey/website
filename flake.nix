{
  description = "Patrick's Site Flake";

  inputs = { nixpkgs.url = "nixpkgs/nixpkgs-unstable"; };

  outputs = { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { inherit system; };
      cedarDrv = pkgs.callPackage ./nix/default.nix { };
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
        cedar = {
          type = "app";
          program = "${cedarDrv}/bin/cedar";
        };

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
