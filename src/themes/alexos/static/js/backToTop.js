$(function () {
    // browser window scroll (in pixels) after which the "back to top" link is shown
    let offset = 300,
        //browser window scroll (in pixels) after which the "back to top" link opacity is reduced
        offset_opacity = 1200,
        //duration of the top scrolling animation (in ms)
        scroll_top_duration = 700,
        //grab the "back to top" link
        $backToTop = $('#back-to-top');
    let $shareNav = $('#share-nav');

    //hide or show the "back to top" link
    $(window).scroll(function () {
        if ($(this).scrollTop() > offset) {
            $backToTop.fadeIn();
            $shareNav.fadeIn();
        } else {
            $backToTop.fadeOut();
            $backToTop.removeClass('button-fade-out');
            $shareNav.fadeOut();
        }

        if ($(this).scrollTop() > offset_opacity) {
            $backToTop.addClass('button-fade-out');
        }
    });

    //smooth scroll to top
    $backToTop.on('click', function (event) {
        event.preventDefault();
        $('body,html').animate({
            scrollTop: 0,
        }, scroll_top_duration
        );
    });

});
