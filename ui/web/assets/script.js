const digitRegex = /^\d+$/;
var functionsCount = 28;
var speed=0;
var address=3;
var directionForward=true;

function addFunctionRows(document) {
  for (let i = 1; i<=functionsCount; i=i+4) {
    if (i == NaN) {
      break;
    }

    var functionNumbers = [];
    for (let j = i; j < i+4 && j<=functionsCount; j++) {
      functionNumbers.push(j)
    }
    addFunctionRow(document, functionNumbers);
  }
}

function addFunctionRow(document, numbers) {
  var functions = document.getElementById('functions');
  var row = document.createElement('div');
  row.setAttribute('class', 'w3-row');
  functions.appendChild(row);

  for (let i = 0; i < numbers.length; i++) {
    addButton(document, row, numbers[i]);
  }
}

function addButton(document, row, number) {
  var buttonContainer = document.createElement('div');
  buttonContainer.setAttribute('class', 'w3-quarter btn-container');
  row.appendChild(buttonContainer);

  var button = document.createElement('button');
  // TODO update the ID to `f%d` instead of just `%d`
  button.setAttribute('id', number);
  button.setAttribute('value', number);
  button.textContent = 'F' + number;
  buttonContainer.appendChild(button);
}

document.addEventListener("keydown", function(e){
  var previousDirectionForward = directionForward;
  var previousSpeed = speed;

  console.log('key pressed: ' + e.which);
  // TODO: `which` is deprecated. figure out what to use instead
  if(e.which == 37) {
    // left arrow key
    directionForward = false;
  } else if(e.which == 38) {
    // up arrow key
    speed++;
  } else if(e.which == 39) {
    // right arrow key
    directionForward = true;
  } else if(e.which == 40) {
    // down arrow key
    speed--;
  }

  // if speed or direction changed, update the ui
  if (previousDirectionForward != directionForward || previousSpeed != speed) {
    updateSpeedRequest();
  }
});

$(document).ready(function(){
  const urlParams = new URLSearchParams(window.location.search);
  var addrFromURLParams = urlParams.get('address');
  if (addrFromURLParams !== null) {
    // update the global `address` var
    address = addrFromURLParams;
    $("#addr-input-box").val(address);
  }

  addFunctionRows(document);

  refreshState();
  window.setInterval(refreshState, 1000);

  var input = "";
  var len = 0;
  var curr="";
  var lastinput = "";
  var flag = 0; // for counting the number of decimals have already there in the input
  $("#speed-input-box").on('change', function(){
    var speedInput = $("#speed-input-box").val();
    console.log('changing speed to ' + speedInput);
    speed = speedInput;
    updateSpeedRequest();
  });
  $("#addr-input-box").on('change', function(){
    var addrInput = $("#addr-input-box").val();
    console.log('changing address to ' + addrInput);
    // update the global `address` var
    address = addrInput;

    // clear all active class buttons when address changes, let refresh state do it's thing
    $('button').removeClass('active');

    // update the browser url to have the address
    history.pushState({pageID: 'ddc-go ' + addrInput}, 'dcc-go ' + addrInput, '?address=' + addrInput);
  });
  $("button").on('click', function(){
      var button = $(this);
      curr = button.attr("value");
      console.log(curr);

      if (curr == 'power') {
        $.ajax({
          url: 'http://10.0.1.121:8080/power',
          success: function(data){
            console.log(JSON.stringify(data));
            var json = JSON.parse(data)
            if (json['power'] == true) {
              console.log('power on');
              button.addClass('active');
            } else {
              console.log('power off');
              button.removeClass('active');
            }
          },
          error: function(error){
            console.log("error: " + JSON.stringify(error));
          }
        });
      } else if (curr == 'stop') {
        speed = 0;
        $("#stop").addClass('active');
        $("#speed-input-box").val('0');
        $.ajax({
          url: 'http://10.0.1.121:8080/'+address+'/stop',
          success: function(data){
            console.log(JSON.stringify(data));
          },
          error: function(error){
            console.log("error: " + JSON.stringify(error));
          }
        });
      } else if (curr == '+' || curr == '-' || curr == 'forward' || curr == 'backward') {
        if (curr == '+') {
          speed++;
        } else if (curr == '-') {
          speed--;
        } else if (curr == 'forward') {
          directionForward = true;
        } else if (curr == 'backward') {
          directionForward = false;
        }
        updateSpeedRequest();
      } else if (curr.match(digitRegex)) {
        console.log(curr + " is a number");
        $.ajax({
          url: 'http://10.0.1.121:8080/'+address+'/function?function=' + curr,
          success: function(data){
            console.log(JSON.stringify(data));
            var json = JSON.parse(data)
            if (json['functions'][curr] == true) {
              console.log('function active');
              button.addClass('active')
            } else {
              console.log('function inactive');

              button.removeClass('active')
            }
          },
          error: function(error){
            console.log("error: " + JSON.stringify(error));
          }
        });
      }
   });
});

function updateSpeedRequest() {
  $("#stop").removeClass('active');
  if (speed <= 0) {
    speed = 0;
    $("#stop").addClass('active');
  }
  $("#speed-input-box").val(speed);

  if (directionForward) {
    $('#backward').removeClass('active');
    $('#forward').addClass('active');
  } else {
    $('#backward').addClass('active');
    $('#forward').removeClass('active');
  }

  $.ajax({
    url: 'http://10.0.1.121:8080/'+address+'/speed?forward=' + directionForward + '&speed=' + speed,
    success: function(data){
      console.log(JSON.stringify(data));
    },
    error: function(error){
      console.log("error: " + JSON.stringify(error));
    }
  });
}

function refreshState() {
  $.ajax({
    url: 'http://10.0.1.121:8080/state',
    success: function(data){
      var json = JSON.parse(data);

      // power is not throttle specific since it appplies to the whole track
      // update power first, then check for throttle address specific state.
      var power = json.power || false;
      if (power) {
        activeButton('power');
      } else {
        inactiveButton('power');
      }

      json.throttles = json.throttles || {};
      if(typeof json.throttles[address] === 'undefined') {
        activeButton('stop');
        return;
      }

      // update the global `speed` var
      speed = json.throttles[address].speed || 0;
      $("#speed-input-box").val(speed);
      if (speed == 0) {
        activeButton('stop');
      } else {
        inactiveButton('stop');
      }

      var direction = json.throttles[address].direction || 0;
      // update the global `directionForward` var
      directionForward = (direction == 1);
      if (directionForward) {
        inactiveButton('backward');
        activeButton('forward');
      } else {
        inactiveButton('forward');
        activeButton('backward');
      }

      // update functions states
      var functions = json.throttles[address].functions || {};
      for (let i = 0; i<=functionsCount; i++) {
        var functionEnabled = functions[i] || false;
        if (functionEnabled) {
          activeButton(i);
        } else {
          inactiveButton(i);
        }
      }
    },
    error: function(error){
      console.log("error on fetch state: " + JSON.stringify(error));
    }
  });
}

function activeButton(id) {
  var button = $("button#" + id);
  if (!button.hasClass('active')) {
    console.log("setting " + id + " to active");
    button.addClass('active');
  }
}

function inactiveButton(id) {
  var button = $("button#" + id);
  if (button.hasClass('active')) {
    console.log("setting " + id + " to inactive");
    button.removeClass('active');
  }
}
