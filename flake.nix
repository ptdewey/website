{
  description = "Patrick's Site (Astro)";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs }:
    let
      supportedSystems = [ "x86_64-linux" "aarch64-linux" "aarch64-darwin" "x86_64-darwin" ];
      forAllSystems = f:
        nixpkgs.lib.genAttrs supportedSystems (system: f {
          inherit system;
          pkgs = import nixpkgs { inherit system; };
        });
    in {
      devShells = forAllSystems ({ pkgs, ... }: {
        default = pkgs.mkShell {
          packages = [ pkgs.nodejs_22 pkgs.pnpm_10 pkgs.jujutsu ];
        };
      });

      apps = forAllSystems ({ pkgs, ... }:
        let
          runtimeInputs = [ pkgs.nodejs_22 pkgs.pnpm_10 pkgs.jujutsu ];
          buildSite = pkgs.writeShellApplication {
            name = "build-site";
            inherit runtimeInputs;
            text = ''
              pnpm install --frozen-lockfile
              pnpm build
            '';
          };
          devSite = pkgs.writeShellApplication {
            name = "dev-site";
            inherit runtimeInputs;
            text = ''
              pnpm install --frozen-lockfile
              pnpm dev
            '';
          };
        in {
          default = {
            type = "app";
            program = "${buildSite}/bin/build-site";
          };
          dev = {
            type = "app";
            program = "${devSite}/bin/dev-site";
          };
        });
    };
}
