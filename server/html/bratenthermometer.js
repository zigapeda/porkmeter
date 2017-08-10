var pushToken = null;

var config = {
    apiKey: "AIzaSyD3TZlq2TIS9C719NhvekAL2LfA-Fsr3h4",
    authDomain: "porkmeter-b3cf2.firebaseapp.com",
    databaseURL: "https://porkmeter-b3cf2.firebaseio.com",
    projectId: "porkmeter-b3cf2",
    storageBucket: "porkmeter-b3cf2.appspot.com",
    messagingSenderId: "615966299981"
};

firebase.initializeApp(config);

const messaging = firebase.messaging();

function addZero(i) {
    if (i < 10) {
        i = "0" + i;
    }
    return i;
}

function formatDate(d) {
    var h = d.getHours();
    var m = addZero(d.getMinutes());
    var s = addZero(d.getSeconds());
    return h + ":" + m + ":" + s;
}

function progress(data){
    if(data.error == null) {
        $("#fleischCards").html("");
        $("#smokerCards").html("");
        for(var i in data.success.Temps) {
            var t = data.success.Temps[i];
            if(t.Meter.Type == "fleisch") {
                if(t.Temp >= "95") {
                    $("#fleischCards").append("<div class='col-md-2 col-xs-6'><div class='tempcard tempgreencard'>"+t.Temp+" °C</div></div>");
                } else {
                    $("#fleischCards").append("<div class='col-md-2 col-xs-6'><div class='tempcard'>"+t.Temp+" °C</div></div>");
                }
            } else {
                if(t.Temp <= "100" || t.Temp > "130") {
                    $("#smokerCards").append("<div class='col-md-2 col-xs-6'><div class='tempcard tempredcard'>"+t.Temp+" °C</div></div>");
                } else {
                    $("#smokerCards").append("<div class='col-md-2 col-xs-6'><div class='tempcard'>"+t.Temp+" °C</div></div>");
                }
            }
        }
        var date = new Date(data.success.Time);
        $("#status").html("Letztes Update: " + formatDate(date));
    }
}

function getData() {
    $.get("/api/GetTemps", progress);
}

function handlePush() {
    console.log("handlePush gecalled");
    if($("#pushswitch").is(":checked") == true) {
        enablePush();
    } else {
        disablePush();
    }
}

function changePushSwitch(value) {
    $('#pushswitch').unbind("change");
    $('#pushswitch').bootstrapToggle(value);
    $('#pushswitch').change(handlePush);
}

function checkPushState() {
    if ('serviceWorker' in navigator && 'PushManager' in window) {
        //push is available
        messaging.requestPermission()
        .then(function() {
            messaging.getToken()
            .then(function(currentToken) {
                if (currentToken) {
                    pushToken = currentToken;
                    $.get("/api/CheckKey?key=" + currentToken, function(data) {
                        if(data.success == "on") {
                            $('#pushswitch').bootstrapToggle('enable');
                            changePushSwitch("on");
                        } else if(data.success == "off") {
                            $('#pushswitch').bootstrapToggle('enable');
                            changePushSwitch("off");
                        } else {
                            console.log("Key an den Server uebermittelt");
                        }
                    });
                } else {
                    console.log('No Instance ID token available. Request permission to generate one.');
                }
            })
            .catch(function(err) {
                console.log('An error occurred while retrieving token. ', err);
            });
        })
        .catch(function(err) {
            console.log('Unable to get permission to notify.', err);
        });
    } else {
        //push is not available
    }
}

function enablePush() {
    $('#pushswitch').bootstrapToggle('disable');
    $.get("/api/RegisterKey?key=" + pushToken, function(data) {
        $('#pushswitch').bootstrapToggle('enable');
        if(data.success != "ok") {
            console.log("Fehler: " + data.error);
            changePushSwitch("off");
            alert("Fehler: " + data.error);
        } else {
            console.log("Key an den Server uebermittelt");
        }
    });
}

function disablePush() {
    $('#pushswitch').bootstrapToggle('disable');
    $.get("/api/RemoveKey?key=" + pushToken, function(data) {
        $('#pushswitch').bootstrapToggle('enable');
        if(data.success != "ok") {
            console.log("Fehler: " + data.error);
            changePushSwitch("on");
            alert("Fehler: " + data.error);
        } else {
            console.log("Key an den Server uebermittelt");
        }
    });
}

$(document).ready(function (){
    $("#pushswitch").bootstrapToggle("off");
    $("#pushswitch").bootstrapToggle("disable");
    getData();
    setInterval(getData, 5000);
    checkPushState();
    $("#pushswitch").change(handlePush);
});