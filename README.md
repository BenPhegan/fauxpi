fauxpi fake API proxy server
============================

[![Build Status](https://travis-ci.org/BenPhegan/fauxpi.svg?branch=master)](https://travis-ci.org/BenPhegan/fauxpi)

When working with APIs, more often than not during development you want to work
with a fixed version of the responses provided. This is especially true if the
API is still under development, and maybe even still needs input on how to
output something. This is what fauxpi is for!

fauxpi is aimed at reimplementing some of the functionality from [canned](https://github.com/sideshowcoder/canned) in Go with a few bells and whistles.  Much of the documentation here is copied directly from Canned, and will initially be used as feature requirements and tracking.

Basically provides extensions on top of [goproxy](https://github.com/elazarl/goproxy)

Currently very rough and feature incomplete.

What does it do?
----------------
fauxpi maps a folder structure to API responses

    /comment/any.get.json
    /comment/index.get.html

requests like

    GET /comment/:id

are served as

    Content-Type: application/json
    { "content": "I am a comment", "author": "sideshowcoder" }

requests like

    GET /content/

are served as

    Content-Type: text/html
    <html>
      <body>Some html in here</body>
    </html>

Awesome! so what is supported?
------------------------------
Currently Canned supports the basic REST-API mapping, as well as custom method
mapping with nested endpoints.

    file                            | request					|	implemented
    ---------------------------------------------------------------------------
	/index.get.json                 | GET /						|	yes
    /any.get.json                   | GET /:id					|	no
    /_search.get.json               | GET /search				|	yes
    /comments/index.get.json        | GET /comments/			|	yes
    /comments/any.get.json          | GET /comments/:id			|	no
    /comments/_search.get.json      | GET /comments/search		|	yes
	
You can also ask it to repond to various hosts and protocols via command line options, providing per-host based configuration.

    file                            					| request									|			implemented
    ------------------------------------------------------------------------------------------------------------------
	/http/www.test.io/index.get.json                 	| GET http://www.test.io/					|	yes
    /http/www.test.io/any.get.json                   	| GET http://www.test.io/:id				|	no
    /http/www.blah.io/_search.get.json               	| GET http://www.blah.io/search			|	yes
    /https/www.test.io/comments/index.get.json       	| GET https://www.test.io/comments/		|	yes
    /http/www.test.io/comments/any.get.json          	| GET http://www.test.io/comments/:id		|	no
    /http/www.test.io/comments/_search.get.json      	| GET http://www.test.io/comments/search	|	yes

You can even add query parameters to your filenames to return different responses on the same route. If the all query params in a filename match the incoming request, this file will be returned. It will fall back to returning the file with no query params if it exists.

    file                            | request											|	implemented
    -----------------------------------------------------------------------------------------------
    /index?name=Superman.get.json   | GET /?name=Superman&NotAllParams=NeedToMatch	|		no
    /_search?q=hello.get.json       | GET /comments/search?q=hello					|		no
    /_search.get.json               | GET /comments/search?iam=soignored				|		no

Same support is available for PUT, POST, etc.

    file                        | request											|	implemented
    --------------------------------------------------------------------------------------------
    /index.post.json            | POST serves /... + CORS Headers				|		no
    /index.put.json             | PUT serves /... + CORS Headers					| 		no

If CORS support is enabled additionally options will be available as a http verb
and all requests will serve the CORS Headers as well

    /                           | OPTIONS serve all the options needed for CORS
    /index.get.json             | GET serves /... + CORS Headers

If you need some custum return codes, just add them to the file via adding a
file header like so

    //! statusCode: 201
    <html>
      <body>Created something successfully! Happy!</body>
    </html>

The header will be stripped before sending and the statusCode will be set.

You can also override the default content types by adding a custom content type to the file header:

    //! contentType: "application/vnd.custom+xml"
    <xml>
        <created>1</created>
    </xml>

This will be returned with a `Content-type: application/vnd.custom+xml` header.

Multiple headers need to be written on one single line and comma-separated, like so:

    //! statusCode:201, contentType: "application/vnd.custom+xml"



