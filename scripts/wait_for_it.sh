#/bin/bash

max_retry=6
count=0
until [ ! $(nc -z localhost 8080) ]; do
   sleep 2
   [ $count -eq $max_retry ] && echo "Failed all api checks!" && exit 1
   echo "Trying again. Try #${count}"
   ((count++))
done

echo "Health check ready"
echo "API started"