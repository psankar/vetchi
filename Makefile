dev:
	tilt down
	helm repo add cnpg https://cloudnative-pg.github.io/charts
	# Install CloudNativePG operator first
	helm upgrade --install cnpg \
		--namespace cnpg-system \
		--create-namespace \
		cnpg/cloudnative-pg
	# Wait for operator to be ready
	kubectl -n cnpg-system wait --for=condition=ready pod -l app.kubernetes.io/name=cloudnative-pg --timeout=60s
	# Now install our infrastructure
	cd charts/vetchidev && helm dependency build && cd ../..
	helm uninstall vetchidev -n vetchidev --ignore-not-found
	time kubectl delete pvc -n vetchidev --all --ignore-not-found
	kubectl delete pv -n vetchidev --all --ignore-not-found
	kubectl delete namespace vetchidev --ignore-not-found --force --grace-period=0
	kubectl create namespace vetchidev
	helm upgrade --install vetchidev ./charts/vetchidev -n vetchidev
	sleep 5
	kubectl -n vetchidev wait --for=condition=ready pod -l cnpg.io/podRole=instance --timeout=240s
	kubectl -n vetchidev port-forward service/postgres-rw 5432:5432 &
	kubectl -n vetchidev port-forward svc/mailpit 8025:8025 &
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
