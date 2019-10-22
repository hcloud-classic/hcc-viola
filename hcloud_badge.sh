#!/bin/bash

goreportcard-cli -v > hcloud-badge/goreport_$1
sed -n '/Grade/p' hcloud-badge/goreport_$1 > goreport_grade
if [ $(uname -s) == "FreeBSD" ]
then
  sed -i '' 's/Grade:\ //g' goreport_grade
  sed -i '' 's/\%/\%25/g' goreport_grade
  sed -i '' 's/\ /\%20/g' goreport_grade
  sed -i '' 's/\+/\%2B/g' goreport_grade
else
  sed -i 's/Grade:\ //g' goreport_grade
  sed -i 's/\%/\%25/g' goreport_grade
  sed -i 's/\ /\%20/g' goreport_grade
  sed -i 's/\+/\%2B/g' goreport_grade
fi
wget -O hcloud-badge/hcloud-badge_$1.svg 'https://img.shields.io/badge/go%20report-'$(cat goreport_grade)'-brightgreenv'

cd hcloud-badge/
git add .
git commit -m "Update goreport and badge for $1"
git push origin feature/dev
