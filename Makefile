run: shutdown
	docker run -d --name etcd -p 2379:2379 -p 2380:2380 quay.io/coreos/etcd:v3.0.15 etcd --listen-client-urls http://0.0.0.0:2379 --advertise-client-urls http://0.0.0.0:2380
	sleep 0
	ETCD_ENDPOINT=$$(docker-machine ip):2379 go run main.go

shutdown:
	docker stop etcd
	docker rm etcd