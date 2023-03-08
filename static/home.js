$(document).ready(function(){
  $('.dropdown-item').click(function(){
    $('.dropdown-item').removeClass('active');
    $(this).addClass('active');
    $('#stateDropdownButton').text($(this).text());
  });

  var all_urls = ["/information/shelters", "/information/food"];

  $("#submitButton").click(function(event) {
    event.preventDefault();

    var state = $(".dropdown-menu .dropdown-item.active").data("value");
    var city = $("#cityInput").val();

    $.each(all_urls, function(index, url) {
      $.ajax({
        url: url,
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
});
