#/bin/bash

for i in {1..5}
do
   nc -z localhost 8080
   if [ $? == 0 ]
   then
      echo "Connection ready"
      exit 0
   fi
   echo "Check #${i} failed"
   echo "Retry in 3s..."
   sleep 3
done 

echo "All checks failed."
exit 1 