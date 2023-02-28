$(document).ready(function(){
  $('.dropdown-item').click(function(){
      $('#stateDropdownButton').text($(this).text());
  });

  $("#submitButton").click(function(event) {
      event.preventDefault();
      var state = $("#stateDropdown .dropdown-item.active").data("value");
      var city = $("#cityInput").val();
      $.ajax({
        url: "/",
        type: "POST",
        data: {state: state, city: city},
        success: function(response) {
          // handle successful response
        },
        error: function(xhr, status, error) {
          // handle error response
        }
      });
    });

});