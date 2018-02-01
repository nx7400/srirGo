REPO = github.com/nx7400/srirGo/

.PHONY: doc get

doc: $(GOPATH)/$(REPO)/server $(GOPATH)/$(REPO)/client doc/server.png doc/client.png

$(GOPATH)/$(REPO)/server:
	go get -d github.com/nx7400/srirGo/server 

$(GOPATH)/$(REPO)/client:
	go get -d github.com/nx7400/srirGo/client 

doc/server.png:
	goviz -i github.com/nx7400/srirGo/server -l | dot -Tpng -o doc/server.png

doc/client.png:
	goviz -i github.com/nx7400/srirGo/client -l | dot -Tpng -o doc/client.png
