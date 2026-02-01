git_path="./.git"

if [ ! -d "$git_path" ]; then
    git init
fi

git add .
git ls-files '*.go' '*.sh' '*.shell' '*.bash' '*.luau' '*.squishy' | xargs wc -l

rm -rf $git_path