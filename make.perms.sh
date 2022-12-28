#!/usr/bin/env bash
#------------------------------------------------------------------------------
#  Makefile : make perms 
# -----------------------------------------------------------------------------

perms(){
	echo "=== chmod 0600 @ all FILEs under '$1'"
	find "$1" -type f -execdir /bin/bash -c 'chmod 0600 "$@"' _ {} \+
}

perms infra
perms assets
chmod 0700 assets/HOME/.bin/ec2*

echo "=== chmod 0660 @ all *.{md,html,png,...} FILEs under PWD"
find . -type f \( -iname '*.md' -or -iname '*.png' -or -name '*.lnk' \
    -or -iname '*.html' -or -iname '*.json' \
    -or -iname '*.yml' -or -iname '*.log' \
    -or -iname '*.ico' -or -iname '*.ini' \
    -or -iname '*.webp' -or -iname '*.out' \
    -or -iname '*.txt' -or -iname '*.svg' \
    -or -iname '*.inf' -or -iname '*.jpg' \) \
    -execdir /bin/bash -c 'chmod 0660 "$@"' _ {} \+ &

echo "=== chmod 0444 @ all 'LICENSE' files under PWD"
find . -type f -iname 'LICENSE' -execdir /bin/bash -c 'chmod 0444 "$@"' _ {} \+ &

echo "=== chmod 0774 @ all *.shm *.go FILEs under PWD"
find . -type f \( -iname '*.sh' -or -iname '*.go' \) \
    -execdir /bin/bash -c 'chmod 0774 "$@"' _ {} \+

echo '=== chmod 0755 @ all DIRs under PWD'
find . -type d -execdir /bin/bash -c 'chmod 0755 "$@"' _ {} \+ &

sleep 2
# Wait until all (background) proceses (@ '-execdir') complete
while [[ $(ps aux |grep -- -execdir |grep -v grep |awk 'NR == 1 {print $2}') ]]; do sleep 2; done

ls -ahl --color=auto --group-directories-first . 
ls -ahlR --color=auto --group-directories-first app/pwa
ls -ahl --color=auto --group-directories-first assets/sql

exit 0

