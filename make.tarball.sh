#!/usr/bin/env bash
#------------------------------------------------------------------------------
#  make : tarball
# -----------------------------------------------------------------------------

echo '=== Markdown => Markup'
make markup 
echo "=== tarball @ ${PWD}"
/c/HOME/.bin/tarball
export dname="$(basename "$PWD")"
pushd ./.. 
echo "=== move to dir @ '${PWD[@]:0:3}'"
mv "$(find . -iname "${dname}.*.tgz")" "${PWD[@]:0:3}"
popd

