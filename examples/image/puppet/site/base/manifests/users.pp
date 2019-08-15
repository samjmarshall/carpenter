class base::users(
  Boolean $www_data = false,
) {

  # DEPRECATED - ubuntu user SSH from multiple keys
  file { '/home/ubuntu/.ssh/authorized_keys':
    mode    => '0600',
    content => template("${module_name}/authorized_keys.erb"),
  }

  if $www_data {
    file { '/var/www':
      ensure => directory,
      owner  => 'www-data',
      group  => 'www-data',
      mode   => 'g+s',
    }
  }

  $base::users.each |String $user, Hash $attributes| {
    if $attributes['groups'] {
      $groups = $attributes['groups']
    } else {
      $groups = []
    }

    $user_name = $user.regsubst(/@.+/, '', 'G')

    user { $user:
      name       => $user_name,
      home       => "/home/${user_name}",
      managehome => true,
      groups     => $groups,
    }

    if $attributes['key'] {
      file { "/home/${user_name}/.ssh":
        ensure => directory,
        owner  => $user_name,
        group  => $user_name,
        mode   => '0700',
      }

      file { "/home/${user_name}/.ssh/authorized_keys":
        mode    => '0600',
        content => "ssh-rsa ${$attributes['key']} ${user}",
      }

      User[$user]
        -> File["/home/${user_name}/.ssh"]
        -> File["/home/${user_name}/.ssh/authorized_keys"]
    }

    if 'admin' in $groups {
      file { "/etc/sudoers.d/90-${user_name}":
        owner   => 'root',
        group   => 'root',
        mode    => '0440',
        content => "${user_name} ALL=(ALL) NOPASSWD:ALL",
      }
    }
  }

}
