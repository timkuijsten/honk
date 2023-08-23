
all: honk

honk: .preflightcheck schema.sql *.go go.mod
	env CGO_ENABLED=1 go build -mod=`ls -d vendor 2> /dev/null` -o honk

.preflightcheck: preflight.sh
	@sh ./preflight.sh

clean:
	rm -f honk

test:
	go test
