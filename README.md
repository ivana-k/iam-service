# iam-service

Da bi se servis uspesno pokrenuo sa oort-om, potrebno je oort i magnetar servis postaviti u isti folder gde je i iam-service.
iam-service i oort komuniciraju preko network1 external mreze.

Komande za kreiranje cassandra keyspace (ukoliko ne postoji):
- docker exec -it cassandra cqlsh
- CREATE KEYSPACE IF NOT EXISTS apollo WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };

Port na kom je iam-service: 8002
