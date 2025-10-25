TARGETS = chatio-cat

chatio-cat: chatio-cat.go
	GOARCH=arm64 GOOS=linux go build -o bootstrap

chatio-cat.zip: chatio-cat
	zip chatio-cat.zip bootstrap

deploy: chatio-cat.zip
	aws --profile deploy-chatio-cat --no-cli-pager lambda update-function-code --function-name chatio-cat --zip-file fileb://chatio-cat.zip

clean:
	rm -f ${TARGETS} ${TARGETS}.zip