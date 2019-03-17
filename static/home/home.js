
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
let result ="";
let full_result="";
// on page load
// get all available contacts and append 
// them to the list of possible message receivers
$(document).ready(function () {
	fetch("/api/contacts/all").then(function (response) {
		return response.json();
	}).then(function (data) {
		for (var i = 0; i < data.length; i++){
			$("#unique_id").append("<option value=" + data[i].ID + ">"+data[i].Name +" - " + data[i].ID+ "</option>");
		}
	})

	fetch("/api/self").then(function (response) {	
		return response.text();
	}).then(function (data) {
		let count = 0;
		let i = 26;
		while (count <= 10){
			result += data[i];
			i++;
			count++;
		}
		for (i = 0 ; i < data.length; i++ ){
			full_result+=data[i];
		}

		result += "...";
		$("#personal-key").text(result);
	})	
})

$( "#personal-key" ).hover(
	function() {
		$( this ).append(full_result);
	}, function() {
		$( this ).find($(full_result-result)).remove();
	}
);