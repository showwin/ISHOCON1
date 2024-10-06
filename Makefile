WORKLOAD = 3
ifeq ($(UNAME),)
	UNAME = $(shell whoami)
endif

ifeq ($(ARCH),)
	ARCH = $(shell uname -m)
endif

LOCAL_ISHOCON_BASE_IMAGE = ishocon1-app-base:latest

build-base:
	docker build \
	-f ./docker/app/base/Dockerfile \
	-t $(LOCAL_ISHOCON_BASE_IMAGE) \
	-t $(UNAME)/ishocon1-app-base:latest \
	.;

build-bench:
	docker build \
	-f ./docker/benchmarker/Dockerfile \
	-t ishocon1-bench:latest \
	-t $(UNAME)/ishocon1-bench:latest \
	.;

build-app: change-lang build-base
	ISHOCON_APP_LANG=$(ISHOCON_APP_LANG:python)
	docker build \
	--build-arg BASE_IMAGE=$(LOCAL_ISHOCON_BASE_IMAGE) \
	-f ./docker/app/$(ISHOCON_APP_LANG)/Dockerfile \
	-t ishocon1-app-$(ISHOCON_APP_LANG):latest \
	-t $(UNAME)/ishocon1-app-$(ISHOCON_APP_LANG):latest \
	.;

build: build-bench build-app
	@echo "Build done."

pull-bench:
	docker pull $(UNAME)/ishocon1-bench:latest;
	docker tag $(UNAME)/ishocon1-bench:latest ishocon1-bench:latest;

pull-base:
	docker pull $(UNAME)/ishocon1-app-base:latest;
	docker tag $(UNAME)/ishocon1-app-base:latest $(LOCAL_ISHOCON_BASE_IMAGE);

pull-app: check-lang
	docker pull $(UNAME)/ishocon1-app-$(ISHOCON_APP_LANG):latest;
	docker tag $(UNAME)/ishocon1-app-$(ISHOCON_APP_LANG):latest ishocon1-app-$(ISHOCON_APP_LANG):latest;

pull: pull-bench pull-base pull-app
	@echo "Pull done."

push:
	docker push $(UNAME)/ishocon1-bench:latest;
	docker push $(UNAME)/ishocon1-app-base:latest;
	docker push $(UNAME)/ishocon1-app-$(ISHOCON_APP_LANG):latest;

up:
	docker compose up -d;

up-nod:
	docker compose up;

down:
	docker compose down;

bench:
	docker exec -i ishocon1-app-1 sh -c "./benchmark --workload ${WORKLOAD}"

bench-with-db-init:
	docker exec -i ishocon1-bench-1 sh -c " \
		service mysql restart && \
		tar -zxvf ~/admin/ishocon1.dump.tar.gz -C ~/admin && mysql -u root -pishocon ishocon1 < ~/admin/ishocon1.dump && \
		./benchmark --ip app:80 --workload ${WORKLOAD} \
	";

check-lang:
	if echo "$(ISHOCON_APP_LANG)" | grep -qE '^(ruby|python|go)$$'; then \
        echo "ISHOCON_APP_LANG is valid."; \
    else \
        echo "Invalid ISHOCON_APP_LANG. It must be one of: ruby, python, go."; \
        exit 1; \
    fi;

change-lang: check-lang
	if sed --version 2>&1 | grep -q GNU; then \
		echo "GNU sed"; \
		sed -i 's/\(ruby\|python\|go\)/'"$(ISHOCON_APP_LANG)"'/g' ./docker-compose.yml; \
	else \
		echo "BSD sed"; \
		sed -i '' -E 's/(ruby|python|go)/'"$(ISHOCON_APP_LANG)"'/g' ./docker-compose.yml; \
	fi;
