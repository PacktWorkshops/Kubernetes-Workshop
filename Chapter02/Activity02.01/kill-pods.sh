#!/usr/bin/env bash

# Sleep 1 second in each cycle.
SLEEP_INTERVAL=1

while true; do
    # Always keep one Pod running.
    running_pod_kept=
    # In each cycle, get the running Pods, and fetch their names.
    kubectl get pod -l tier=frontend --no-headers | grep Running | awk '{print $1}' | while read pod_name; do
        # Keep the 1st Pod running.
        if [[ -z $running_pod_kept ]]; then
            running_pod_kept=yes
            echo "Keeping Pod $pod_name running"
        # Delete all other running Pods.
        else
            echo "Killing Pod $pod_name"
            kubectl delete pod $pod_name 
        fi
    done
    sleep $SLEEP_INTERVAL
done
