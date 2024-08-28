up:
	sudo docker-compose up -d --build --remove-orphans

upterm:
	sudo docker-compose up --build --remove-orphans

down:
	sudo docker-compose down

tag-push:
	make tag-push-frontend
	make tag-push-api
	make tag-push-email

tag-push-frontend:
	sudo docker tag taxarific-frontend-service:latest morineth/taxarific-frontend:latest && sudo docker push morineth/taxarific-frontend:latest 

tag-push-api:
	sudo docker tag taxarific-backend-users-api:latest morineth/taxarific-backend-user:latest && sudo docker push morineth/taxarific-backend-user:latest

tag-push-email:
	sudo docker tag taxarific-email-service:latest morineth/taxarific-backend-email:latest && sudo docker push morineth/taxarific-backend-email:latest