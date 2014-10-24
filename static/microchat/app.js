$(function () {
  var s = document.location.search.substring(1).split('&');
  //document.write(document.location)

  var GET = {};

  for (var i = 0; i < s.length; i++) {
    var parts = s[i].split('=');
    GET[unescape(parts[0])] = unescape(parts[1]);
  }

  var userUrl = '/api/user?sessionId=' + GET.sessionId
  var usersUrl = '/api/users?sessionId=' + GET.sessionId
  //document.write('<br><a href="' + userUrl + '">user</a>')
  //document.write('<br><a href="/api/users">users</a>')
  //document.write('<br><a href="/api/messages?sessionId="' + GET.sessionId + '">messages</a>')

  $.ajax(userUrl, {
    dataType: 'json',
    success: function (data) {
      $('#username').text(data.username);
    }
  });

  var getUsers = function () {
    $.ajax(usersUrl, {
      dataType: 'json',
      success: function (data) {
        $('#users').text(data);
      }
    });
  }

  getUsers();

  if (!window.location.origin) { // Some browsers (mainly IE) do not have this property, so we need to build it manually...
    window.location.origin = window.location.protocol + '//' + window.location.hostname + (window.location.port ? (':' + window.location.port) : '');
  }


  var sock = new SockJS(window.location.origin + '/ws')

  sock.onopen = function () {
    // console.log('connection open');
    sock.send(GET.sessionId);
    document.getElementById("status").innerHTML = "connected";
    document.getElementById("send").disabled = false;
  };

  sock.onmessage = function (e) {
    if (e.data.indexOf("[info]") === 0) {
      getUsers();
    }
    document.getElementById("output").value += e.data + "\n";
  };

  sock.onclose = function () {
    // console.log('connection closed');
    document.getElementById("status").innerHTML = "disconnected";
    document.getElementById("send").disabled = true;
  };

  window.sendMessage = function () {
    $input = $('#input');
    message = $input.val();
    sock.send(message);
    $input.val('');
  };

});
