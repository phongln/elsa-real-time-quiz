DC = docker compose

start:
	$(DC) build app
	$(DC) up -d

stop:
	$(DC) down

rebuild-app:
	$(DC) build app

restart-app: rebuild-app
	$(DC) restart app

db-setup:
	docker cp schema.sql mysql://schema.sql
	$(DC) exec mysql bash -c "mysql -u root -p quizdb < schema.sql"

.PHONY: start stop rebuild-app restart-app