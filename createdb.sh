mysql -h 127.0.0.1 -P 43306 -u root -pmysqlvotmdnjem <<MYSQL_SCRIPT
CREATE DATABASE IF NOT EXISTS payhere DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci;
MYSQL_SCRIPT

mysql -h 127.0.0.1 -P 43306 -u root -pmysqlvotmdnjem <<MYSQL_SCRIPT
CREATE USER 'payhere'@'%' IDENTIFIED BY 'payherepass';
GRANT ALL PRIVILEGES ON payhere.* TO 'payhere'@'%';
FLUSH PRIVILEGES;
MYSQL_SCRIPT