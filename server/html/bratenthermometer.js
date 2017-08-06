
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