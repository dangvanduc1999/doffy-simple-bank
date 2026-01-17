#!/bin/bash

# This script calculates the next version based on PR title and creates a git tag.
set -e

PR_TITLE=$1

if [ -z "$PR_TITLE" ]; then
    echo "Error: PR title is required."
    echo "Usage: $0 \"<prefix> : <message>\""
    exit 1
fi

# Validate pattern: must contain a colon
if [[ ! "$PR_TITLE" == *":"* ]]; then
    echo "Error: PR title does NOT match pattern '<prefix> : <message>'"
    echo "PR Title: $PR_TITLE"
    exit 1
fi

# Extract prefix (everything before the first colon, trimmed)
PREFIX=$(echo "$PR_TITLE" | cut -d':' -f1 | xargs)
PREFIX=$(echo "$PREFIX" | tr '[:upper:]' '[:lower:]')

if [ -z "$PREFIX" ]; then
    echo "Error: Prefix is empty."
    echo "PR Title: $PR_TITLE"
    exit 1
fi

echo "Detected prefix: $PREFIX"

# Determine bump type based on user rules:
# feature, chore, patch -> patch
# fix -> minor
# major -> major
BUMP_TYPE=""

case "$PREFIX" in
    feature|chore|patch)
        BUMP_TYPE="patch"
        ;;
    fix)
        BUMP_TYPE="minor"
        ;;
    major)
        BUMP_TYPE="major"
        ;;
    *)
        echo "Error: Unknown prefix '$PREFIX'. Expected feature, chore, patch, fix, or major."
        exit 1
        ;;
esac

echo "Bump type: $BUMP_TYPE"

# Ensure we have the latest tags from main
git fetch --tags

# Get the latest tag
LATEST_TAG=$(git tag --sort=-v:refname | head -n 1)

if [ -z "$LATEST_TAG" ]; then
    echo "No tags found. Starting from v0.0.0"
    LATEST_TAG="v0.0.0"
fi

echo "Current version: $LATEST_TAG"

# Remove 'v' prefix if exists
VERSION=${LATEST_TAG#v}

# Parse version components
IFS='.' read -r MAJOR MINOR PATCH <<< "$VERSION"

# Default to 0 if components are missing
MAJOR=${MAJOR:-0}
MINOR=${MINOR:-0}
PATCH=${PATCH:-0}

# Increment version
if [ "$BUMP_TYPE" == "major" ]; then
    MAJOR=$((MAJOR + 1))
    MINOR=0
    PATCH=0
elif [ "$BUMP_TYPE" == "minor" ]; then
    MINOR=$((MINOR + 1))
    PATCH=0
elif [ "$BUMP_TYPE" == "patch" ]; then
    PATCH=$((PATCH + 1))
fi

NEW_VERSION="v$MAJOR.$MINOR.$PATCH"
echo "New version: $NEW_VERSION"

# Create and push the tag
echo "Creating tag $NEW_VERSION..."
git tag "$NEW_VERSION"
git push origin "$NEW_VERSION"

echo "Successfully created and pushed tag $NEW_VERSION"

# Output for GitHub Actions
if [ -n "$GITHUB_OUTPUT" ]; then
    echo "new_tag=$NEW_VERSION" >> "$GITHUB_OUTPUT"
fi