$(document).ready(function(){
  $('.dropdown-item').click(function(){
    $('.dropdown-item').removeClass('active'); // remove active class from all dropdown items
    $(this).addClass('active'); // add active class to selected dropdown item
    $('#stateDropdownButton').text($(this).text());
  });

  all_urls = "/information/shelters", "/information/food"
  $("#submitButton").click(function(event) {
    //event.preventDefault();
    var state = $(".dropdown-menu .dropdown-item.active").data("value");
    var city = $("#cityInput").val();
    $.ajax({
      url: all_urls,
      type: "POST",
      data: {"state": state, "city": city},
      success: function(response) {
        window.location.href = "resources.html";
      },
      error: function(xhr, status, error) {
        // handle error response
      }
    });
  });
  

});