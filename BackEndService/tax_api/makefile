codegen:
	oapi-codegen \
	-generate gin,spec \
	-package api config/taxarific_users.yml > api/api.gen.go
	oapi-codegen \
	-generate types \
	-package models config/taxarific_users.yml > models/models.gen.go
	go mod tidy
