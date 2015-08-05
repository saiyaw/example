$(document).ready(function() {

	$('#btnping').click(function() {
		$.ajax({
			type : "POST",
			url : "/dbping",
			success : function(result) {
				alert(result);

			}
		});
	});

	$('#btndemo').click(function() {
		$.ajax({
			type : "POST",
			url : "/dbgetdemo",
			data : {
				"date" : $('#indate').val(),
				"symbol" : $('#insymbol').val()
			},
			success : function(result) {
				$('#ptext').text(result);
			}
		});
	});

	$('#btnstock').click(function() {
		$.ajax({
			type : "POST",
			url : "/getstock",
			data : {
				"date" : $('#indate').val(),
				"symbol" : $('#insymbol').val()
			},
			success : function(result) {
				// alert(result);

				$('#ptext').text(result);

			}
		});
	});

	$('#btncopy').click(function() {
		$.ajax({
			type : "POST",
			url : "/datacopy",
			success : function(result) {
				// alert(result);

				$('#ptext').text(result);

			}
		});
	});

});