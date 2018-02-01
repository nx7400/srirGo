REPO = github.com/nx7400/srirGo/

DOC_IMAGES = doc/images

.PHONY: doc get

doc: $(GOPATH)/$(REPO)/server $(GOPATH)/$(REPO)/client $(DOC_IMAGES)/server.png $(DOC_IMAGES)/client.png

$(GOPATH)/$(REPO)/server:
	go get -d github.com/nx7400/srirGo/server 

$(GOPATH)/$(REPO)/client:
	go get -d github.com/nx7400/srirGo/client 

$(DOC_IMAGES)/server.png:
	goviz -i github.com/nx7400/srirGo/server -l | dot -Tpng -o $(DOC_IMAGES)/server.png

$(DOC_IMAGES)/client.png:
	goviz -i github.com/nx7400/srirGo/client -l | dot -Tpng -o $(DOC_IMAGES)/client.png
