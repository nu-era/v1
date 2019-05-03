// Glowing animation
var glow = $('.glowing');
setInterval(function () {
  glow.toggleClass('glow');
}, 1000);


/*
////////////////////////////
////////////////////////////
  Connect to websocket
///////////////////////////
///////////////////////////
*/

const hostUrl = "wss://trivia.bfranzen.me/v1/ws?auth="
const auth = localStorage.getItem('auth');
const user = JSON.parse(localStorage.getItem('user'));
const triviaUrl = "https://trivia.bfranzen.me/v1/trivia"
var currentLobby;




socket = new WebSocket(hostUrl + auth)
socket.onopen = function () {
  console.log('Websocket connected');
}

socket.onmessage = function (event) {
  var receivedMsg = JSON.parse(event.data);
  var type = receivedMsg.type;

  if (type === 'lobby-new') {
    currentLobby = receivedMsg.lobby.lobbyId;
    createAddLobby(receivedMsg.lobby);
  } else if (type === 'game-start') {
    startGameHandler(receivedMsg);
  } else if (type === 'game-question') {
    newQuestionHandler(receivedMsg);
  } else if (type === 'game-won') {
    alert('Congrats! You won!')
    leaveGameHandler();
  } else if (type === 'game-over') {
    alert('Bummer, you lost...');
    leaveGameHandler();
  } else if (type === 'lobby-add') {
    updatePlayers(receivedMsg.userIDs);
  } else if (type === 'lobby-update') {
    console.log(receivedMsg);
    getAllLobbies();
    updateOptions(receivedMsg.lobby);
  }
}



socket.onclose = function (event) {
  alert("Please Sign-In")
  window.location.replace("https://jmatray.me/app");
}







/*
////////////////////////////
////////////////////////////
  DOM Functions
///////////////////////////
///////////////////////////
*/

// Switch from landing page to lobby
function switchToLobby() {
  var landing = document.querySelector(".landing");
  var game = document.querySelector(".game");
  if (landing.style.display === "none") {
    landing.style.display = "flex";
    game.style.display = "none"
  } else {
    landing.style.display = "none";
    game.style.display = "flex"
  }
};

// Switch from lobby to game
function switchToGame() {
  $('.waiting').hide();
  $('.playing').show();
  $('.board').css('height', 'auto')
}



// Switch from game to landing
function leaveGameHandler() {
  $('.game').hide();
  $('.waiting').show();
  $('.playing').hide();
  $('.landing').show();
  $('.board').css('height', '80vh');
  getAllLobbies();
}



/*
////////////////////////////
////////////////////////////
  Landing Page Functions
///////////////////////////
///////////////////////////
*/

function createAddLobby(lob) {
  var newLob = document.createElement('DIV');
  newLob.setAttribute('class', 'lobby')
  var img = document.createElement('IMG')
  img.setAttribute('class', 'lobby-pic')
  img.setAttribute('src', '/app/public/imgs/Drawing.png')

  var creat = document.createElement('P')
  creat.innerHTML = "Creator: " + lob.creator.userName

  var cat = document.createElement('P')
  cat.innerHTML = "Category: " + lob.options.category

  var diff = document.createElement('P')
  diff.innerHTML = "Difficulty: " + lob.options.difficulty

  var join = document.createElement('button')
  join.setAttribute('id', lob.lobbyId)
  join.setAttribute('class', 'join button')

  if (!lob.inProgress) {
    $(join).text("Join")
    $(join).on('click', joinGame)
  } else {
    $(join).text("In Progress").addClass("disabled");

  }

  newLob.appendChild(img)
  newLob.appendChild(creat)
  newLob.appendChild(cat)
  newLob.appendChild(diff)
  newLob.appendChild(join)
  $('.lobbies').append(newLob);
}

function joinGame() {
  /*
  Step 1: Post request to /trivia/id
  Step 2: Wait for response with lobby struct
  Step 3: Switch to show lobby
*/
  var lobbyId = $(this).attr('id');
  currentLobby = lobbyId;
  $.ajax({
    type: "POST",
    url: triviaUrl + '/' + lobbyId + '?type=add',
    contentType: 'application/json',
    headers: {
      "Authorization": auth
    },
    success: function (data, textStatus, response) {
      $('.form-control').prop('disabled', true);
      switchToLobby()
    },
    error: function (jqXhr, textStatus, errorThrown) {
      alert(jqXhr.responseText);
    }
  })

}

function createGame() {
  /*
    Step 1: Post request to /trivia
    Step 2: Wait for response of new lobby created, with lobby struct
    Step 3: Track the creator somehow
    Step 4: Switch to show lobby
  */
  var category = parseInt($("#category").val());
  var diff = $('#difficulty').val()
  var players = parseInt($('#players').val())
  var questions = parseInt($('#questions').val());

  var options = { "numQuestions": questions, "category": category, "difficulty": diff, "maxPlayers": players }
  var optionsJson = JSON.stringify(options);
  $.ajax({
    type: "POST",
    url: triviaUrl,
    contentType: 'application/json',
    headers: {
      "Authorization": auth
    },
    data: optionsJson,
    success: function (data, textStatus, response) {
      $('.form-control').prop('disabled', false);
      $('.num-players').html('1');
      switchToLobby()
    },
    error: function (jqXhr, textStatus, errorThrown) {
      alert(jqXhr.responseText);
    }
  })


}
$('.new-lobby').on('click', createGame);




// Get lobbies

getAllLobbies();
function getAllLobbies() {
  /*
    Step 1: remove all elements from lobbies DOM
    Step 2: send get request to get all lobby structs
    Step 3: Wait for response with all lobbies
    Step 4: Loop through each lobby, adding to DOM
  */
  $('.lobby').remove();
  $.ajax({
    type: "GET",
    url: triviaUrl,
    headers: {
      "Authorization": auth
    },
    dataType: 'text',
    success: function (data, textStatus, response) {
      var lobs = JSON.parse(data);
      var values = Object.values(lobs);
      values.forEach(function(item) {
        createAddLobby(item);
      })
      
    },
    error: function (jqXhr, textStatus, errorThrown) {

      console.log(errorThrown);
    }
  })  

}

var lobbies = document.querySelectorAll(".join");
var gameStart = document.querySelectorAll(".start-game");
var leaveGame = document.querySelector(".leave-game");
leaveGame.addEventListener('click', leaveGameHandler, false);
for (var i = 0; i < gameStart.length; i++) {
  gameStart[i].addEventListener('click', startGameHandler, false)
}


/*
////////////////////////////
////////////////////////////
  Inside Lobby Functions
///////////////////////////
///////////////////////////
*/
$(".form-control").on('change', function () {
  var category = parseInt($("#category").val());
  var diff = $('#difficulty').val()
  var players = parseInt($('#players').val())
  var questions = parseInt($('#questions').val());

  var options = { "numQuestions": questions, "category": category, "difficulty": diff, "maxPlayers": players }
  var optionsJson = JSON.stringify(options);
  $.ajax({
    type: "PATCH",
    url: triviaUrl  + '/' + currentLobby,
    contentType: 'application/json',
    headers: {
      "Authorization": auth
    },
    data: optionsJson,
    success: function (data, textStatus, response) {
      updateOptions()
    },
    error: function (jqXhr, textStatus, errorThrown) {
      alert(jqXhr.responseText);
    }
  })
});

function updatePlayers(players) {
  $('.num-players').html(players.length);
}

function updateOptions(msg) {
  
  var newOptions = msg.options;
  $('#' + msg.lobbyId)
  $('#category').val(newOptions.category);
  $('#difficulty').val(newOptions.difficulty);
  $('#maxPlayers').val(newOptions.maxPlayers);
  $('#questions').val(newOptions.numQuestions); 
}


// Handle the start of a new game

function startGameHandler(msg) {
  /*
    Step 1: Run newQuestionHandler
    Step 2: Switch to show game DOM
    Step 3: (Might have to wait for question)
  */
  $('.timer').html('');
  currentLobby = msg.lobby.lobbyId;
  $('#' + msg.lobby.lobbyId).addClass("disabled").text("In Progress").off('click', switchToLobby);
  switchToGame()
}






/*
////////////////////////////
////////////////////////////
  Inside Game Functions
///////////////////////////
///////////////////////////
*/

// Handle new question messages from server

var currentQuestion;
function newQuestionHandler(msg) {
  currentQuestion = msg.question.questionID;
  $('.num-players').html(msg.userIDs.length);
  // TEMPORARY TIME FOR QUESTION
  var now = new Date().getTime()
  var timeLeft = 30

  clearInterval(timerId);
  $('.timer').html('');

  var timerId = setInterval(countdown, 1000);

  function countdown() {
    $('.timer').html(timeLeft + 's')
    if (timeLeft == 0) {
      clearTimeout(timerId);
      doSomething();
    } else {
      timeLeft--;
    }
  }
  $('.question-text').html(msg.question.question)
  var choices = msg.question.choices;
  var answersDom = document.querySelector('.answers');
  answersDom.innerHTML = '';
  choices.forEach(function (item) {
    var newButton = document.createElement('button');
    newButton.setAttribute('class', 'answer button button-block');
    newButton.innerHTML = item;
    $(newButton).on('click', submitAnswer);
    answersDom.appendChild(newButton);
  })

}

function submitAnswer() {
  /*
    Step 1: Send post request to answer
    Step 2: After successful post request, show message that answer was submitted
    Step 3: Disable answer buttons
  */
  var chosen = $(this).html();
  var answer = { "lobbyId": currentLobby, "userId": user.id, "questionId": currentQuestion, "answer": chosen }
  var answerJson = JSON.stringify(answer);
  $.ajax({
    type: "POST",
    url: triviaUrl + '/' + currentLobby + '?type=answer',
    contentType: 'application/json',
    headers: {
      "Authorization": auth
    },
    data: answerJson,
    success: function (data, textStatus, response) {
      $('.answer').addClass('disabled');
      $('.answer').off('click', submitAnswer);
      alert('Answer Submitted');
    },
    error: function (jqXhr, textStatus, errorThrown) {
      alert(jqXhr.responseText);
    }
  })


}