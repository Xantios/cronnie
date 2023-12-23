.PHONY: shell import

shell:
	docker exec -it beastmode-app-1 sh

import:
	docker exec -it beastmode-app-1 go run *.go import