var containers = [];

function reload() {
    $.ajax({
          url: "/ping",
          dataType: "json",
          cache: false,
          success: function(data){
              if ($.inArray(data.instance, containers) == -1) {
                  containers.push(data.instance);
              }

              for (var i=0; i<containers.length; i++) {
                  var instanceName = containers[i];
                  var el = $("#instance-" + instanceName);
                  if (el.length == 0) {
                      console.log("creating instance " + instanceName);
                      var elData = '<div id="instance-' + instanceName + '" class="card container-instance"><div class="image"><img width="25%" height="25%" src="static/img/container.png"></div><div class="content"><h2 class="header">' + instanceName + '</h2></div></div>';
                      $("div.container-group").append(elData);
                      el = $("#instance-" + instanceName);
                  }
              }

              // "pulse"
              var el = $("#instance-" + data.instance);
              $(el).fadeToggle(250, function(){
                  $(el).fadeToggle(250);
              });
          }
    });
}

setInterval(reload, 1000);
