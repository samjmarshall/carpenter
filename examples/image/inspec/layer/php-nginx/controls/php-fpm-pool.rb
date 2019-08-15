# encoding: utf-8

title 'php-fpm'

php_version = '7.3'

control 'php-fpm-pool-00' do
  impact 1.0
  title 'PHP-FPM pool listen config'

  describe ini("/etc/php/#{php_version}/fpm/pool.d/www.conf") do
    its('www.listen') { should eq "/run/php/php#{php_version}-fpm.sock" }
    its(['www', 'listen.backlog']) { should eq '-1' }
    its(['www', 'listen.owner']) { should eq 'www-data' }
    its(['www', 'listen.group']) { should eq 'www-data' }
    its(['www', 'listen.mode']) { should eq '0660' }
  end

  describe file('/run/php') do
    it { should be_directory }
    its('mode') { should cmp '0755' }
    its('owner') { should eq 'www-data' }
    its('group') { should eq 'www-data' }
  end

  describe file("/run/php/php#{php_version}-fpm.sock") do
    it { should be_socket }
    its('mode') { should cmp '0660' }
    its('owner') { should eq 'www-data' }
    its('group') { should eq 'www-data' }
  end
end

control 'php-fpm-pool-01' do
  impact 1.0
  title 'PHP-FPM pool process manager settings'

  describe ini("/etc/php/#{php_version}/fpm/pool.d/www.conf") do
    its('www.pm') { should eq 'ondemand' }
    its(['www', 'pm.max_children']) { should eq '150' }
    its(['www', 'pm.process_idle_timeout']) { should eq '10s' }
    its(['www', 'pm.max_requests']) { should eq '0' }
  end
end