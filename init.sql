-- Initialize databases for both Ruby and Go APIs
CREATE DATABASE IF NOT EXISTS campsite_development;
CREATE DATABASE IF NOT EXISTS campsite_test;
CREATE DATABASE IF NOT EXISTS campsite_go;

-- Grant permissions to campsite user
GRANT ALL PRIVILEGES ON campsite_development.* TO 'campsite'@'%';
GRANT ALL PRIVILEGES ON campsite_test.* TO 'campsite'@'%';
GRANT ALL PRIVILEGES ON campsite_go.* TO 'campsite'@'%';

FLUSH PRIVILEGES;