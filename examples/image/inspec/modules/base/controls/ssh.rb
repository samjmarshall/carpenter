# encoding: utf-8

title 'ssh'

include_controls 'ssh-baseline' do
  # Allow sshd to listen on the default address `0.0.0.0`
  skip_control 'sshd-09'

  # DH primes in Ubuntu 16.04 do not pass
  skip_control 'sshd-48'
end

control 'sshd-config-00' do
  impact 1.0
  title 'Public key authentication via Bastion'

  describe file('/etc/ssh/sshd_config') do
    it { should be_file }
    its('content') { should match(%r{^UsePAM yes$}) }
  end
end

control 'ssh-authorized-keys-00' do
  impact 1.0
  title 'Administrators authorized SSH keys'
  desc 'Ensure only expected public keys are included in admin user SSH authorized_keys files'

  {
    'test': '51b6340ea349a5048c45038b12a73307',
  }.each do |user, md5sum|
    describe file("/home/#{user}/.ssh/authorized_keys") do
      it { should be_file }
      its('mode') { should cmp '0600' }
      its('md5sum') { should eq md5sum }
      its('content') { should match(%r{^ssh-rsa .+ #{user}@mydomain.com$}) }
    end
  end

  describe file('/home/ubuntu/.ssh/authorized_keys') do
    it { should be_file }
    its('mode') { should cmp '0600' }
    its('md5sum') { should eq '2cb8421047941feb77b5f3b6090bf8ed' }
    its('content') { should match(%r{^ssh-rsa .+ test@mydomain.com$}) }
  end
end