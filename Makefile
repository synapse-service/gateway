
# include scripts/*.mk

.PHONY: run
run:
	@go run \
		cmd/main/main.go \
		-config config.yaml


.PHONY: proto
proto:
	@protoc \
		--proto_path api/ \
		--proto_path api/third/ \
		--go_out . \
		--go-grpc_out . \
		--govalidators_out . \
		--doc_out ./api \
		--doc_opt html,index.html \
		api/gateway/*.proto

ca:
	@rm -f keys/*.pem
	@make cert:ca cert:client cert:client:sign

cert\:ca:
	@openssl req \
		-x509 -nodes -new -sha256 \
		-days 365 \
		-newkey rsa:2048 \
    	-keyout keys/ca.key.pem \
		-out keys/ca.cert.pem \
		-addext "subjectAltName=DNS:0.0.0.0,DNS:localhost,IP:0.0.0.0" \
		-subj '/CN=0.0.0.0'

cert\:ca\:show:
	@openssl x509 \
		-text -noout \
		-in keys/ca.cert.pem

cert\:client:
	@openssl req \
		-nodes -new -sha256 \
		-newkey rsa:2048 \
		-keyout keys/client.key.pem \
		-out keys/client.req.pem \
		-addext "subjectAltName=DNS:0.0.0.0,DNS:localhost,IP:0.0.0.0" \
		-subj '/CN=0.0.0.0'

cert\:client\:sign:
	@openssl x509 \
		-req -sha256 \
		-days 365 \
		-in keys/client.req.pem \
		-signkey keys/ca.key.pem \
		-out keys/client.cert.pem
	@rm keys/client.req.pem

cert\:client\:show:
	@openssl x509 \
		-text -noout \
		-in keys/client.cert.pem
