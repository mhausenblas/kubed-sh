#!/usr/bin/env python

from http.server import BaseHTTPRequestHandler,HTTPServer

class AppServerHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.end_headers()
        self.wfile.write(bytes("Hello from Python 3","utf-8"))
        return

try:
    PORT = 8080
    server = HTTPServer(('', PORT), AppServerHandler)
    print ("Serving on port " , PORT)
    server.serve_forever()

except KeyboardInterrupt:
    server.socket.close()
