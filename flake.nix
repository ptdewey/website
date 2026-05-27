{
  description = "patrick's static site";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs =
    { self, nixpkgs }:
    let
      systems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];

      forAllSystems = nixpkgs.lib.genAttrs systems;
    in
    {
      packages = forAllSystems (
        system:
        let
          pkgs = import nixpkgs { inherit system; };
          beamPackages = pkgs.beamPackages;

          src = pkgs.lib.cleanSourceWith {
            src = ./.;
            filter = path: type:
              let
                rel = pkgs.lib.removePrefix (toString ./. + "/") (toString path);
              in
              pkgs.lib.cleanSourceFilter path type
              && ! pkgs.lib.hasPrefix "output/" rel
              && ! pkgs.lib.hasPrefix "_build/" rel
              && ! pkgs.lib.hasPrefix "deps/" rel;
          };

          mixFodDeps = beamPackages.fetchMixDeps {
            pname = "site-deps";
            version = "0.1.0";
            inherit src;
            hash = "sha256-1JbQujRXl0KrMU98DEN7jnaoGC7m02X8tIHk3Tuw+u4=";
          };

          nativeArtifacts = {
            x86_64-linux = {
              mdex = {
                file = "libcomrak_nif-v0.12.2-nif-2.15-x86_64-unknown-linux-gnu.so.tar.gz";
                hash = "sha256:a13679dd322957e415839b0ca4d4b759241dea725c8616bff0906eca515cd77f";
              };
              lumis = {
                file = "liblumis_nif-v0.5.0-nif-2.15-x86_64-unknown-linux-gnu.so.tar.gz";
                hash = "sha256:02c58f7c27dbedd2ef577145a9327317c1d8e75f512eb9c6279256dafb1947f5";
              };
            };
            aarch64-linux = {
              mdex = {
                file = "libcomrak_nif-v0.12.2-nif-2.15-aarch64-unknown-linux-gnu.so.tar.gz";
                hash = "sha256:9111c332fbac1d10bd027c9251ce0a7c6f8c0665cc71497d4562ddd03c70a98c";
              };
              lumis = {
                file = "liblumis_nif-v0.5.0-nif-2.15-aarch64-unknown-linux-gnu.so.tar.gz";
                hash = "sha256:d83a55467895f86cd03c9062bc2f00686d945849ef47d5497a07d374ee6ec06f";
              };
            };
            x86_64-darwin = {
              mdex = {
                file = "libcomrak_nif-v0.12.2-nif-2.15-x86_64-apple-darwin.so.tar.gz";
                hash = "sha256:a9f5f7fe297ccacb4d521662f8ccb136ef2b9cc2526992f7b9f4336278a176d2";
              };
              lumis = {
                file = "liblumis_nif-v0.5.0-nif-2.15-x86_64-apple-darwin.so.tar.gz";
                hash = "sha256:f164c99fa1532f2e9d589558956a70a6140e5bb8407ffdf06883fe37ab659717";
              };
            };
            aarch64-darwin = {
              mdex = {
                file = "libcomrak_nif-v0.12.2-nif-2.15-aarch64-apple-darwin.so.tar.gz";
                hash = "sha256:05a8e8b77b4181491478bbbc61ec907b28f8fabec2577673b83445be87f88597";
              };
              lumis = {
                file = "liblumis_nif-v0.5.0-nif-2.15-aarch64-apple-darwin.so.tar.gz";
                hash = "sha256:b32b7b703d180620c011fc13efaedf8f5ea3a5a744bcdcc7ae0666194dda2a12";
              };
            };
          }.${system};

          rustlerCache = pkgs.runCommand "site-rustler-precompiled-cache" { } ''
            mkdir -p "$out"
            cp ${pkgs.fetchurl {
              url = "https://github.com/leandrocp/mdex/releases/download/v0.12.2/${nativeArtifacts.mdex.file}";
              hash = nativeArtifacts.mdex.hash;
            }} "$out/${nativeArtifacts.mdex.file}"
            cp ${pkgs.fetchurl {
              url = "https://github.com/leandrocp/lumis/releases/download/hex-lumis/v0.5.0/${nativeArtifacts.lumis.file}";
              hash = nativeArtifacts.lumis.hash;
            }} "$out/${nativeArtifacts.lumis.file}"
          '';
        in
        {
          default = pkgs.stdenv.mkDerivation {
            pname = "site";
            version = "0.1.0";
            inherit src;

            nativeBuildInputs = [
              beamPackages.elixir
              beamPackages.hex
              beamPackages.rebar
              beamPackages.rebar3
            ];

            env = {
              MIX_ENV = "prod";
              HEX_OFFLINE = "1";
              RUSTLER_PRECOMPILED_GLOBAL_CACHE_PATH = rustlerCache;
              MIX_REBAR = "${beamPackages.rebar}/bin/rebar";
              MIX_REBAR3 = "${beamPackages.rebar3}/bin/rebar3";
              LANG = if pkgs.stdenv.hostPlatform.isLinux then "C.UTF-8" else "C";
              LC_CTYPE = if pkgs.stdenv.hostPlatform.isLinux then "C.UTF-8" else "UTF-8";
            };

            configurePhase = ''
              runHook preConfigure

              export MIX_HOME="$TMPDIR/mix"
              export HEX_HOME="$TMPDIR/hex"
              export MIX_DEPS_PATH="$TMPDIR/deps"
              export REBAR_GLOBAL_CONFIG_DIR="$TMPDIR/rebar3"
              export REBAR_CACHE_DIR="$TMPDIR/rebar3.cache"

              cp --no-preserve=mode -R "${mixFodDeps}" "$MIX_DEPS_PATH"
              ln -s "$MIX_DEPS_PATH" deps
              mix deps.compile --no-deps-check --skip-umbrella-children

              runHook postConfigure
            '';

            buildPhase = ''
              runHook preBuild
              mix compile --no-deps-check
              mix build
              runHook postBuild
            '';

            installPhase = ''
              runHook preInstall
              mkdir -p "$out"
              cp -R output/. "$out/"
              runHook postInstall
            '';
          };
        }
      );

      devShells = forAllSystems (
        system:
        let
          pkgs = import nixpkgs { inherit system; };
        in
        {
          default = pkgs.mkShell {
            packages = with pkgs.beamPackages; [
              elixir
              hex
              rebar
              rebar3
            ];
          };
        }
      );
    };
}
