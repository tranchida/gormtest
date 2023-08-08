docker run -d \
--name mysql-8 \
-p 3306:3306 \
-v ~/mysql_data:/var/lib/mysql \
-e MYSQL_ROOT_PASSWORD='MyStr0ngP@ssw0rd' \
-e MYSQL_USER=go1 \
-e MYSQL_PASSWORD='go1' \
-e MYSQL_DATABASE=sample \
mysql:8