variable "cluster-name" {
  default = "terraform-eks-demo"
  type = string
}

# This data source is included for ease of sample architecture deployment
# and can be swapped out as necessary.
data "aws_availability_zones" "available" {}

resource "aws_vpc" "demo" {
  cidr_block = "10.0.0.0/16"
  enable_dns_hostnames = true
  tags = map(
  "Name", "terraform-eks-demo-node",
  "kubernetes.io/cluster/${var.cluster-name}", "shared",
  )
}

resource "aws_subnet" "demo" {
  count = 3

  availability_zone = data.aws_availability_zones.available.names[count.index]
  cidr_block = "10.0.${count.index}.0/24"
  vpc_id = aws_vpc.demo.id

  tags = map(
  "Name", "terraform-eks-demo-node",
  "kubernetes.io/cluster/${var.cluster-name}", "shared",
  )
}

resource "aws_internet_gateway" "demo" {
  vpc_id = aws_vpc.demo.id

  tags = {
    Name = "terraform-eks-demo"
  }
}

resource "aws_route_table" "demo" {
  vpc_id = aws_vpc.demo.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.demo.id
  }
}

resource "aws_route_table_association" "demo" {
  count = 3

  subnet_id = aws_subnet.demo.*.id[count.index]
  route_table_id = aws_route_table.demo.id
}

resource "aws_iam_role" "demo-cluster" {
  name = "terraform-eks-demo-cluster"

  assume_role_policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "eks.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
POLICY
}

resource "aws_iam_role_policy_attachment" "demo-cluster-AmazonEKSClusterPolicy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
  role = aws_iam_role.demo-cluster.name
}

resource "aws_iam_role_policy_attachment" "demo-cluster-AmazonEKSServicePolicy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSServicePolicy"
  role = aws_iam_role.demo-cluster.name
}

resource "aws_security_group" "demo-cluster" {
  name = "terraform-eks-demo-cluster"
  description = "Cluster communication with worker nodes"
  vpc_id = aws_vpc.demo.id

  egress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = [
      "0.0.0.0/0"]
  }

  tags = {
    Name = "terraform-eks-demo"
  }
}

resource "aws_security_group_rule" "demo-cluster-ingress-workstation-https" {
  cidr_blocks = ["0.0.0.0/0"]
  description = "Allow workstation to communicate with the cluster API Server"
  from_port = 443
  protocol = "tcp"
  security_group_id = aws_security_group.demo-cluster.id
  to_port = 443
  type = "ingress"
}

resource "aws_eks_cluster" "demo" {
  name = var.cluster-name
  role_arn = aws_iam_role.demo-cluster.arn

  vpc_config {
    security_group_ids = [aws_security_group.demo-cluster.id]
    subnet_ids = aws_subnet.demo.*.id
  }

  depends_on = [
    aws_iam_role_policy_attachment.demo-cluster-AmazonEKSClusterPolicy,
    aws_iam_role_policy_attachment.demo-cluster-AmazonEKSServicePolicy,
  ]
}

locals {
  kubeconfig = <<KUBECONFIG


apiVersion: v1
clusters:
- cluster:
    server: ${aws_eks_cluster.demo.endpoint}
    certificate-authority-data: ${aws_eks_cluster.demo.certificate_authority.0.data}
  name: kubernetes
contexts:
- context:
    cluster: kubernetes
    user: aws
  name: aws
current-context: aws
kind: Config
preferences: {}
users:
- name: aws
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1alpha1
      command: aws-iam-authenticator
      args:
        - "token"
        - "-i"
        - "${var.cluster-name}"
KUBECONFIG
}

output "kubeconfig" {
  value = local.kubeconfig
}

resource "aws_iam_role" "demo-node" {
  name = "terraform-eks-demo-node"

  assume_role_policy = <<POLICY
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
POLICY
}

resource "aws_iam_role_policy_attachment" "demo-node-AmazonEKSWorkerNodePolicy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"
  role = aws_iam_role.demo-node.name
}

resource "aws_iam_role_policy_attachment" "demo-node-AmazonEKS_CNI_Policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"
  role = aws_iam_role.demo-node.name
}

