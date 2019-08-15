# encoding: utf-8

title 'chrony'

control 'ntp-00' do
  impact 1.0
  title 'NTP not installed'

  describe package('ntp') do
    it { should_not be_installed }
  end
end

control 'chrony-00' do
  impact 1.0
  title 'Chrony settings'

  describe file('/etc/chrony/chrony.conf') do
    its('content') { should match(%r{^server 169\.254\.169\.123 prefer iburst minpoll 4 maxpoll 4$}) }
    its('content') { should match(%r{^rtcsync$}) }
  end
end

control 'chrony-01' do
  impact 1.0
  title 'Chrony sync'

  describe command('chronyc sources') do
    its('stdout') { should match (%r{^\^[\*\?] 169\.254\.169\.123 \s* \d \s* 4 .*$}) }
  end
end

control 'chrony-02' do
  impact 1.0
  title 'Chrony service'

  describe service('chrony') do
    it { should be_running }
    it { should be_enabled }
  end
end