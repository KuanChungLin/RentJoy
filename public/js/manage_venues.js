$(document).ready(function() {
    // 獲取所有卡片功能區塊
    $('.card-func').each(function() {
        const $cardFunc = $(this);
        const $funcNavIcon = $cardFunc.find('.func-nav-icon');
        const $funcList = $cardFunc.find('.func-list');
        
        // 點擊圖標時顯示功能列表
        $funcNavIcon.on('click', function(e) {
            e.stopPropagation(); // 阻止事件冒泡
            
            // 先隱藏所有其他功能列表
            $('.func-list').not($funcList).hide();
            
            // 切換當前功能列表的顯示狀態
            $funcList.toggle();
        });
        
        // 防止功能列表內的點擊事件冒泡到文檔
        $funcList.on('click', function(e) {
            e.stopPropagation();
        });
    });
    
    // 點擊文檔其他地方時隱藏所有功能列表
    $(document).on('click', function() {
        $('.func-list').hide();
    });

    // 處理刪除場地
    $('.delete-venue').on('click', function() {
        const venueId = $(this).data('venue-id');

        swal({
            title: "確定要刪除此場地嗎？",
            icon: "warning",
            buttons: ["取消", "確定刪除"],
        })
        .then((willDelete) => {
            if (willDelete) {
                $.ajax({
                    url: '/Manage/DeleteVenue',
                    type: 'POST',
                    data: {
                        venueId: venueId
                    },
                    success: function(data) {
                        if (data === 'Success') {
                            swal({
                                title: "場地已成功刪除！",
                                icon: "success",
                            }).then(() => {
                                location.reload();
                            });
                        } else {
                            swal({
                                title: "場地刪除失敗",
                                text: data,
                                icon: "error",
                            });
                        }
                    },
                    error: function(xhr, status, error) {
                        swal({
                            title: "操作失敗",
                            text: error,
                            icon: "error",
                        });
                    }
                });
            }
        });
    });
    
    // 處理下架場地
    $('.delist-venue').on('click', function() {
        const venueId = $(this).data('venue-id');
        
        swal({
            title: "確定要下架此場地嗎？",
            icon: "warning",
            buttons: ["取消", "確定下架"],
        })
        .then((willDelist) => {
            if (willDelist) {
                $.ajax({
                    url: '/Manage/DelistVenue',
                    type: 'POST',
                    data: {
                        venueId: venueId
                    },
                    success: function(data) {
                        if (data === 'Success') {
                            swal({
                                title: "場地已成功下架！",
                                icon: "success",
                            }).then(() => {
                                location.reload();
                            });
                        } else {
                            swal({
                                title: "場地下架失敗",
                                text: data,
                                icon: "error",
                            });
                        }
                    },
                    error: function(xhr, status, error) {
                        swal({
                            title: "操作失敗",
                            text: error,
                            icon: "error",
                        });
                    }
                });
            }
        });
    });
});