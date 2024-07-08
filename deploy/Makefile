.PHONY: config start stop clean
# define two variables
total=7
hackers=2

config:
	./genconfig -val ${total} -target ./config && \
	./gen-compose.sh ${total} ${hackers}

start:
	@docker compose up -d

stop:
	@docker compose down

clean:
	@sudo rm -rf data/*
