// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Modified by Tuomas Jaakola <tuomas.jaakola@iki.fi>

// +build ignore

package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Idle Timeout Tester</title>
    <!-- Bootstrap -->
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">
    <!-- HTML5 shim and Respond.js for IE8 support of HTML5 elements and media queries -->
    <!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
    <!--[if lt IE 9]>
      <script src="https://oss.maxcdn.com/html5shiv/3.7.3/html5shiv.min.js"></script>
      <script src="https://oss.maxcdn.com/respond/1.4.2/respond.min.js"></script>
    <![endif]-->
  </head>
  <body>
    <div class="container">
      <div class="row">
        <div class="col-sm-12">
          <h1>Idle Timeout Tester</h1>
        </div>
      </div> <!-- .row -->

      <div class="row">
        <div class="col-sm-12">
	  <p>Test if idle TCP connections are timed out in your internet connection. Choose time to and click Start. TCP connection using WebSocket is created, given time is idled and after that a test message is sent to server. The result will be shown after the time has elapsed.</p>
	</div>
      </div>

      <div class="row" style="margin-bottom: 25px;">
        <div class="col-sm-12">
            <select class="form-control" id="duration" style="width: 100px; display: inline; margin-right: 25px">
              <option value="70">1 min</option>
              <option value="310">5 min</option>
              <option value="610">10 min</option>
              <option value="910">15 min</option>
              <option value="1810">30 min</option>
              <option value="3610">60 min</option>
            </select>

          <button type="button" class="btn btn-lg btn-primary" id="start">
            Start
          </button>
        </div>
      </div> <!-- .row -->

      <div class="row" id="progress" style="display: none;">
        <div class="col-sm-12">
          Idling, please wait...
          <div class="progress">
            <div class="progress-bar progress-bar-striped active" role="progressbar" style="width:0%">
              <span id="time-left"></span>
            </div>
          </div>
        </div>
      </div>

      <div class="row" id="success" style="display: none;">
        <div class="col-sm-12">
          <div class="alert alert-success" role="alert">
            <strong>Good news!</strong> Idling connection was not timed out.
          </div>
        </div>
      </div>

      <div class="row" id="failure" style="display: none;">
        <div class="col-sm-12">
          <div class="alert alert-danger" role="alert">
            <strong>Sorry!</strong> Idling connection was timed out. Reload and try smaller duration.
          </div>
        </div>
      </div>

    </div> <!-- .container -->

    <!-- jQuery (necessary for Bootstrap's JavaScript plugins) -->
    <script src="https://ajax.googleapis.com/ajax/libs/jquery/1.12.4/jquery.min.js"></script>
    <!-- Include all compiled plugins (below), or include individual files as needed -->
    <script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/js/bootstrap.min.js" integrity="sha384-Tc5IQib027qvyjSMfHjOMaLkfuWVxZxUPnCJA7l2mCWNIpG9mGCD8wGNIcPD7Txa" crossorigin="anonymous"></script>
    <script>

var started;
var duration;
var ws;

function updateProgressBar(txt, progress) {
    $('#time-left').html(txt);
    $('.progress-bar').css({ width: '' + progress + '%' });
}

function connect() {
    updateProgressBar('Connecting...', 100);
    if (ws) {
        return false;
    }
    ws = new WebSocket("{{.}}");
    ws.onopen = function(evt) {
	console.debug('Connected');
	startIdling();
    }
    ws.onclose = function(evt) {
        ws = null;
	console.debug('Websocket closed');
    }
    ws.onmessage = function(evt) {
        console.debug('Message from server: ', evt.data);
	ws.close();
	onSuccess();
    }
    ws.onerror = function(evt) {
        onFailure(evt);
    }
    return false;
}

function updateProgress() {
    var now = new Date();
    var secs = (now.getTime() - started.getTime()) / 1000;
    var left = Math.round(duration - secs);
    if (left > 0) {
        var minsLeft = Math.floor(left / 60);
        var secsLeft = left % 60;

        var txt = minsLeft < 10 ? '0' + minsLeft : minsLeft;
        txt += ':';
        txt += secsLeft < 10 ? '0' + secsLeft : secsLeft;

	updateProgressBar(txt, secs / duration * 100);
    } else {
        end();
    }
}

function startIdling() {
    t = setInterval(function() {
        updateProgress();
    }, 1000);
}

function start() {
    $('#success').hide();
    $('#failure').hide();
    duration = parseInt($('#duration').val());
    started = new Date();
    $('#start').attr('disabled', 'disabled');
    $('#progress').show();
    connect();
}

function end() {
    $('#progress').hide();
    clearInterval(t);
    t = null;
    // test if connection still works by sending a message to server
    ws.send('a');
}

function onSuccess() {
    $('#success').show();
    $('#start').removeAttr('disabled');
}

function onFailure(evt) {
    $('#failure').show();
    console.error('Websocket error:', evt);
    $('#start').removeAttr('disabled');
}

function test(duration) {
    console.log('Testing ', duration, ' sec timeout');
    start();
}

$(document).ready(function() {
    $('#start').on('click', function() {
	start();
    });
});

    </script>
  </body>
</html>
`))
