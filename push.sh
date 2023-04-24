info=$1
git add .
git commit -m  $info
git push origin main

git tag $info
git push origin $info
