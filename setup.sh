
./makeCertificates.sh

docker-compose up -d mysql rabbitmq

sleep 10

./load_ip_address.sh

docker-compose up -d
