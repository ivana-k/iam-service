# iam-service

Da bi se servis uspesno pokrenuo sa oort-om, potrebno je oort i magnetar servis postaviti u isti folder gde je i iam-service.
iam-service i oort komuniciraju preko network1 external mreze.

Pokretanje servisa:
- docker network create network1
- docker compose up za oort
- docker compose up za iam-service

  Port na kom je iam-service: 8002
