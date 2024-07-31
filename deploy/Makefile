.PHONY: config reset start stop clean
# define two variables
total=22
hackers=7
tps=0

config:
	@./genconfig -val ${total} -target ./config && \
	./gen-compose.sh ${total} ${hackers} ${tps}

reset:
	@rm -rf config/keys config/genesis.json config/accounts.json config/master.keys docker-compose.yml


start:
	@docker compose up -d

stop:
	@docker compose down

clean:
	@sudo rm -rf data/*
