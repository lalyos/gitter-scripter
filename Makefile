docker-img:
	docker build -t eu.gcr.io/involuted-smile-221511/gitter .

docker-push: docker-img
	docker push eu.gcr.io/involuted-smile-221511/gitter

docker-push-local:
	GOOS=linux go build -v -o gitter-scripter-linux .
	docker build -f Dockerfile.local -t eu.gcr.io/involuted-smile-221511/gitter .
	docker push eu.gcr.io/involuted-smile-221511/gitter

redeploy: docker-push-local 
	kubectl delete -f gitter.yaml
	kubectl apply -f gitter.yaml