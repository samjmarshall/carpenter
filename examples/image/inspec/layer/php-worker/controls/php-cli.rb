# encoding: utf-8

title 'php-cli'

php_version = '7.3'

control 'php-cli-00' do
  impact 0.7
  title 'PHP worker config'

  describe ini("/etc/php/#{php_version}/cli/php.ini") do
    its('PHP.max_execution_time') { should eq '900' }
    its('PHP.memory_limit') { should eq '384M' }
  end
end