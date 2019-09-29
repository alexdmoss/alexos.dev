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

function mobileMenu() {
	var x = document.getElementById("header-menu");
	if (x.className === "links") {
	  x.className += " responsive";
	} else {
	  x.className = "links";
	}
  }

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


// clipboard
var clipInit = false;
$('code').each(function () {
	var code = $(this),
		text = code.text();

	if (text.length > 20) {
		if (!clipInit) {
			var text, clip = new Clipboard('.copy-to-clipboard', {
				text: function (trigger) {
					text = $(trigger).prev('code').text();
					return text.replace(/^\$\s/gm, '');
				}
			});

			var inPre;
			clip.on('success', function (e) {
				e.clearSelection();
				inPre = $(e.trigger).parent().prop('tagName') == 'PRE';
				$(e.trigger).attr('aria-label', 'Copied to clipboard!').addClass('tooltipped tooltipped-' + (inPre ? 'w' : 's'));
			});

			clip.on('error', function (e) {
				inPre = $(e.trigger).parent().prop('tagName') == 'PRE';
				$(e.trigger).attr('aria-label', fallbackMessage(e.action)).addClass('tooltipped tooltipped-' + (inPre ? 'w' : 's'));
				$(document).one('copy', function () {
					$(e.trigger).attr('aria-label', 'Copied to clipboard!').addClass('tooltipped tooltipped-' + (inPre ? 'w' : 's'));
				});
			});

			clipInit = true;
		}

		code.after('<span class="copy-to-clipboard" title="Copy to clipboard" />');
		code.next('.copy-to-clipboard').on('mouseleave', function() {
			$(this).attr('aria-label', null).removeClass('tooltipped tooltipped-s tooltipped-w');
		});
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

		var $body = $('body');

		// Search (header).
		var $search = $('#search');
		var $search_input = $search.find('input');
		var $search_results = $('#search-results');

		$body.on('click', '[href="#search"]', function (event) {
			event.preventDefault();
			if (!$search.hasClass('visible')) {
				$search[0].reset();
				$search.addClass('visible');
				$search_input.focus();
			}
		});
		$search_input
			.on('keydown', function (event) {
				if (event.keyCode == 27) {
					$search_input.blur();
					window.setTimeout(function () {
						$search_results.removeClass('visible');
					}, 100);
				}
			})
			.on('blur', function () {
				window.setTimeout(function () {
					$search.removeClass('visible');
				}, 100);
			});

		// Share Menu (header).
		var $share = $('#share');
		$body
			.on('click', '[href="#share-menu"]', function (event) {
				event.preventDefault();
				if (!$share.hasClass('visible')) {
					$share.addClass('visible');
				}
			})
			.on('keydown', function (event) {
				if (event.keyCode == 27)
					window.setTimeout(function () {
						$share.removeClass('visible');
					}, 100);
			})
			.on('click', '[href="#close-share"]', function (event) {
				event.preventDefault();
				window.setTimeout(function () {
					$share.removeClass('visible');
				}, 100);
			});

	});

})(jQuery);

jQuery(document).ready(function () {

	// anchor links for headings

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
		$(e.trigger).attr('aria-label', 'Link copied to clipboard!').addClass('tooltipped tooltipped-se');
	});

    /** 
    * Fix anchor scrolling that hides behind top nav bar
    * Courtesy of https://stackoverflow.com/a/13067009/28106
    **/
	(function (document, history, location) {
		var HISTORY_SUPPORT = !!(history && history.pushState);

		var anchorScrolls = {
			ANCHOR_REGEX: /^#[^ ]+$/,
			OFFSET_HEIGHT_PX: 80,

			// Establish events, and fix initial scroll position if a hash is provided.
			init: function () {
				this.scrollToCurrent();
				$(window).on('hashchange', $.proxy(this, 'scrollToCurrent'));
				$('body').on('click', 'a', $.proxy(this, 'delegateAnchors'));
			},

			// Return the offset amount to deduct from the normal scroll position. Modify as appropriate to allow for dynamic calculations
			getFixedOffset: function () {
				return this.OFFSET_HEIGHT_PX;
			},

			// If the provided href is an anchor which resolves to an element on the page, scroll to it
			scrollIfAnchor: function (href, pushToHistory) {
				var match, anchorOffset;

				if (!this.ANCHOR_REGEX.test(href)) {
					return false;
				}

				match = document.getElementById(href.slice(1));

				if (match) {
					anchorOffset = $(match).offset().top - this.getFixedOffset();
					$('html, body').animate({ scrollTop: anchorOffset });

					// Add the state to history as-per normal anchor links
					if (HISTORY_SUPPORT && pushToHistory) {
						history.pushState({}, document.title, location.pathname + href);
					}
				}

				return !!match;
			},
			// Attempt to scroll to the current location's hash.
			scrollToCurrent: function (e) {
				if (this.scrollIfAnchor(window.location.hash) && e) {
					e.preventDefault();
				}
			},
			// If the click event's target was an anchor, fix the scroll position.
			delegateAnchors: function (e) {
				var elem = e.target;

				if (this.scrollIfAnchor(elem.getAttribute('href'), true)) {
					e.preventDefault();
				}
			}
		};

		$(document).ready($.proxy(anchorScrolls, 'init'));
	})(window.document, window.history, window.location);

});