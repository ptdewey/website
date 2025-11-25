{ pkgs }:

pkgs.writeShellApplication {
  name = "run";

  runtimeInputs = [ pkgs.tailwindcss_4 ];

  text = ''
    #!/usr/bin/env bash
    set -e
    echo "building tailwind styles..."
    tailwindcss -i static/app.css -o public/style.css -m
    echo "done."
  '';
}
