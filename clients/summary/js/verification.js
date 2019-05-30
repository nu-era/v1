const phoneNum = localStorage.getItem('phone');
console.log(phoneNum);


$('#phone-verification').submit(function(e) {
    var verificationURL = new URL("https://api.bfranzen.me/setup");

    e.preventDefault();
    var formInputs = $('#phone-verification :input');
    var values = {};
    formInputs.each(function() {
        values[this.name] = $(this).val();
    });
    values.phone = phoneNum;
    // console.log(values)

    var valJson = JSON.stringify(values);
    $.ajax({
        type: "PATCH",
        url: verificationURL,
        contentType: 'application/json',
        data: valJson,
        success: function( data, textStatus, response) {
            window.location.replace("./html/alert.html");
        },
        error: function(jqXhr, textStatus, errorThrown) {
            alert(jqXhr.responseText);
        }
    })
});