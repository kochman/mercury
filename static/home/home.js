
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
	console.log($("#unique_id").val())
	fetch("/send", {
		method: "POST",
		body: JSON.stringify({
			TargetUserID: Number($("#unique_id").val()),
			Message: $("#message").val()
		})
	}).then(() => {
		console.log("dopne")
	})
	// $.ajax({
	// 	url: "localhost:3000",
	// 	method: "POST",
	// 	data: {
	// 		TargetUserID: $("#unique_id").val(),
	// 		Message: $("#message").val()
	// 	},
	// 	success: function (result) {
	// 		console.log("success")
	// 		console.log(result)
	// 	}
	// })
}