resource "aws_iam_role_policy_attachment" "demo-node-AmazonEC2ContainerRegistryReadOnly" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
  role = aws_iam_role.demo-node.name
}

resource "aws_iam_instance_profile" "demo-node" {
  name = "terraform-eks-demo"
  role = aws_iam_role.demo-node.name
}

resource "aws_security_group" "demo-node" {
  name = "terraform-eks-demo-node"
  description = "Security group for all nodes in the cluster"
  vpc_id = aws_vpc.demo.id

  egress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = [
      "0.0.0.0/0"]
  }

  tags = map(
     "Name", "terraform-eks-demo-node",
     "kubernetes.io/cluster/${var.cluster-name}", "owned",
    )
}

resource "aws_security_group_rule" "demo-node-ingress-self" {
  description = "Allow node to communicate with each other"
  from_port = 0
  protocol = "-1"
  security_group_id = aws_security_group.demo-node.id
  source_security_group_id = aws_security_group.demo-node.id
  to_port = 65535
  type = "ingress"
}

resource "aws_security_group_rule" "demo-node-ingress-cluster" {
  description = "Allow worker Kubelets and pods to receive communication from the cluster control plane"
  from_port = 1025
  protocol = "tcp"
  security_group_id = aws_security_group.demo-node.id
  source_security_group_id = aws_security_group.demo-cluster.id
  to_port = 65535
  type = "ingress"
}

resource "aws_security_group_rule" "demo-cluster-ingress-node-https" {
  description = "Allow pods to communicate with the cluster API Server"
  from_port = 443
  protocol = "tcp"
  security_group_id = aws_security_group.demo-node.id
  source_security_group_id = aws_security_group.demo-cluster.id
  to_port = 443
  type = "ingress"
}

data "aws_ami" "eks-worker" {
  filter {
    name = "name"
    values = [
      "amazon-eks-node-${aws_eks_cluster.demo.version}-v*"]
  }

  most_recent = true
  owners = [
    "602401143452"]
  # Amazon EKS AMI Account ID
}

# This data source is included for ease of sample architecture deployment
# and can be swapped out as necessary.
data "aws_region" "current" {}

# EKS currently documents this required userdata for EKS worker nodes to
# properly configure Kubernetes applications on the EC2 instance.
# We implement a Terraform local here to simplify Base64 encoding this
# information into the AutoScaling Launch Configuration.
# More information: https://docs.aws.amazon.com/eks/latest/userguide/launch-workers.html
locals {
  demo-node-userdata = <<USERDATA
#!/bin/bash
set -o xtrace
/etc/eks/bootstrap.sh --apiserver-endpoint '${aws_eks_cluster.demo.endpoint}' --b64-cluster-ca '${aws_eks_cluster.demo.certificate_authority.0.data}' '${var.cluster-name}'
USERDATA
}

resource "aws_launch_configuration" "demo" {
  associate_public_ip_address = true
  iam_instance_profile = aws_iam_instance_profile.demo-node.name
  image_id = data.aws_ami.eks-worker.id
  instance_type = "m4.large"
  name_prefix = "terraform-eks-demo"
  security_groups = [
    aws_security_group.demo-node.id]
  user_data_base64 = base64encode(local.demo-node-userdata)

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_autoscaling_group" "demo" {
  desired_capacity = 2
  launch_configuration = aws_launch_configuration.demo.id
  max_size = 2
  min_size = 1
  name = "terraform-eks-demo"
  vpc_zone_identifier = aws_subnet.demo.*.id

  tag {
    key = "Name"
    value = "terraform-eks-demo"
    propagate_at_launch = true
  }

  tag {
    key = "kubernetes.io/cluster/${var.cluster-name}"
    value = "owned"
    propagate_at_launch = true
  }
}

locals {
  config_map_aws_auth = <<CONFIGMAPAWSAUTH


apiVersion: v1
kind: ConfigMap
metadata:
  name: aws-auth
  namespace: kube-system
data:
  mapRoles: |
    - rolearn: ${aws_iam_role.demo-node.arn}
      username: system:node:{{EC2PrivateDNSName}}
      groups:
        - system:bootstrappers
        - system:nodes
CONFIGMAPAWSAUTH
}

output "config_map_aws_auth" {
  value = local.config_map_aws_auth
}