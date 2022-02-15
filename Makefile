plugin:
	GO111MODULE="on" go build cmd/oc-clusteroperator.go

install:
	sudo mv ./oc-clusteroperator /usr/local/bin/