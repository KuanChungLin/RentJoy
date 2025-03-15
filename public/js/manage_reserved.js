$(document).ready(function() {
    // 處理接受按鈕點擊事件
    $(".accept-btn").click(function() {
        const orderId = $(this).data("order-id");
        
        $.ajax({
            url: "/Manage/ReservedAccept",
            type: "POST",
            data: {
                orderId: orderId
            },
            success: function(response) {
                // 成功處理邏輯
                swal({
                    title: "訂單接受成功！",
                    icon: "success",
                }).then((value) => {
                    location.reload();
                });
            },
            error: function() {
                // 錯誤處理邏輯
                swal({
                    title: "操作失敗，請重新操作或聯絡客服",
                    icon: "warning",
                });
            }
        });
    });
    
    // 處理拒絕按鈕點擊事件
    $(".reject-btn").click(function() {
        const orderId = $(this).data("order-id");
        
        $.ajax({
            url: "/Manage/ReservedReject",
            type: "POST",
            data: {
                orderId: orderId
            },
            success: function(response) {
                // 成功處理邏輯
                swal({
                    title: "訂單拒絕成功！",
                    icon: "success",
                }).then((value) => {
                    location.reload();
                });
            },
            error: function() {
                // 錯誤處理邏輯
                swal({
                    title: "操作失敗，請重新操作或聯絡客服",
                    icon: "warning",
                });
            }
        });
    });
});