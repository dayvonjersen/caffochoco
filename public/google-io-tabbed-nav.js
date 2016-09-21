$(function() {
  
  // Setup Line
  $('.papertabs').append("<li id='papertabs-line'></li>");
  var $line = $('#papertabs-line');
  var $activeItem =  $('.papertabs .active').parent();
  var $activeX = $('.papertabs .active').parent().position().left;
	$line.width($activeItem.width()).css("left", $activeX);

  // Click Event
  $('.papertabs a').click( function(e){

    var $el = $(this);
    var $offset = $el.offset();
   	var $clickX = e.pageX - $offset.left;
   	var $clickY = e.pageY - $offset.top;
    var $parentX = $el.parent().position().left;
    var $elWidth = $el.parent().width();

    e.preventDefault();
    
    $('.papertabs .active').removeClass('active');
    $el.addClass('pressed active');

    $el.find('.circle').css({
    	left: $clickX + 'px', top: $clickY + 'px'
    });

    $line.animate({
      left: $parentX, width: $elWidth
    });
    
    $el.on("animationend webkitAnimationEnd oAnimationEnd MSAnimationEnd", function(){
      $el.removeClass('pressed').addClass('focused');
      setTimeout( function(){
      	$el.removeClass('focused');
      }, 800);
    });

  });
  
}); // End Of Document Ready Function