
// Function to validate that the message form is not blank 
function validateForm(){
	let unique_id = document.getElementById("unique_id").value;	
	let message = document.getElementById("message").value;
	if (message.length == 0){
		return false;
	}
	sendMessage();
	return true;
}

// Function to send the message 
function sendMessage() {
	
	$.ajax({
		url: "127.0.0.1:3000/send",
		method: "POST",
		data: {
			TargetUserID: $("#unique_id"),
			Message: $("#message")
		},
		dataType:"json",
		success: function (result) {
			console.log("success")
			console.log(result)
		}
	})
	console.log("asdf")
	return false;
}


