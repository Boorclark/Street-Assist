$(document).ready(function(){
  $('.dropdown-item').click(function(){
    $('.dropdown-item').removeClass('active'); // remove active class from all dropdown items
    $(this).addClass('active'); // add active class to selected dropdown item
    $('#stateDropdownButton').text($(this).text());
  });

  $("#submitButton").click(function(event) {
      //event.preventDefault();
      var state = $(".dropdown-menu .dropdown-item.active").data("value");
      var city = $("#cityInput").val();
      $.ajax({
        url: "/information.html",
        type: "POST",
        data: {"state": state, "city": city},
        success: function(response) {
          // handle successful response
          $("#shelterList").html(response);
        },
        error: function(xhr, status, error) {
          // handle error response
        }
      });
    });

});