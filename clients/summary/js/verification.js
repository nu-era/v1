const phoneNum = localStorage.getItem('phone');
console.log(phoneNum);


$('#phone-verification').click(function(e) {
    var verificationURL = new URL("https://api.bfranzen.me/device-info");

    e.preventDefault();
    var formInputs = $('#phone-verification :input');
    var values = {};
    formInputs.each(function() {
        values[this.name] = $(this).val();
    });
    values.phone = phoneNum;
    console.log(values)
    // var values = {};
    // formInputs.each(function() {
    //     values[this.name] = $(this).val();
    // });
    // var valJson = JSON.stringify(values);
    // $.ajax({
    //     type: "POST",
    //     url: retUserUrl,
    //     contentType: 'application/json',
    //     data: valJson,
    //     success: function( data, textStatus, response) {
    //         var auth = response.getResponseHeader('Authorization');
    //         var userData = JSON.stringify(data);
    //         localStorage.setItem('auth', auth);
    //         localStorage.setItem('user', userData);
    //         switchToVerification();
    //         window.location.replace("./html/alert.html");
    //     },
    //     error: function(jqXhr, textStatus, errorThrown) {
    //         alert(jqXhr.responseText);
    //     }
    // })
});