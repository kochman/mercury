// Function to validate that the message form is not blank 
function validateForm(){
	let message = document.getElementById("message").value;
	if (message.length == 0) {
		alert("invalid message")
		return false;
	}
	sendMessage();
	return true;
}

// Function to send the message 
function sendMessage() {
	fetch("/send", {
		method: "POST",
		body: JSON.stringify({
			TargetUserID: Number($("#unique_id").val()),
			Message: $("#message").val()
		})
	}).then(() => {
		console.log("here")
		$("#messagesentinfo").show()
		setTimeout(()=> {
			$("#messagesentinfo").hide()
		}, 1000)
	})
}

// on page load
// get all available contacts and append 
// them to the list of possible message receivers
$(document).ready(function () {
	$("#messagesentinfo").hide()
	fetch("/api/contacts/all").then(function (response) {
		return response.json();
	}).then(function (data) {
		for (var i = 0; i < data.length; i++){
			$("#unique_id").append("<option value=" + data[i].ID + ">" + data[i].Name + " - " + data[i].ID + "</option>");
			
		}
		
		//parse arguments in the bar
		// may be used for future url params
		$.urlParam = function (name) {
			var results = new RegExp('[\?&]' + name + '=([^&#]*)')
							.exec(window.location.search);
			return (results !== null) ? results[1] || 0 : false;
		}
		if ($.urlParam('user')) {
			$('#unique_id').val($.urlParam('user'))
		}
	})
})



