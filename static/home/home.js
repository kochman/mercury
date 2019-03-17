
// Function to validate that the message form is not blank 
function validateForm(){
	let unique_id = document.getElementById("unique_id").value;	
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
	
	$.ajax({
		url: "localhost:3000",
		method: "POST",
		data: {
			TargetUserID: $("#unique_id").val(),
			Message: $("#message").val()
		},
		success: function (result) {
			console.log("success")
			console.log(result)
		}
	})
}


