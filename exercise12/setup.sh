#!/bin/bash
# sets up the current directory for testing exercise 12

if [[ -d "./sample/old" ]]; then
  rm -rf "./sample/old"
fi

mkdir -p "./sample/old";
cd "sample" || exit 1

for i in {1..4}
do
  touch "birthday_00" + $i + ".txt"
  touch "christmas_2016_" + $i + "_of_100.txt"
done

cd "old" || exit 1
for i in {5..8}
do
  touch "n_00" + $i + ".txt"
  touch "birthday_00" + $i + ".txt"
done