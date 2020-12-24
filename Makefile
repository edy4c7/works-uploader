.PHONY: all
all:	
	cd web \
	&& yarn \
	&& yarn generate
	cd functions \
	&& yarn \
	&& yarn build

emulator: all
	firebase emulators:start
	&& yarn run build

.PHONY: dev
dev:
	cd web && yarn dev

.PHONY: test
test:
	cd web && yarn test
