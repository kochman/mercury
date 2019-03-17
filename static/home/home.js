
// Function to validate that the message form is not blank 
function validateForm(){
	let unique_id = document.getElementById("unique_id").value;	
	let message = document.getElementById("message").value;

	if (message.length == 0){
		console.log("too short");
	}
	
	return true;


}


// Function to send the message 
function sendMessage(){

}