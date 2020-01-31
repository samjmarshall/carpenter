# encoding: utf-8

title 'php'

php_version = '7.3'

control 'php-00' do
  impact 1.0
  title "PHP #{php_version}"

  ['common', 'cli'].each do |package|
    describe package("php#{php_version}-#{package}") do
      it { should be_installed }
    end
  end
end

control 'php-dev-00' do
  impact 0.7
  title 'Ensure unwanted package are not present'

  describe package("php#{php_version}-dev") do
    it { should_not be_installed }
  end
end

control 'php-fpm-00' do
  impact 1.0
  title "PHP #{php_version}"

  only_if('nginx is installed') do
    package('nginx').installed?
  end

  describe package("php#{php_version}-fpm") do
    it { should be_installed }
  end

  describe service("php#{php_version}-fpm") do
    it { should be_installed }
  end
end

control 'php-fpm-01' do
  impact 1.0
  title "PHP #{php_version}"

  only_if('nginx is not installed') do
    !package('nginx').installed?
  end

  describe package("php#{php_version}-fpm") do
    it { should_not be_installed }
  end

  describe service("php#{php_version}-fpm") do
    it { should_not be_installed }
  end
end

control 'php-extensions-00' do
  impact 1.0
  title "PHP #{php_version} extensions"

  [
    'ctype',
    'curl',
    'json',
    'mysqlnd',
    'openssl',
    'PDO',
    'readline',
    'tokenizer',
    'xml',
  ].each do |extension|
    describe command('php -m') do
      its('stdout') { should match(%r{^#{extension}$}) }
    end
  end
end