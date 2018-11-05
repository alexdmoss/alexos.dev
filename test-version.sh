#!/bin/bash

# hostname="http://35.187.52.146"
hostname="https://www.mosstech.io"

i=1
one=0
two=0 


while [[ $i -lt 21 ]]; do

  echo "Iteration $i"
  version=$(curl -s ${hostname}/version.html)

  if [[ $version == "1.0" ]]; then
    one=$(expr $one + 1)
  elif [[ $version == "2.0" ]]; then
    two=$(expr $two + 1)
  else
    echo "ERROR - no version!"
  fi

  i=$(expr $i + 1)

done

echo
echo "-----------------------------------"
echo " Version 1.0: $one"
echo " Version 2.0: $two"
echo "-----------------------------------"
echo
