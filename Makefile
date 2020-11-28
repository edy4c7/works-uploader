all:	
	cd web \
	&& yarn \
	&& yarn generate
	cd functions \
	&& yarn \
	&& yarn build

deploy: all
	firebase deploy
