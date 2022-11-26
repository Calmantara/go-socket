# building the image
docker build -t go-sock .

# running local docker container
docker run \
	--rm \
	--network none \
	--memory 1g \
	-v $(pwd)/sock:/var/run/dev-test \
	go-sock /var/run/dev-test/sock