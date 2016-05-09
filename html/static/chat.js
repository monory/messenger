function SendMessage(socket) {
    var data = $(".message-text").val();

    setTimeout(function () {
        $(".message-text").val("")
    }, 0);

    if (data.length == 0) {
        return;
    }

    var message = {"body": data};
    socket.send(JSON.stringify(message));
}

function ReceiveMessage(event) {
    var message = $.parseJSON(event.data);

    var messageString = $("<li class=\"message\"></li>");
    messageString.append($("<strong></strong>").text(message.author + ": "));
    messageString.append($("<span></span>").text(message.body));
    $("#chat-messages").append(messageString);

    $(".chat-messages").scrollTop($(".chat-messages")[0].scrollHeight);
}

function CloseConnection() {
    $(".message-text").prop("disabled", true);
    $(".message-text").val("Connection closed. Refresh to chat again!");

    $(".message-button").prop("disabled", true);
}

function getCookie(name) {
    var re = new RegExp(name + "=([^;]+)");
    var value = re.exec(document.cookie);
    return (value != null) ? unescape(value[1]) : 0;
}

function SendToken(socket) {
    socket.send(getCookie("chat_token"));
}

$(function() {
    var socket = new WebSocket("wss://chat.monory.org/ws");

    socket.onopen = function() {
        SendToken(socket);
    };
    socket.onmessage = ReceiveMessage;
    socket.onclose = CloseConnection;

    $('.message-button').click(function() {
        SendMessage(socket);
    });

    $('.message-text').keydown(function(event) {
        if (event.keyCode == 13 && !event.shiftKey) {
            SendMessage(socket);
        }
    });
});
