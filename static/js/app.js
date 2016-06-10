var containers = [];

function reload() {
    $.getJSON("/ping", function(data){
        if ($.inArray(data.instance, containers) == -1) {
            containers.push(data.instance);
        }

        for (var i=0; i<containers.length; i++) {
            var instanceName = containers[i];
            var el = $("#instance-" + instanceName);
            console.log(el);
            if (el.length == 0) {
                console.log("creating instance " + instanceName);
                console.log(el);
                var elData = '<div id="instance-' + instanceName + '" class="card"><div class="image"><img width="25%" height="25%" src="static/img/container.png"></div><div class="content"><h2 class="header">' + instanceName + '</h2></div></div>';
                $("div.container-group").append(elData);
                el = $("#instance-" + instanceName);
            } else {
                // "pulse"
                $(el).fadeToggle("fast", function(){
                    $(el).fadeToggle(250);
                });
            }
        }
    });
}

setInterval(reload, 2000);
