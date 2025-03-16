document.addEventListener('DOMContentLoaded', function() {
    // 獲取所有卡片功能區塊
    const cardFuncs = document.querySelectorAll('.card-func');
    
    // 為每個卡片功能區塊添加點擊事件
    cardFuncs.forEach(function(cardFunc) {
        const funcNavIcon = cardFunc.querySelector('.func-nav-icon');
        const funcList = cardFunc.querySelector('.func-list');
        
        // 點擊圖標時顯示功能列表
        funcNavIcon.addEventListener('click', function(e) {
            e.stopPropagation(); // 阻止事件冒泡
            
            // 先隱藏所有其他功能列表
            document.querySelectorAll('.func-list').forEach(function(list) {
                if (list !== funcList) {
                    list.style.display = 'none';
                }
            });
            
            // 切換當前功能列表的顯示狀態
            funcList.style.display = funcList.style.display === 'block' ? 'none' : 'block';
        });
        
        // 防止功能列表內的點擊事件冒泡到文檔
        funcList.addEventListener('click', function(e) {
            e.stopPropagation();
        });
    });
    
    // 點擊文檔其他地方時隱藏所有功能列表
    document.addEventListener('click', function() {
        document.querySelectorAll('.func-list').forEach(function(funcList) {
            funcList.style.display = 'none';
        });
    });
});