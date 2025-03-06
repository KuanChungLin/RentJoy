$(document).ready(function () {
    $('.details-scheduled.dropdown').on('click', '.dropdown-toggle', function () {
        var $dropdownMenu = $(this).siblings('.dropdown-menu');
        if ($dropdownMenu.is(':hidden')) {
            $dropdownMenu.show();
        } else {
            $dropdownMenu.hide();
        }
    });

    var titleStyle = {
        "font-weight": "700",
        "padding-bottom": "10px",
        "border-bottom": "3px solid #4a4a4a"
    }

    var url = window.location.href;
    if (url.endsWith('OrderReserved')) {
        $('.title-reserved').css(titleStyle)
    } else if (url.endsWith('OrderProcessing')) {
        $('.title-processing').css(titleStyle)
    } else if (url.endsWith('OrderCancel')) {
        $('.title-cancel').css(titleStyle)
    } else if (url.endsWith('OrderFinished')) {
        $('.title-finished').css(titleStyle)
    };

});