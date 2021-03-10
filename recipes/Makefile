$(SPECIALRESOURCE):
	kubectl apply -f recipes/$(SPECIALRESOURCE)/0000-$(SPECIALRESOURCE)-cr.yaml

assets:
	cd recipes/$(SPECIALRESOURCE)/manifests && $(KUSTOMIZE) edit set namespace $(SPECIALRESOURCE)
	kubectl create ns $(SPECIALRESOURCE) --dry-run=client -o yaml | kubectl apply -f -
	kubectl apply -k recipes/$(SPECIALRESOURCE)/manifests
