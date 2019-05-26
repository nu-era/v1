const data = localStorage.getItem('phone');
console.log(data);
$('#phone-verification').submit(function(e) {
    let numTo = data.phone;
    var verificationURL = new URL("https://api.bfranzen.me/device-info");
    console.log(numTo);
    // e.preventDefault();
    // var formInputs = $('#phone-verification :input');
    // console.log(formInputs);
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