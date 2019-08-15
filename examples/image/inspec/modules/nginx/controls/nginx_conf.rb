# encoding: utf-8

title 'nginx-conf'

# include_controls 'nginx-baseline'

nginx_params = nginx_conf.params

control 'nginx-conf-00' do
  impact 0.7
  title 'NGINX basic server checks'

  describe command('nginx -t') do
    its('exit_status') { should eq 0 }
    its('stderr') { should match(%r{syntax is ok}) }
    its('stderr') { should match(%r{test is successful}) }
  end
end

control 'nginx-conf-01' do
  title 'NGINX service config'

  {
    'user'                 => ['www-data', 'www-data'],
    'worker_processes'     => ['auto'],
    'worker_rlimit_nofile' => ['1024'],
    'pid'                  => ['/var/run/nginx.pid'],
  }.each do |param, value|
    describe param do
      subject { nginx_params[param].flatten }
      it { should eq value }
    end
  end
end

control 'nginx-conf-02' do
  title 'NGINX event config'

  {
    'accept_mutex'       => ['on'],
    'accept_mutex_delay' => ['500ms'],
    'worker_connections' => ['1024'],
  }.each do |param, value|
    describe param do
      subject { nginx_params['events'].first[param].flatten }
      it { should eq value }
    end
  end
end

control 'nginx-conf-03' do
  title 'NGINX http basic config'

  {
    'default_type'                  => ['application/octet-stream'],
    'access_log'                    => ['/var/log/nginx/access.log'],
    'error_log'                     => ['/var/log/nginx/error.log', 'error'],
    'sendfile'                      => ['on'],
    'tcp_nopush'                    => ['on'],
    'server_tokens'                 => ['off'],
    'types_hash_max_size'           => ['1024'],
    'types_hash_bucket_size'        => ['512'],
    'server_names_hash_bucket_size' => ['64'],
    'server_names_hash_max_size'    => ['512'],
    'keepalive_timeout'             => ['65s'],
    'keepalive_requests'            => ['100'],
    'client_body_timeout'           => ['60s'],
    'client_body_temp_path'         => ['/var/nginx/client_body_temp'],
    'client_max_body_size'          => ['10m'],
    'client_body_buffer_size'       => ['128k'],
    'send_timeout'                  => ['60s'],
    'lingering_timeout'             => ['5s'],
    'tcp_nodelay'                   => ['on'],
  }.each do |param, value|
    describe param do
      subject { nginx_params['http'].first[param].flatten }
      it { should eq value }
    end
  end
end

control 'nginx-conf-04' do
  title 'NGINX http gzip config'

  {
    'gzip'              => ['on'],
    'gzip_comp_level'   => ['1'],
    'gzip_disable'      => ['msie6'],
    'gzip_min_length'   => ['20'],
    'gzip_http_version' => ['1.1'],
    'gzip_proxied'      => ['off'],
    'gzip_vary'         => ['off'],
  }.each do |param, value|
    describe param do
      subject { nginx_params['http'].first[param].flatten }
      it { should eq value }
    end
  end
end

control 'nginx-conf-05' do
  title 'NGINX http proxy config'

  {
    'proxy_temp_path'                => ['/var/nginx/proxy_temp'],
    'proxy_connect_timeout'          => ['90s'],
    'proxy_send_timeout'             => ['90s'],
    'proxy_read_timeout'             => ['90s'],
    'proxy_buffers'                  => ['32', '4k'],
    'proxy_buffer_size'              => ['8k'],
    'proxy_set_header'               => ['Host', '$host', 'X-Real-IP', '$remote_addr', 'X-Forwarded-For', '$proxy_add_x_forwarded_for', 'Proxy', {:value=>[]}],
    'proxy_headers_hash_bucket_size' => ['64'],
  }.each do |param, value|
    describe param do
      subject { nginx_params['http'].first[param].flatten }
      it { should eq value }
    end
  end
end
