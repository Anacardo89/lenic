### Upload and start app on AWS
Connect to EC2
- ssh -i .ssh/<ssh_key> ec2-user@<ec2_ip>

Inside EC2
1. sudo dnf update && sudo dnf upgrade -y
2. sudo dnf install -y postgresql17
3. sudo mkdir /opt/lenic
4. sudo chown ec2-user /opt/lenic
5. mkdir /opt/lenic/bin /opt/lenic/config

In Project root
1. GOOS=linux GOARCH=amd64 go build -o lenic ./cmd/main
2. GOOS=linux GOARCH=amd64 go build -o migrate ./cmd/migrate
3. scp -i ~/.ssh/<ssh_key> ./lenic ./migrate ec2-user@<ec2_ip>:/opt/lenic/bin
4. scp -i ~/.ssh/<ssh_key> -r ./db/users ./db/migrations ./frontend/templates ./frontend/static prod.env ec2-user@<ec2_ip>:/opt/lenic
5. scp -i ~/.ssh/<ssh_key> ./config/config.yaml ec2-user@<ec2_ip>:/opt/lenic/config
6. scp -i ~/.ssh/<ssh_key> ./lenic.service ./lenic-migrator.service ec2-user@<ec2_ip>:/tmp
7. rm lenic migrate

Back In EC2
1. sudo mv /tmp/lenic.service /etc/systemd/system/lenic.service
2. sudo chown root:root /etc/systemd/system/lenic.service
3. sudo chmod 644 /etc/systemd/system/lenic.service
4. sudo mv /tmp/lenic-migrator.service /etc/systemd/system/lenic-migrator.service
5. sudo chown root:root /etc/systemd/system/lenic-migrator.service
6. sudo chmod 644 /etc/systemd/system/lenic-migrator.service
7. psql -h <rds_endpoint> -p 5432 -U lenic_admin -d lenicDB -f /opt/lenic/users/migrator.sql -f /opt/lenic/users/runner.sql
8. rm -rf /opt/lenic/users
9. sudo setcap 'cap_net_bind_service=+ep' /opt/lenic/bin/lenic
9. sudo systemctl daemon-reload
10. sudo systemctl enable lenic-migrator
11. sudo systemctl enable lenic
12. sudo systemctl start lenic-migrator
13. sudo systemctl restart lenic
14. sudo systemctl status lenic

### Check logs
sudo journalctl -u lenic -f
