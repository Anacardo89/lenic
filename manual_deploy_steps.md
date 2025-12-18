Connect to EC2
`ssh -i .ssh/<ssh_key> ec2-user@<ec2_ip>`

Inside EC2
1. `sudo dnf update && sudo dnf upgrade -y`
2. `sudo dnf install -y postgresql17`
3. `sudo mkdir /opt/lenic`
4. `sudo chown ec2-user /opt/lenic`
5. `mkdir /opt/lenic/bin /opt/lenic/config`

In Project root
1. `GOOS=linux GOARCH=amd64 go build -o lenic ./cmd/main`
2. `GOOS=linux GOARCH=amd64 go build -o migrate ./cmd/migrate`
3. `scp -i ~/.ssh/<ssh_key> ./lenic ./migrate ec2-user@<ec2_ip>:/opt/lenic/bin`
4. `scp -i ~/.ssh/<ssh_key> -r ./db/users ./db/migrations ./frontend/templates ./frontend/static prod.env ec2-user@<ec2_ip>:/opt/lenic`
5. `scp -i ~/.ssh/<ssh_key> ./config/config.yaml ec2-user@<ec2_ip>:/opt/lenic/config`
6. `rm lenic migrate`

Back In EC2
1. `source /opt/lenic/prod.env`
2. `rm /opt/lenic/prod.env`
2. `psql -h <rds_endpoint> -p 5432 -U lenic_admin -d lenicDB -f /opt/lenic/users/migrator.sql`
4. `/opt/lenic/bin/migrate`
2. `psql -h <rds_endpoint> -p 5432 -U lenic_admin -d lenicDB -f /opt/lenic/users/runner.sql`
3. `rm -rf /opt/lenic/users`
5. `sudo -E /opt/lenic/bin/lenic &`