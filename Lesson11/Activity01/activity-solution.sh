#Create the test namespace
kubectl create namespace my-app-test

#Create the a folder for test namespace
mkdir -p activity-configmap-secret/test
cd activity-configmap-secret/test

# create application.properties file
cat > application-data.properties << EOF
external-system-location=https://testvendor.example.com
external-system-basic-auth-username=user123
EOF

#create the configmap in the test namespace
kubectl create configmap my-app-data --from-file=./application-data.properties --namespace my-app-test


# create the file that contains passwoed
cat > application-secure.properties << EOF
external-system-basic-auth-password=password123
EOF

#create the secret in the test namespace
kubectl create secret generic my-app-secret --from-file=application-secure.properties --namespace my-app-test

#run the pod
kubectl create -f activity-solution-pod-spec.yaml --namespace my-app-test

#see the logs
kubectl logs -f configmap-secrets-activity-pod --namespace my-app-test


=======
#Create the test namespace
kubectl create namespace my-app-prod

#Create the a folder for prod namespace
mkdir -p activity-configmap-secret/prod
cd activity-configmap-secret/prod


# create application.properties file
cat > application-data.properties << EOF
external-system-location=https://vendor.example.com
external-system-basic-auth-username=activityapplicationuser
EOF

#create the configmap in the prod namespace
kubectl create configmap my-app-data --from-file=./application-data.properties --namespace my-app-prod

# create the file that contains passwoed
cat > application-secure.properties << EOF
external-system-basic-auth-password=A#4b*(1=B88%tFr3
EOF

#create the secret in the prod namespace
kubectl create secret generic my-app-secret --from-file=application-secure.properties --namespace my-app-prod

#run the pod
kubectl create -f activity-solution-pod-spec.yaml --namespace my-app-prod

#see the logs
kubectl logs -f configmap-secrets-activity-pod --namespace my-app-prod


