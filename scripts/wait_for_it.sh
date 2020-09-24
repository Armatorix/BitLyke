#/bin/bash

max_retry=10
counter=0
until [ -n $(curl -sf ${1}) ]
do
   sleep 2
   [[ counter -eq $max_retry ]] && echo "Failed all api checks!" && exit 1
   echo "Trying again. Try #$counter"
   ((counter++))
done

echo "Health check ready"
echo "API started"