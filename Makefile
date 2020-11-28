all:	
	cd web \
	&& yarn \
	&& yarn generate
	cd functions \
	&& yarn \
	&& yarn run build
