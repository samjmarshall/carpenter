# encoding: utf-8

title 'php-fpm'

php_version = '7.3'

control 'php-fpm-00' do
  impact 1.0
  title 'PHP-FPM is installed and running'

  describe package("php#{php_version}-fpm") do
    it { should be_installed }
  end

  describe service("php#{php_version}-fpm") do
    it { should be_installed }
    it { should be_running }
    it { should_not be_enabled }
  end
end

control 'php-fpm-01' do
  impact 0.7
  title 'PHP-FPM emergency restart conditions for failed child processes'

  describe ini("/etc/php/#{php_version}/fpm/php-fpm.conf") do
    its('global.emergency_restart_threshold') { should eq '10' }
    its('global.emergency_restart_interval') { should eq '1m' }
    its('global.process_control_timeout') { should eq '10s' }
  end
end

control 'php-fpm-02' do
  impact 1.0
  title 'PHP short open tags enabled'

  describe ini("/etc/php/#{php_version}/fpm/php.ini") do
    its('PHP.short_open_tag') { should eq 'On' }
  end
end

control 'php-fpm-03' do
  impact 0.7
  title 'PHP inefficient legacy pages'

  describe ini("/etc/php/#{php_version}/fpm/php.ini") do
    its('PHP.max_execution_time') { should eq '30' }
    its('PHP.max_input_time') { should eq '60' }
    its('PHP.memory_limit') { should eq '128M' }
  end
end

control 'php-fpm-04' do
  impact 0.7
  title 'PHP file upload and post body limits'

  describe ini("/etc/php/#{php_version}/fpm/php.ini") do
    its('PHP.post_max_size') { should eq '8M' }
    its('PHP.upload_max_filesize') { should eq '2M' }
    its('PHP.max_file_uploads') { should eq '20' }
  end
end

control 'php-fpm-05' do
  impact 0.7
  title 'PHP session file upload progress config'

  describe ini("/etc/php/#{php_version}/fpm/php.ini") do
    its(['Session', 'session.upload_progress.enabled']) { should eq 'On' }
    its(['Session', 'session.upload_progress.cleanup']) { should eq 'On' }
    its(['Session', 'session.upload_progress.prefix']) { should eq 'upload_progress_' }
    its(['Session', 'session.upload_progress.name']) { should eq 'PHP_SESSION_UPLOAD_PROGRESS' }
    its(['Session', 'session.upload_progress.freq']) { should eq '1%' }
    its(['Session', 'session.upload_progress.min_freq']) { should eq '1' }
  end
end