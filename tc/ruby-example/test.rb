require 'bundler/setup'

Bundler.require(:default)

pp Nokogiri::XML("<root>TEST</root>")
