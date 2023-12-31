NOW := $(shell echo "`date +%Y-%m-%d`")

#
# Display help
# 
define help_info
	@echo "\nUsage:\n"
	@echo ""
	@echo "  $$ make setVersion version=9.9.9     - Used to set the version number."
	@echo ""


endef

help:
	$(call help_info)


setVersion:
	@echo Set version "$(version)";\
	ver="$(version)" yq e '.image.tag = strenv(ver)' ./charts/serial-gateway/values.yaml --inplace;\
	ver="v$(version)" yq e '.version = strenv(ver)' ./charts/serial-gateway/Chart.yaml  --inplace;\
	ver="v$(version)" yq e '.appVersion = strenv(ver)' ./charts/serial-gateway/Chart.yaml  --inplace;


