#!/usr/bin/env python

from BaseHTTPServer import BaseHTTPRequestHandler, HTTPServer
import SocketServer

class S(BaseHTTPRequestHandler):
    def _set_headers(self):
        self.send_response(200)
        self.send_header('Content-type', 'text/html')
        self.end_headers()

    def do_GET(self):
        print("handling GET request")
        self._set_headers()
        self.wfile.write("<html><body><p>Hello from Python!</p></body></html>")

def run(server_class=HTTPServer, handler_class=S):
    server_address = ('', 80)
    httpd = server_class(server_address, handler_class)
    print('Starting simple Python server...')
    httpd.serve_forever()

if __name__ == "__main__":
    from sys import argv
    run()
