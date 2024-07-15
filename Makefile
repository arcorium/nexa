
PRIVATE_KEY=privkey.pem
PUBLIC_KEY=pubkey.pem

prepare: create.key distribute.key

create.key:
	openssl genpkey -algorithm RSA -out $(PRIVATE_KEY)
	openssl rsa -pubout -in $(PRIVATE_KEY) -out $(PUBLIC_KEY)

distribute.key:
	@echo "Distribute public key"
	@for svc in ./services/*; do \
		if [ -d $$svc ]; then \
			cp $(PUBLIC_KEY) $$svc/; \
			echo "Public key copied to $$svc"; \
	  	fi \
	done
	@rm $(PUBLIC_KEY)
	@echo "Distribute private key"
	@cp $(PRIVATE_KEY) ./services/authentication/
	@cp $(PRIVATE_KEY) ./token_generator/
	@rm $(PRIVATE_KEY)
	@echo "Distributing done"

create.service:
	mkdir services\$(service)
	mkdir services\$(service)\cmd
	mkdir services\$(service)\config
	mkdir services\$(service)\constant
	mkdir services\$(service)\cmd\server
	mkdir services\$(service)\cmd\migrate
	mkdir services\$(service)\cmd\seed
	mkdir services\$(service)\internal
	mkdir services\$(service)\internal\api
	mkdir services\$(service)\internal\api\grpc\handler
	mkdir services\$(service)\internal\api\grpc\mapper
	mkdir services\$(service)\internal\api\grpc\interceptor
	mkdir services\$(service)\internal\app\service
	mkdir services\$(service)\internal\domain
	mkdir services\$(service)\internal\domain\dto
	mkdir services\$(service)\internal\domain\mapper
	mkdir services\$(service)\internal\domain\entity
	mkdir services\$(service)\internal\domain\service
	mkdir services\$(service)\internal\domain\repository
	mkdir services\$(service)\internal\infra
	mkdir services\$(service)\internal\infra\repository
	mkdir services\$(service)\util

gen.go:
	buf generate

run.compose:
	NEXA_RELEASE=1 SMTP_PATH=smtp.env \
	docker compose -f ./docker-compose.yml -f ./docker-compose.override.yml \
    -p nexa up -d

run.compose.dev:
	NEXA_RELEASE=0 SMTP_PATH=services/mailer/dev.smtp.env \
	docker compose -f ./docker-compose.yml -f ./docker-compose.override.yml \
        -p nexa up -d

stop.compose:
	NEXA_RELEASE=0 SMTP_PATH=services/mailer/dev.smtp.env \
	docker compose down --remove-orphans

run: prepare run.compose

run.dev: prepare run.compose.dev