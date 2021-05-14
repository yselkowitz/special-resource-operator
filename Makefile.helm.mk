HELM_REPOS = $(shell ls -d charts/*/)


helm-lint:
	echo $(HELM_REPOS)
	@for repo in $(HELM_REPOS); do                          \
		cd $$repo;                              \
		helm lint -f ../global-values.yaml `ls -d */`; \
		cd ../..;                                      \
	done

helm-repo-index: helm-lint
	@for repo in $(HELM_REPOS); do                          \
		cd $$repo;                              \
		helm package `ls -d */`;                       \
		helm repo index . --url=file:///$$repo; \
		cd ../..;                                      \
	done
