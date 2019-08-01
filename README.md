[![GitHub issues](https://img.shields.io/github/issues/TrainingByPackt/Kubernetes-Workshop.svg)](https://github.com/TrainingByPackt/Kubernetes-Workshop/issues)
[![GitHub forks](https://img.shields.io/github/forks/TrainingByPackt/Kubernetes-Workshop)](https://github.com/TrainingByPackt/Kubernetes-Workshop/network)
[![GitHub stars](https://img.shields.io/github/stars/TrainingByPackt/Kubernetes-Workshop.svg)](https://github.com/TrainingByPackt/Kubernetes-Workshop/stargazers)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/TrainingByPackt/Kubernetes-Workshop/pulls)

# Kubernetes-Workshop
“Intro to Kubernetes” will cover a wide variety of topics surrounding the Kubernetes ecosystem. Because Kubernetes is a dense (intellectually speaking) tool with quite a few moving parts, we are going to survey the landscape in which Kubernetes operates in addition to covering the operation of clusters and the applications they contain. We are choosing to do this because we believe it will aid the reader in understanding some of the purpose built design assumptions and decisions that are part of the Kubernetes core API and featureset. In addition to highlighting aspects of the ecosystem, we will also cover the overview of a Kubernetes-oriented software delivery pipeline, provisioning and managing a cluster (from the hardware and software perspective) and finally cover running applications inside of Kubernetes architected for scale and stability.

# What you will learn:
* Will be able to create a container image from an image definition manifest (such as a Dockerfile)
* Will be able to push that image to a registry such as Docker Hub, Quay.io, or ECR/GCR
* Potentially give advice for containerizing popular programming frameworks
* Understand the locus of control shift from you as the operator to you as the person who instructs the operator (benevolent botnet)
* Will be able to build a Kubernetes cluster hosted in a managed service provider (such as AKS, EKS, GKE, IKS, Openshift) or using an open source provisioning tool such as Kubeadm, Kops, or Kubespray
* Understand what the following terms mean and how to use them to construct an orchestrated application running atop Kubernetes:
* CustomResourceDefinition, AdmissionController, Deployment, ReplicaSet, DaemonSet, Pod, Service, Secret, ConfigMap, Namespace, ServiceAccount, ClusterRole, ClusterRoleBinding, Role, RoleBinding, NetworkPolicy, PodSecurityPolicy, Volume, PersistentVolume, PersistentVolumeClaim
* Give plenty of examples of YAML above
* Understand what the following building blocks of Kubernetes are and how to leverage them for benefit:
