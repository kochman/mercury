
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
function sendMessage(){
	console.log("sldnfds");
}