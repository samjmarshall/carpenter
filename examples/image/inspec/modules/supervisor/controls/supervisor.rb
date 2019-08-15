# encoding: utf-8

title 'supervisor'

control 'supervisor-00' do
  impact 1.0
  title 'Supervisord application binaries'

  describe file('/usr/local/bin/supervisord') do
    it { should be_file }
    its('mode') { should cmp '0755' }
    its('owner') { should eq 'root' }
    its('group') { should eq 'root' }
  end

  describe file('/usr/local/bin/supervisorctl') do
    it { should be_file }
    its('mode') { should cmp '0755' }
    its('owner') { should eq 'root' }
    its('group') { should eq 'root' }
  end

  describe command('/usr/local/bin/supervisord --version') do
    its('exit_status') { should eq 0 }
    its('stdout') { should match(%r{4.0.\d+}) }
  end
end

control 'supervisor-01' do
  impact 1.0
  title 'Supervisord service'

  describe service('supervisord') do
    it { should be_installed }
    it { should be_enabled }
    it { should be_running }
  end
end

control 'supervisor-02' do
  impact 1.0
  title 'Supervisord include configs'

  describe ini('/etc/supervisord.conf') do
    its('include.files') { should eq '/etc/supervisor.d/*.conf' }
  end
end