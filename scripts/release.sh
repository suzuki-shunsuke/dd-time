#!/usr/bin/env bash
# Usage
#   bash scripts/tag.sh v0.3.2

REMOTE=https://github.com/suzuki-shunsuke/dd-time

echoEval() {
  echo "+ $*"
  eval "$@"
}

BRANCH=$(git branch | grep "^\* " | sed -e "s/^\* \(.*\)/\1/")
if [ "$BRANCH" != "master" ]; then
  read -r -p "The current branch isn't master but $BRANCH. Are you ok? (y/n)" YN
  if [ "${YN}" != "y" ]; then
    echo "cancel to release"
    exit 0
  fi
fi

TAG="$1"
echo "TAG: $TAG"
VERSION="${TAG#v}"

if [ "$TAG" = "$VERSION" ]; then
  echo "the tag must start with 'v'" >&2
  exit 1
fi

echoEval cd "$(dirname "$0")/.." || exit 1

VERSION_FILE=pkg/constant/version.go

echo "create $VERSION_FILE"
cat << EOS > "$VERSION_FILE" || exit 1
package constant

// Don't edit this file.
// This file is generated by the release command.

// Version is the dd-time's version.
const Version = "$VERSION"
EOS

echoEval git add "$VERSION_FILE" || exit 1
echo "+ git commit -m \"build: update version to $TAG\""
git commit -m "build: update version to $TAG" || exit 1
echoEval git tag "$TAG" || exit 1
echoEval git push "$REMOTE" "$BRANCH" || exit 1
echoEval git push "$REMOTE" "$TAG"
