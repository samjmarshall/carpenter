# encoding: utf-8

title 'users'

control 'users-00' do
  impact 1.0
  title 'Administrators'
  desc 'Ensure administrators are configured'

  [
    'test'
  ].each do |user|
    describe user(user) do
      its('groups') { should include 'admin' }
    end

    describe file("/etc/sudoers.d/90-#{user}") do
      it { should exist }
      its('owner') { should eq 'root' }
      its('group') { should eq 'root' }
      its('mode') { should cmp '0440' }
      its('content') { should eq "#{user} ALL=(ALL) NOPASSWD:ALL" }
    end
  end
end