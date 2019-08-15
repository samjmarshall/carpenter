class base(
  Hash $users              = {},
  Hash $ssh_authorized_key = {},
  Hash $class_dependencies = {},
) {

  include base::users
  include base::syslog

  $class_dependencies.each |String $class, Array $dependencies| {
    $dependencies.each |String $dependency| {
      Class[$class]
        -> Class[$dependency]
    }
  }

}
