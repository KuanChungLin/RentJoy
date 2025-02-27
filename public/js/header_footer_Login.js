let navItems = document.querySelectorAll('.nav-item');

window.addEventListener('load',function () {
    // 初始化時根據當前螢幕尺寸設置子選單的顯示狀態
    initializeSubMenuDisplay();
    
    navItems.forEach(navItem => {
        navItem.onclick = function (event) {
            let screenWidth = window.innerWidth;
                
            if (screenWidth < 1025 && this.querySelector('.sub-menu')) {
                let color = getComputedStyle(navItem).backgroundColor;
                navItem.style.backgroundColor = color === "rgba(0, 0, 0, 0)" ? "#f6f6f6" : "rgba(0, 0, 0, 0)";
            }


            let allSubMenus = document.querySelectorAll('.sub-menu');
            allSubMenus.forEach(subMenu => {
                if (subMenu !== this.querySelector('.sub-menu')) {
                    subMenu.style.display = 'none';
                }
            });


            // 切換當前子選單的顯示狀態
            let currentSubMenu = this.querySelector('.sub-menu');
            if (currentSubMenu) {
                // 在小螢幕下，我們總是切換顯示狀態
                if (screenWidth < 1025) {
                    currentSubMenu.style.display = currentSubMenu.style.display === 'block' ? 'none' : 'block';
                } else {
                    // 在大螢幕下，我們只切換顯示狀態
                    currentSubMenu.style.display = currentSubMenu.style.display === 'block' ? 'none' : 'block';
                }
            }
            

            event.stopPropagation();
        };
    });
    

    
    // 當觸發.faqBtn:hover時，關閉所有.sub-menu
    let faqBtn = document.querySelector('.faq-btn');
    faqBtn.onmouseover = function () {
        if (window.innerWidth > 1024) {
            let subMenus = document.querySelectorAll('.sub-menu');
            subMenus.forEach(subMenu => subMenu.style.display = 'none');
        }
    };

    // 點擊在.nav-item以外的地方，則隱藏所有.sub-menu
    document.onclick = function () {
        if (window.innerWidth > 1024) {
            let subMenus = document.querySelectorAll('.sub-menu');
            subMenus.forEach(subMenu => subMenu.style.display = 'none');
        }
    }

    // 視窗大小變動時，隱藏所有 .sub-menu
    window.onresize = function () {
        if (window.innerWidth > 1025) {
            let subMenus = document.querySelectorAll('.sub-menu');
            subMenus.forEach(subMenu => subMenu.style.display = 'none');
        }
    }
});

// 根據螢幕尺寸初始化子選單顯示狀態
function initializeSubMenuDisplay() {
    const subMenus = document.querySelectorAll('.sub-menu');
    
    if (window.innerWidth <= 1023) {
        // 在小螢幕下，預設隱藏所有子選單
        subMenus.forEach(subMenu => {
            // 只在菜單未被點擊時隱藏
            if (!document.getElementById('menu-switch').checked) {
                subMenu.style.display = 'none';
            }
        });
    } else {
        // 在大螢幕下，隱藏所有子選單
        subMenus.forEach(subMenu => {
            subMenu.style.display = 'none';
        });
    }
}