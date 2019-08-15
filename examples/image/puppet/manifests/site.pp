File { backup => false }

node default {

  class { lookup('modules'): }

}
