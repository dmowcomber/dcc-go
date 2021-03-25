function addFunctionRows(document) {
  for (let i = 1; i<=32; i=i+4) {
    if (i == NaN) {
      break;
    }

    var functionNumbers = [];
    for (let j = i; j < i+4 && j<=28; j++) {
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
  button.setAttribute('id', number);
  button.setAttribute('value', number);
  button.textContent = 'F' + number;
  buttonContainer.appendChild(button);
}

const digitRegex = /^\d+$/;
var speed=0;
var address=3;
var directionForward=true;

document.addEventListener("keydown", function(e){
  console.log('key pressed: ' + e.which);
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
  updateSpeedRequest();
});

$(document).ready(function(){
  addFunctionRows(document);

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
    address = addrInput;
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
              button.addClass('active')
            } else {
              console.log('power off');
              button.removeClass('active')
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
