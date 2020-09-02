#!/bin/bash
while [ true ]; do

NUM=$[ ( $RANDOM % 10 ) ]
kubectl get pod -n kubernetes-dashboard | awk -v num=$NUM '(NR==num && NR > 1) {print $1}' | xargs kubectl delete pod -n kubernetes-dashboard
sleep "$NUM"s

done
