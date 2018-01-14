#!/usr/bin/env ruby
require 'socket'
hostname = Socket.gethostname

puts "Hello from Ruby, running on host #{hostname}"
