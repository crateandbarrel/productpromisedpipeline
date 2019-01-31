function error_exit
{
        echo "$1" 1>&2
        exit 1
}

elbname="crateAPIE-PrivateE-1CX8B6DKCL82E"
asgname="crateAPIExternalPartners-DockerServerScalingGroup-16CUNZKSZBO28"

#assumes there are two active servers that need to be replaced, one at a time...

#clear local state
~/go/bin/awsdeploy -elb-name $elbname -task clearAndRefreshStatus -no-prompt true

#display current state
~/go/bin/awsdeploy -elb-name $elbname -task pollAndWaitUntilStable -no-prompt true

#increase capacity above the desired capacity to allow for removing and adding servers
if ~/go/bin/awsdeploy -elb-name $elbname -asg-name $asgname  -task increaseCapacity -desired-capacity 2 -incr-desired-capacity 1 -no-prompt true
then
echo "waiting 30 seconds to let the capacity to increase to trigger"
sleep 30s
#poll until new server is up
if ~/go/bin/awsdeploy -elb-name $elbname -task pollAndWaitUntilStable -no-prompt true -polling-wait-time 60
then
#remove first server (wait for comfirmation) - abort whole script if no confirmation is given
if ~/go/bin/awsdeploy -elb-name $elbname -asg-name $asgname  -task removeServers -elb-removal-incr 1
then
echo "waiting 30 seconds to let the first server trigger"
sleep 30s
#poll until first new server is up
if ~/go/bin/awsdeploy -elb-name $elbname -task pollAndWaitUntilStable -no-prompt true -polling-wait-time 60
then
#remove second server (wait for comfirmation) - abort whole script if no confirmation is given
if ~/go/bin/awsdeploy -elb-name $elbname -asg-name $asgname  -task removeServers -elb-removal-incr 1
then
echo "waiting 30 seconds to let the second server trigger"
sleep 30s
#poll until second new server is up
if ~/go/bin/awsdeploy -elb-name $elbname -task pollAndWaitUntilStable -no-prompt true -polling-wait-time 60
then
#decrease capacity back to original desired capacity
if ~/go/bin/awsdeploy -elb-name $elbname -asg-name $asgname  -task decreaseCapacity -desired-capacity 2 -no-prompt true
then
  #wait 60 seconds to let the autoscaler terminate a server
  echo "waiting 60 seconds to let autoscaler decrease servers"
  sleep 60s
  #display final status with what should be 3 servers with PendingTerm status
  ~/go/bin/awsdeploy -elb-name $elbname -task pollAndWaitUntilStable -no-prompt true
  echo "Finished Successfully - don't forget to terminate any unused EC2 servers"
  echo "EXAMPLE:  ~/go/bin/awsdeploy -elb-name $elbname -task terminateAPendingTermServer"
else
  error_exit "Aborting decreasing capacity"
fi
else
  error_exit "Aborting polling waiting on second server"
fi
else
  error_exit "Aborting second server replace"
fi
else
  error_exit "Aborting polling waiting on first server"
fi
else
  error_exit "Aborting first server replace"
fi
else
  error_exit "Aborting polling waiting on capacity to increase"
fi
else
  error_exit "Aborting error increasing capacity"
fi