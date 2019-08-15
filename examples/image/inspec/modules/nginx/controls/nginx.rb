# encoding: utf-8

title 'nginx'

control 'nginx-00' do
  impact 1.0
  title 'NGINX is installed and running'

  describe package('nginx') do
    it { should be_installed }
    its('version') { should match(%r{1\.}) }
  end

  describe service('nginx') do
    it { should be_installed }
    it { should be_running }
    it { should_not be_enabled }
  end
end