
function progress(data){
    $("#fleischCards").html("");
    $("#smokerCards").html("");
    for(var i in data.success.Temps) {
    var t = data.success.Temps[i];
    if(t.Meter.Type == "fleisch") {
        if(t.Temp >= "95") {
        $("#fleischCards").append("<div class='col-sm-2'><div class='tempcard tempgreencard'>"+t.Temp+" °C</div></div>");                       
        } else {
        $("#fleischCards").append("<div class='col-sm-2'><div class='tempcard'>"+t.Temp+" °C</div></div>");                       
        }
    } else {
        if(t.Temp <= "100" || t.Temp > "130") {
        $("#smokerCards").append("<div class='col-sm-2'><div class='tempcard tempredcard'>"+t.Temp+" °C</div></div>");
        } else {
        $("#smokerCards").append("<div class='col-sm-2'><div class='tempcard'>"+t.Temp+" °C</div></div>");
        }
    }
    }
}
function getData() {
    $.get("/api/GetTemps", progress);
}
$(document).ready(function (){
    getData();
    setInterval(getData, 5000);
});

let swRegistration = null;
// Initialize Firebase

var config = {
    apiKey: "AIzaSyD3TZlq2TIS9C719NhvekAL2LfA-Fsr3h4",
    authDomain: "porkmeter-b3cf2.firebaseapp.com",
    databaseURL: "https://porkmeter-b3cf2.firebaseio.com",
    projectId: "porkmeter-b3cf2",
    storageBucket: "porkmeter-b3cf2.appspot.com",
    messagingSenderId: "615966299981"
};

firebase.initializeApp(config);

// Retrieve Firebase Messaging object.
const messaging = firebase.messaging();

messaging.requestPermission()
.then(function() {
    console.log('Notification permission granted.');
    messaging.getToken()
    .then(function(currentToken) {
        if (currentToken) {
            console.log("current token", currentToken)
            $.get("/api/RegisterKey?key=" + currentToken, function(data) {
                if(data.success != "ok") {
                    console.log("Fehler: " + data.error);
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