brew cask install minikube
#minikube config set vm-driver vmware
minikube config set vm-driver vmwarefusion
minikube config set memory 8192
minikube config set cpus 4

minikube start --cpus=4 --memory=4096 --disk-size=30g --vm-driver=hyperkit --kubernetes-version=v1.10.12 --extra-config=apiserver.authorization-mode=RBAC -p new



#############
minikube start

kubectl config set-context confgimap-context --namespace=configmap-test
kubectl config use-context confgimap-context

kubectl create namespace configmap-test

kubectl create configmap singlevalue-map --from-literal=partner-url=https://www.auppost.com.au


cat > application.properties << EOF
partner-url=https://www.fedex.com
partner-key=1234
EOF


kubectl get configmaps --namespace configmap-test

kubectl create configmap full-file-map --from-file=./application.properties --namespace configmap-test


kubectl get configmaps --namespace configmap-test

kubectl get configmap singlevalue-map -o yaml --namespace configmap-test
kubectl get configmap full-file-map -o yaml --namespace configmap-test


kubectl create -f mount-configmap-as-volume.yaml --namespace configmap-test
kubectl logs -f configmap-test-pod --namespace configmap-test

kubectl create -f mount-configmap-as-env.yaml --namespace configmap-test
kubectl logs -f configmap-test-pod --namespace configmap-test



kubectl create secret --help 

kubectl create secret docker-registry test-docker-registry-secret --docker-username=test --docker-password=testpassword --docker-email=example@a.com --namespace configmap-test

kubectl get secrets test-docker-registry-secret --namespace configmap-test

kubectl describe secrets test-docker-registry-secret --namespace configmap-test