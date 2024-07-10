
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
	go mod init -C ./services/$(service) nexa/services/$(service)
	mkdir services\$(service)\cmd
	mkdir services\$(service)\internal
	mkdir services\$(service)\internal\api
	mkdir services\$(service)\internal\app\config
	mkdir services\$(service)\internal\app\service
	mkdir services\$(service)\internal\domain
	mkdir services\$(service)\internal\domain\service
	mkdir services\$(service)\internal\domain\repository
	mkdir services\$(service)\internal\infra
	mkdir services\$(service)\internal\infra\model
	mkdir services\$(service)\internal\infra\repository
	mkdir services\$(service)\shared
	mkdir services\$(service)\shared\domain
	mkdir services\$(service)\shared\domain\dto
	mkdir services\$(service)\shared\domain\entity
	mkdir services\$(service)\shared\domain\mapper
	mkdir services\$(service)\shared\proto
	mkdir services\$(service)\test
	go work use ./services/$(service)

gen.go:
	buf generate