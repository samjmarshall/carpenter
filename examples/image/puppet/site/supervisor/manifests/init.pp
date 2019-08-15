class supervisor {

  package { 'python3-pip': }

  class { 'supervisord': package_provider => 'pip3' }

  Package['python3-pip']
    -> Class['supervisord']

}
