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

build: change-lang build-base
	ISHOCON_APP_LANG=$(ISHOCON_APP_LANG:python)
	docker build \
	--build-arg BASE_IMAGE=$(LOCAL_ISHOCON_BASE_IMAGE) \
	-f ./docker/app/$(ISHOCON_APP_LANG)/Dockerfile \
	-t ishocon1-app-$(ISHOCON_APP_LANG):latest \
	-t $(UNAME)/ishocon1-app-$(ISHOCON_APP_LANG):latest \
	.;
	@echo "Build done."

pull-base:
	docker pull $(UNAME)/ishocon1-app-base:latest;
	docker tag $(UNAME)/ishocon1-app-base:latest $(LOCAL_ISHOCON_BASE_IMAGE);

pull-app: check-lang
	docker pull $(UNAME)/ishocon1-app-$(ISHOCON_APP_LANG):latest;
	docker tag $(UNAME)/ishocon1-app-$(ISHOCON_APP_LANG):latest ishocon1-app-$(ISHOCON_APP_LANG):latest;

pull: pull-base pull-app
	@echo "Pull done."

push:
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

bench-from-scratch:
	docker compose up -d;
	sleep 30;
	docker exec -i ishocon1-app-1 sh -c "./benchmark --workload ${WORKLOAD}"

check-lang:
	if echo "$(ISHOCON_APP_LANG)" | grep -qE '^(ruby|python|go|nodejs)$$'; then \
        echo "ISHOCON_APP_LANG is valid."; \
    else \
        echo "Invalid ISHOCON_APP_LANG. It must be one of: ruby, python, go, nodejs."; \
        exit 1; \
    fi;

change-lang: check-lang
	if sed --version 2>&1 | grep -q GNU; then \
		echo "GNU sed"; \
		sed -i 's/\(ruby\|python\|go\|nodejs\)/'"$(ISHOCON_APP_LANG)"'/g' ./docker-compose.yml; \
	else \
		echo "BSD sed"; \
		sed -i '' -E 's/(ruby|python|go|nodejs)/'"$(ISHOCON_APP_LANG)"'/g' ./docker-compose.yml; \
	fi;
