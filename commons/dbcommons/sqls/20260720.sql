-- Allow more than one package configuration per app and push channel.
-- The statements are intentionally idempotent because deployments can retry a
-- migration after a process or network interruption.

SET @create_android_package_index = IF(
  (SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'androidpushconfs' AND index_name = 'uniq_channel_package') = 0,
  'CREATE UNIQUE INDEX uniq_channel_package ON androidpushconfs (app_key, push_channel, package)',
  'SELECT 1'
);
PREPARE push_multi_stmt FROM @create_android_package_index;
EXECUTE push_multi_stmt;
DEALLOCATE PREPARE push_multi_stmt;

SET @drop_android_channel_index = IF(
  (SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'androidpushconfs' AND index_name = 'uniq_channel') > 0,
  'DROP INDEX uniq_channel ON androidpushconfs',
  'SELECT 1'
);
PREPARE push_multi_stmt FROM @drop_android_channel_index;
EXECUTE push_multi_stmt;
DEALLOCATE PREPARE push_multi_stmt;

SET @create_ios_package_index = IF(
  (SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'ioscertificates' AND index_name = 'uniq_app_package') = 0,
  'CREATE UNIQUE INDEX uniq_app_package ON ioscertificates (app_key, package)',
  'SELECT 1'
);
PREPARE push_multi_stmt FROM @create_ios_package_index;
EXECUTE push_multi_stmt;
DEALLOCATE PREPARE push_multi_stmt;

SET @drop_ios_app_index = IF(
  (SELECT COUNT(*) FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = 'ioscertificates' AND index_name = 'uniq_package') > 0,
  'DROP INDEX uniq_package ON ioscertificates',
  'SELECT 1'
);
PREPARE push_multi_stmt FROM @drop_ios_app_index;
EXECUTE push_multi_stmt;
DEALLOCATE PREPARE push_multi_stmt;
