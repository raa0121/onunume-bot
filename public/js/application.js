$(function() {
    $('.content').each(function(i, e) {
        var text;
        text = $(e).text().replace(/\bhttps?:\/\/\S+(?:jpg|png|gif|JPG|PNG|GIF)(\?[^< \t\r\n]+|$)/, function(match) {
            return '<img class="osusume-image thumbnail lazy" src="./img/dummy.jpg" data-original="' + match.replace('&', '&amp;').replace('<', '&lt;').replace('>', '&gt;').replace('"', '&quot;').replace("'", '&apos;') + '">';
        });
        return $(e).html(text.replace(/\n/g, '<br/>'));
    });
    $('.lazy').lazyload({
        effect: 'fadeIn',
        effectspeed: 500
    });
    $('button.manage').click(function() {
        var clicked, enable;
        clicked = $(this);
        enable = clicked.text() === 'Enable';
        return $.ajax('/manage', {
            type: 'POST',
            data: {
                'name': clicked.attr('id'),
                'enable': enable
            },
            error: function(jqXHR, textStatus, errorThrown) {
                return alert("AJAX Error: " + textStatus);
            },
            success: function(data, textStatus, jqXHR) {
                var tr;
                tr = clicked.parents('td').parents('tr');
                if (!enable) {
                    tr.removeClass('enable');
                    tr.addClass('disable');
                    return clicked.text('Enable');
                } else {
                    tr.removeClass('disable');
                    tr.addClass('enable');
                    return clicked.text('Disable');
                }
            }
        });
    });
    $('.enable-filter').click(function(e) {
        if ($(this).attr('data-hide') === 'show') {
            $('tr.enable').addClass('hide');
            return $(this).attr('data-hide', 'hide');
        } else {
            $('tr.enable').removeClass('hide');
            return $(this).attr('data-hide', 'show');
        }
    });
    return $('.disable-filter').click(function(e) {
        if ($(this).attr('data-hide') === 'show') {
            $('tr.disable').addClass('hide');
            return $(this).attr('data-hide', 'hide');
        } else {
            $('tr.disable').removeClass('hide');
            return $(this).attr('data-hide', 'show');
        }
    });
});