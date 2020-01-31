# encoding: utf-8

title 'base'

include_controls 'chrony-aws'

include_controls 'linux-baseline' do
  # Entropy is typically below 1000 '/proc/sys/kernel/random/entropy_avail'
  skip_control 'os-08'
end

include_controls 'linux-patch-baseline'