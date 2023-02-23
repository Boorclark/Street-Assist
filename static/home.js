$(document).ready(function(){

    $('.dropdown-item').click(function(){
        $('#stateDropdownButton').text($(this).text());
    });

});