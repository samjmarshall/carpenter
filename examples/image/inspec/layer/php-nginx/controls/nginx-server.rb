# encoding: utf-8

title 'nginx-server'

php_version = '7.3'

control 'nginx-server-00' do
  title 'NGINX server config'
  server = nginx_conf.servers[0].params

  # Ensure only the default monitoring server and FPM FastCGI server exists
  describe nginx_conf.servers.length do
    it { should eq 2 }
  end

  describe 'location count' do
    subject { server['location'].length }
    it { should eq 3 }
  end

  server_config = {
    'listen'      => [['*:80']],
    'server_name' => [['_']],
    'root'        => [['/var/www/html/public']],
    'add_header'  => [['Strict-Transport-Security', 'max-age=31536000; includeSubDomains;']],
    'index'       => [['index.html', 'index.htm', 'index.php']],
    'access_log'  => [['/var/log/nginx/00-php-fpm.access.log', 'combined']],
    'error_log'   => [['/var/log/nginx/00-php-fpm.error.log']],
    'if'          => [{
      '_'       => ['($http_x_forwarded_proto', '=', 'http', ')'],
      'rewrite' => [['^', 'https://$host$request_uri?', 'permanent']],
    }],
  }

  server_config.each do |param, value|
    describe param do
      subject { server[param] } 
      it { should eq value }
    end
  end

  base_location = {
    '_'         => ['/'],
    'try_files' => [['$uri', '$uri/', '/index.php?$query_string']],
  }

  base_location.each do |param, value|
    describe param do
      subject { server['location'][0][param] }
      it { should eq value }
    end
  end

  php_location = {
    '_'                       => ['~', '\\.php$'],
    'try_files'               => [['$uri', '=404']],
    'fastcgi_pass'            => [['127.0.0.1:9000']],
    'fastcgi_index'           => [['index.php']],
    'fastcgi_split_path_info' => [['^(.+\\.php)(/.+)$']],
    # 'if'                      => [{
    #   '_'      => ['(!-f', '$document_root$fastcgi_script_name)'],
    #   'return' => [['404']]
    # }],
    'fastcgi_param'           => [
      ['SCRIPT_FILENAME', '$document_root$fastcgi_script_name'],
      ['QUERY_STRING', '$query_string'],
      ['REQUEST_METHOD', '$request_method'],
      ['CONTENT_TYPE', '$content_type'],
      ['CONTENT_LENGTH', '$content_length'],
      ['SCRIPT_NAME', '$fastcgi_script_name'],
      ['REQUEST_URI', '$request_uri'],
      ['DOCUMENT_URI', '$document_uri'],
      ['DOCUMENT_ROOT', '$document_root'],
      ['SERVER_PROTOCOL', '$server_protocol'],
      ['REQUEST_SCHEME', '$scheme'],
      ['HTTPS', '$https', 'if_not_empty'],
      ['GATEWAY_INTERFACE', 'CGI/1.1'],
      ['SERVER_SOFTWARE', 'nginx/$nginx_version'],
      ['REMOTE_ADDR', '$remote_addr'],
      ['REMOTE_PORT', '$remote_port'],
      ['SERVER_ADDR', '$server_addr'],
      ['SERVER_PORT', '$server_port'],
      ['SERVER_NAME', '$server_name'],
      ['REDIRECT_STATUS', '200']
    ]
  }

  php_location.each do |param, value|
    describe param do
      subject { server['location'][1][param] }
      it { should eq value }
    end
  end

  ht_location = {
    '_'    => ['~', '/\\.ht'],
    'deny' => [['all']]
  }

  ht_location.each do |param, value|
    describe param do
      subject { server['location'][2][param] }
      it { should eq value }
    end
  end
end