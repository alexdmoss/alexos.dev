// Get Parameters from some url
var getUrlParameter = function getUrlParameter(sPageURL) {
	var url = sPageURL.split('?');
	var obj = {};
	if (url.length == 2) {
		var sURLVariables = url[1].split('&'),
			sParameterName,
			i;
		for (i = 0; i < sURLVariables.length; i++) {
			sParameterName = sURLVariables[i].split('=');
			obj[sParameterName[0]] = sParameterName[1];
		}
		return obj;
	} else {
		return undefined;
	}
};


// Execute actions on images generated from Markdown pages
var images = $("div#body-inner img").not(".inline");
// Wrap image inside a featherlight (to get a full size view in a popup)
images.wrap(function () {
	var image = $(this);
	if (!image.parent("a").length) {
		return "<a href='" + image[0].src + "' data-featherlight='image'></a>";
	}
});

// Change styles, depending on parameters set to the image
images.each(function (index) {
	var image = $(this)
	var o = getUrlParameter(image[0].src);
	if (typeof o !== "undefined") {
		var h = o["height"];
		var w = o["width"];
		var c = o["classes"];
		image.css("width", function () {
			if (typeof w !== "undefined") {
				return w;
			} else {
				return "auto";
			}
		});
		image.css("height", function () {
			if (typeof h !== "undefined") {
				return h;
			} else {
				return "auto";
			}
		});
		if (typeof c !== "undefined") {
			var classes = c.split(',');
			for (i = 0; i < classes.length; i++) {
				image.addClass(classes[i]);
			}
		}
	}
});



(function ($) {

	skel.breakpoints({
		xlarge: '(max-width: 1680px)',
		large: '(max-width: 1280px)',
		medium: '(max-width: 980px)',
		small: '(max-width: 736px)',
		xsmall: '(max-width: 480px)'
	});

	$(function () {

		var $window = $(window),
			$body = $('body'),
			$menu = $('#menu'),
			$shareMenu = $('#share-menu'),
			$sidebar = $('#sidebar'),
			$main = $('#main');

		// // Fix: Placeholder polyfill.
		// 	$('form').placeholder();

		// // Prioritize "important" elements on medium.
		// 	skel.on('+medium -medium', function() {
		// 		$.prioritize(
		// 			'.important\\28 medium\\29',
		// 			skel.breakpoint('medium').active
		// 		);
		// 	});

		// // IE<=9: Reverse order of main and sidebar.
		// 	if (skel.vars.IEVersion <= 9)
		// 		$main.insertAfter($sidebar);

		// $menu.appendTo($body);
		// $shareMenu.appendTo($body);

		// $menu.panel({
		// 	delay: 500,
		// 	hideOnClick: true,
		// 	hideOnEscape: true,
		// 	hideOnSwipe: true,
		// 	resetScroll: true,
		// 	resetForms: true,
		// 	side: 'right',
		// 	target: $body,
		// 	visibleClass: 'is-menu-visible'
		// });

		// $shareMenu.panel({
		// 	delay: 500,
		// 	hideOnClick: true,
		// 	hideOnEscape: true,
		// 	hideOnSwipe: true,
		// 	resetScroll: true,
		// 	resetForms: true,
		// 	side: 'right',
		// 	target: $body,
		// 	visibleClass: 'is-share-visible'
		// });



		// Search (header).
		var $search = $('#search'),
			$search_input = $search.find('input');
		$body
			.on('click', '[href="#search"]', function (event) {
				event.preventDefault();
				// Not visible?
				if (!$search.hasClass('visible')) {
					// Reset form.
					$search[0].reset();
					// Show.
					$search.addClass('visible');
					// Focus input.
					$search_input.focus();
				}
			});
		$search_input
			.on('keydown', function (event) {
				if (event.keyCode == 27)
					$search_input.blur();
			})
			.on('blur', function () {
				window.setTimeout(function () {
					$search.removeClass('visible');
				}, 100);
			});

	});

})(jQuery);


// anchor links for headings
jQuery(document).ready(function () {
	var text, clip = new Clipboard('.anchor');
	$("h2[id],h3[id],h4[id],h5[id],h6[id]").append(function (index, html) {
		var element = $(this);
		var url = encodeURI(document.location.origin + document.location.pathname);
		var link = url + "#" + element[0].id;
		return " <span class='anchor' data-clipboard-text='" + link + "'>" +
			"<i class='fas fa-link fa-lg'></i>" +
			"</span>"
			;
	});

	$(".anchor").on('mouseleave', function (e) {
		$(this).attr('aria-label', null).removeClass('tooltipped tooltipped-s tooltipped-w');
	});

	clip.on('success', function (e) {
		e.clearSelection();
		$(e.trigger).attr('aria-label', 'Link copied to clipboard!').addClass('tooltipped tooltipped-s');
	});

});