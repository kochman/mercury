
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
		console.log("done")
	})
}

$(document).ready(function () {
	console.log("loaded")

	fetch("/api/contacts/all").then(function (response) {
		return response.json();
	}).then(function (data) {
		console.log(data);
		for (var i = 0; i < data.length; i++){
			$("#unique_id").append("<option value=" + data[i].ID + ">"+data[i].Name +"</option>");
		}
	})
})


