{ pkgs, cedarDrv }:

pkgs.writeShellApplication {
  name = "run";

  runtimeInputs = [ pkgs.tailwindcss_4 cedarDrv pkgs.jujutsu ];

  text = ''
    #!/usr/bin/env bash
    set -e
    
    # Generate change info template using jj
    echo "generating change info..."
    CHANGE_ID_SHORT=$(jj log -r @ --no-graph -T 'change_id.short()' 2>/dev/null || echo "unknown")
    CHANGE_ID_FULL=$(jj log -r @ --no-graph -T 'change_id' 2>/dev/null || echo "unknown")
    
    cat > templates/partials/_git-info.tmpl << EOF
{{define "git-info"}}
<span class="git-commit" title="Full change ID: $CHANGE_ID_FULL">$CHANGE_ID_SHORT</span>
{{end}}
EOF
    
    echo "building site static files..."
    cedar
    echo "done."
    echo "building tailwind styles..."
    tailwindcss -i static/app.css -o public/style.css -m
    echo "done."
  '';
}
