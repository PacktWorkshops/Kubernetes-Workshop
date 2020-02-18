build:
	cd with_redis && docker build -t quay.io/zparnold/packt-kubernetes-demo-app:with_redis_v1 .
	cd without_redis && docker build -t quay.io/zparnold/packt-kubernetes-demo-app:without_redis_v1 .
	docker push quay.io/zparnold/packt-kubernetes-demo-app:with_redis_v1
	docker push quay.io/zparnold/packt-kubernetes-demo-app:without_redis_v1