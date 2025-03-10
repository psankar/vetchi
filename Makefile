dev:
	tilt down
	time kubectl delete pvc -n vetchidev --all --ignore-not-found
	kubectl delete pv -n vetchidev --all --ignore-not-found
	kubectl delete namespace vetchidev --ignore-not-found --force --grace-period=0
	kubectl create namespace vetchidev
	tilt up

test:
	@ORIG_URI=$$(kubectl -n vetchidev get secret postgres-app -o jsonpath='{.data.uri}' | base64 -d); \
	MOD_URI=$$(echo $$ORIG_URI | sed 's/postgres-rw.vetchidev/localhost/g'); \
	POSTGRES_URI=$$MOD_URI ginkgo -v ./dolores/...

seed:
	@ORIG_URI=$$(kubectl -n vetchidev get secret postgres-app -o jsonpath='{.data.uri}' | base64 -d); \
	MOD_URI=$$(echo $$ORIG_URI | sed 's/postgres-rw.vetchidev/localhost/g'); \
	cd dev-seed && POSTGRES_URI=$$MOD_URI go run .

lib:
	cd typespec && tsp compile . && npm run build && \
	cd ../harrypotter && npm install ../typespec && \
	cd ../ronweasly && npm install ../typespec

docker:
	$(eval GIT_SHA=$(shell git rev-parse --short HEAD))
	docker build -f Dockerfile-harrypotter -t vetchi/harrypotter:$(GIT_SHA) .
	docker build -f Dockerfile-ronweasly -t vetchi/ronweasly:$(GIT_SHA) .
	docker tag vetchi/harrypotter:$(GIT_SHA) vetchi/harrypotter:latest
	docker tag vetchi/ronweasly:$(GIT_SHA) vetchi/ronweasly:latest
