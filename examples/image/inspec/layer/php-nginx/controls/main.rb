# encoding: utf-8

title 'php-nginx'

include_controls 'php'
include_controls 'nginx'

control 'users-00' do
  impact 1.0
  title 'Ensure www-data home directory is configured'

  describe file('/var/www') do
    it { should be_directory }
    it { should be_setgid }
    its('owner') { should eq 'www-data' }
    its('group') { should eq 'www-data' }
  end
end