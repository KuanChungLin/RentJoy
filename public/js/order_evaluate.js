
var currentOrderId; //填寫評價時對應的orderId
$(document).ready(function () {
    $(".rating button").click(function () {
        // 將目前被點擊的星星以及前面的所有星星更改為實心星星
        $(this).prevAll().addBack().text('★');
        // 將目前被點擊的星星之後的星星都更改為空星
        $(this).nextAll().text('☆');
    });
    $(".evaluate-btn").click(function () {
        currentOrderId = $(this).data('order-id'); // 讀取orderId
        console.log(currentOrderId); 
    });
    $("#submitEvaluateButton").click(function (e) {
        e.preventDefault();
        var stars = $('.rating button:contains("★")').length; //計算星數
        var review = $('#evaluate-text').val();  //評價內容
        console.log(stars);
        $.ajax({
            type: "POST",
            url: "/Order/SaveOrderEvaluateToDB",
            data: {
                id: currentOrderId,
                stars: stars,
                review: review
            },
            success: function (response) {
                if (response === "Success") {
                    swal({
                        title: '提交成功!',
                        text: '感謝您的評價!',
                        icon: 'success'
                    }).then((value) => {
                        location.reload();
                    });
                }
            },
            error: function () {
                swal({
                    title: '提交失敗!',
                    text: '請稍後再試，或使用Email聯絡官方人員',
                    icon: 'error'
                });
            }
        });
    });
}); 