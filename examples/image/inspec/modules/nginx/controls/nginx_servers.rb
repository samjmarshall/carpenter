# encoding: utf-8

title 'nginx-monitoring'

control 'nginx-monitoring-00' do
  impact 0.7
  title 'NGINX monitoring server config'
  desc 'NGINX server status endpoint is available and confgured correctly'

  monitoring_server = {
    'listen'      => [['*:80']],
    'server_name' => [['localhost']],
    'access_log'  => [['off']],
    'location'    => [{
      '_'           => ['/nginx_status'],
      'allow'       => [['127.0.0.1']],
      'deny'        => [['all']],
      'stub_status' => [['on']]
    }]
  }

  describe nginx_conf.servers[1] do
    its('params') { should eq monitoring_server }
  end
end