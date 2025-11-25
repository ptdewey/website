{ pkgs, cedarDrv }:

pkgs.writeShellApplication {
  name = "run";

  runtimeInputs = [ pkgs.tailwindcss_4 cedarDrv ];

  text = ''
    #!/usr/bin/env bash
    set -e
    echo "building site static files..."
    cedar
    echo "done."
    echo "building tailwind styles..."
    tailwindcss -i static/app.css -o public/style.css -m
    echo "done."
  '';
}
